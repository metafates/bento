package bento

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/muesli/cancelreader"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInterrupted = errors.New("interrupted")
	ErrKilled      = errors.New("killed")
)

type Msg any

type Cmd func() Msg

type Model interface {
	Init() Cmd
	Update(msg Msg) (Model, Cmd)
	Draw(frame *Frame)
}

type App struct {
	initialModel Model

	ctx       context.Context
	cancelCtx context.CancelFunc

	cancelReader cancelreader.CancelReader
	readLoopDone chan struct{}

	// handlers is a list of channels that need to be waited on before the
	// program can exit.
	handlers channelHandlers

	msgs     chan Msg
	errs     chan error
	finished chan struct{}

	terminal *Terminal
}

func NewApp(ctx context.Context, initialModel Model) (App, error) {
	backend := NewDefaultBackend()
	terminal, err := NewTerminal(&backend, ViewportFullscreen{})
	if err != nil {
		return App{}, fmt.Errorf("new terminal: %w", err)
	}

	ctx, cancelCtx := context.WithCancel(ctx)

	return App{
		initialModel: initialModel,
		ctx:          ctx,
		cancelCtx:    cancelCtx,
		readLoopDone: make(chan struct{}),
		handlers:     []chan struct{}{},
		msgs:         make(chan Msg),
		errs:         make(chan error),
		finished:     make(chan struct{}),
		terminal:     terminal,
	}, nil
}

func (a *App) Run() (Model, error) {
	a.handlers = channelHandlers{}
	cmds := make(chan Cmd)
	a.errs = make(chan error)
	a.finished = make(chan struct{}, 1)

	// TODO: everything else
	if err := a.init(); err != nil {
		return a.initialModel, fmt.Errorf("init: %w", err)
	}

	model := a.initialModel

	if initCmd := model.Init(); initCmd != nil {
		ch := make(chan struct{})
		a.handlers.add(ch)

		go func() {
			defer close(ch)

			select {
			case cmds <- initCmd:
			case <-a.ctx.Done():
			}
		}()
	}

	a.draw(model.Draw)

	if err := a.initCancelReader(); err != nil {
		return model, fmt.Errorf("init cancel reader: %w", err)
	}

	// Handle resize events.
	a.handlers.add(a.handleResize())

	// Process commands.
	a.handlers.add(a.handleCommands(cmds))

	model, err := a.eventLoop(model, cmds)
	killed := a.ctx.Err() != nil || err != nil
	if killed && err == nil {
		err = fmt.Errorf("%w: %s", ErrKilled, a.ctx.Err())
	}

	if err == nil {
		a.draw(model.Draw)
	}

	if errShutdown := a.shutdown(); errShutdown != nil {
		return a.initialModel, fmt.Errorf("shutdown: %w", errors.Join(errShutdown, err))
	}

	return a.initialModel, err
}

func (a *App) Send(msg Msg) {
	select {
	case <-a.ctx.Done():
	case a.msgs <- msg:
	}
}

func (a *App) initCancelReader() error {
	r, err := cancelreader.NewReader(a.terminal)
	if err != nil {
		return fmt.Errorf("new reader: %w", err)
	}

	a.cancelReader = r
	a.readLoopDone = make(chan struct{})

	go a.readLoop()

	return nil
}

func (a *App) readLoop() {
	defer close(a.readLoopDone)

	err := readInputs(a.ctx, a.msgs, a.cancelReader)
	if !errors.Is(err, io.EOF) && !errors.Is(err, cancelreader.ErrCanceled) {
		select {
		case <-a.ctx.Done():
		case a.errs <- err:
		}
	}
}

// handleCommands runs commands in a goroutine and sends the result to the
// program's message channel.
func (a *App) handleCommands(cmds chan Cmd) chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		for {
			select {
			case <-a.ctx.Done():
				return

			case cmd := <-cmds:
				if cmd == nil {
					continue
				}

				// Don't wait on these goroutines, otherwise the shutdown
				// latency would get too large as a Cmd can run for some time
				// (e.g. tick commands that sleep for half a second). It's not
				// possible to cancel them so we'll have to leak the goroutine
				// until Cmd returns.
				go func() {
					msg := cmd() // this can be long.
					a.Send(msg)
				}()
			}
		}
	}()

	return ch
}

func (a *App) handleResize() chan struct{} {
	ch := make(chan struct{})

	// Get the initial terminal size and send it to the program.
	go a.checkResize()

	// Listen for window resizes.
	go a.listenForResize(ch)

	return ch
}

func (a *App) listenForResize(done chan struct{}) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)

	defer func() {
		signal.Stop(sig)
		close(done)
	}()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-sig:
		}

		a.checkResize()
	}
}

// checkResize detects the current size of the output and informs the program
// via a WindowSizeMsg.
func (a *App) checkResize() {
	size, err := a.terminal.Size()
	if err != nil {
		select {
		case <-a.ctx.Done():
		case a.errs <- err:
		}

		return
	}

	a.Send(WindowSizeMsg(size))
}

func (a *App) eventLoop(model Model, cmds chan Cmd) (Model, error) {
	for {
		select {
		case <-a.ctx.Done():
			return model, nil

		case err := <-a.errs:
			return model, err

		case msg := <-a.msgs:
			if msg == nil {
				continue
			}

			switch msg := msg.(type) {
			case QuitMsg:
				return model, nil
			case sequenceMsg:
				go func() {
					// Execute commands one at a time, in order.
					for _, cmd := range msg {
						if cmd == nil {
							continue
						}

						msg := cmd()
						if batchMsg, ok := msg.(BatchMsg); ok {
							g, _ := errgroup.WithContext(a.ctx)
							for _, cmd := range batchMsg {
								cmd := cmd
								g.Go(func() error {
									a.Send(cmd())
									return nil
								})
							}

							_ = g.Wait() // wait for all commands from batch msg to finish
							continue
						}

						a.Send(msg)
					}
				}()
			}

			var cmd Cmd
			model, cmd = model.Update(msg) // run update
			cmds <- cmd                    // process command (if any)
			a.draw(model.Draw)             // send view to renderer
		}
	}
}

func (a *App) shutdown() error {
	a.cancelCtx()

	a.handlers.shutdown()

	return a.restore()
}

func (a *App) draw(draw func(*Frame)) {
	_, err := a.terminal.Draw(draw)
	if err != nil {
		a.errs <- err
	}
}

func (a *App) init() error {
	if err := a.terminal.EnableRawMode(); err != nil {
		return fmt.Errorf("enable raw mode: %w", err)
	}

	if err := a.terminal.EnableAlternateScreen(); err != nil {
		return fmt.Errorf("enable alt screen buffer: %w", err)
	}

	return nil
}

func (a *App) restore() error {
	if err := a.terminal.DisableRawMode(); err != nil {
		return fmt.Errorf("disable raw mode: %w", err)
	}

	if err := a.terminal.LeaveAlternateScreen(); err != nil {
		return fmt.Errorf("leave alt screen buffer: %w", err)
	}

	if err := a.terminal.ShowCursor(); err != nil {
		return fmt.Errorf("show cursor: %w", err)
	}

	return nil
}

// channelHandlers manages the series of channels returned by various processes.
// It allows us to wait for those processes to terminate before exiting the
// program.
type channelHandlers []chan struct{}

// Adds a channel to the list of handlers. We wait for all handlers to terminate
// gracefully on shutdown.
func (h *channelHandlers) add(ch chan struct{}) {
	*h = append(*h, ch)
}

// shutdown waits for all handlers to terminate.
func (h channelHandlers) shutdown() {
	var wg sync.WaitGroup
	for _, ch := range h {
		wg.Add(1)
		go func(ch chan struct{}) {
			<-ch
			wg.Done()
		}(ch)
	}
	wg.Wait()
}

func readInputs(ctx context.Context, msgs chan<- Msg, input io.Reader) error {
	return readAnsiInputs(ctx, msgs, input)
}
