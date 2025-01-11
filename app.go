package bento

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/charmbracelet/x/term"
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

type _Input interface {
	getInput() (io.Reader, func() error, error)
}

type _InputDefault struct{}

func (_InputDefault) getInput() (io.Reader, func() error, error) {
	var input io.Reader = os.Stdin

	// The user has not set a custom input, so we need to check whether or
	// not standard input is a terminal. If it's not, we open a new TTY for
	// input. This will allow things to "just work" in cases where data was
	// piped in or redirected to the application.
	//
	// To disable input entirely pass nil to the [WithInput] program option.
	f, isFile := input.(term.File)
	if !isFile {
		return input, nil, nil
	}
	if term.IsTerminal(f.Fd()) {
		return input, nil, nil
	}

	tty, err := openInputTTY()
	if err != nil {
		return nil, nil, fmt.Errorf("open tty: %w", err)
	}

	return tty, tty.Close, nil
}

type _InputTTY struct{}

func (_InputTTY) getInput() (io.Reader, func() error, error) {
	// Open a new TTY, by request
	f, err := openInputTTY()
	if err != nil {
		return nil, nil, fmt.Errorf("open tty: %w", err)
	}

	return f, f.Close, nil
}

type _InputCustom struct{ io.Reader }

func (i _InputCustom) getInput() (io.Reader, func() error, error) {
	return i.Reader, nil, nil
}

type App struct {
	initialModel Model

	ctx       context.Context
	cancelCtx context.CancelFunc

	input  _Input
	output io.Writer
}

func NewApp(initialModel Model) App {
	ctx, cancelCtx := context.WithCancel(context.Background())

	return App{
		initialModel: initialModel,
		ctx:          ctx,
		cancelCtx:    cancelCtx,
		input:        _InputDefault{},
		output:       os.Stdout,
	}
}

func (a App) WithContext(ctx context.Context) App {
	a.ctx, a.cancelCtx = context.WithCancel(ctx)

	return a
}

func (a App) Run() (Model, error) {
	input, closeInput, err := a.input.getInput()
	if err != nil {
		return nil, fmt.Errorf("get input: %w", err)
	}

	backend := NewDefaultBackend(input, a.output)
	terminal, err := NewTerminal(&backend, ViewportFullscreen{})
	if err != nil {
		return a.initialModel, fmt.Errorf("new terminal: %w", err)
	}

	runner := appRunner{
		initialModel: a.initialModel,
		ctx:          a.ctx,
		cancelCtx:    a.cancelCtx,

		readLoopDone: make(chan struct{}),
		handlers:     channelHandlers{},
		msgs:         make(chan Msg),
		errs:         make(chan error),
		finished:     make(chan struct{}, 1),
		terminal:     terminal,

		onShutdown: closeInput,
	}

	return runner.Run()
}

type appRunner struct {
	initialModel Model

	ctx       context.Context
	cancelCtx context.CancelFunc

	terminal *Terminal

	cancelReader cancelreader.CancelReader
	readLoopDone chan struct{}

	// handlers is a list of channels that need to be waited on before the
	// program can exit.
	handlers channelHandlers

	msgs     chan Msg
	errs     chan error
	finished chan struct{}

	ttyOutput, ttyInput                     term.File
	previousOutputState, previousInputState *term.State

	onShutdown func() error
}

func (a *appRunner) Run() (model Model, err error) {
	defer func() {
		r := recover()

		shutdownErr := a.shutdown()
		if shutdownErr != nil {
			err = errors.Join(err, shutdownErr)
		}

		if r != nil {
			fmt.Printf("Caught panic:\n\n%s\n\nRestoring terminal...\n\n", r)
			debug.PrintStack()
		}
	}()

	err = a.init()
	if err != nil {
		return a.initialModel, fmt.Errorf("init: %w", err)
	}

	cmds := make(chan Cmd)

	model = a.initialModel

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

	err = a.initCancelReader()
	if err != nil {
		return model, fmt.Errorf("init cancel reader: %w", err)
	}

	// Handle resize events.
	a.handlers.add(a.handleResize())

	// Process commands.
	a.handlers.add(a.handleCommands(cmds))

	model, err = a.eventLoop(model, cmds)
	killed := a.ctx.Err() != nil || err != nil
	if killed && err == nil {
		err = fmt.Errorf("%w: %s", ErrKilled, a.ctx.Err())
	}

	if err != nil {
		return model, err
	}

	a.draw(model.Draw)

	return model, nil
}

func (a *appRunner) Send(msg Msg) {
	select {
	case <-a.ctx.Done():
	case a.msgs <- msg:
	}
}

// handleCommands runs commands in a goroutine and sends the result to the
// program's message channel.
func (a *appRunner) handleCommands(cmds chan Cmd) chan struct{} {
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

func (a *appRunner) handleResize() chan struct{} {
	ch := make(chan struct{})

	// Get the initial terminal size and send it to the program.
	go a.checkResize()

	// Listen for window resizes.
	go a.listenForResize(ch)

	return ch
}

func (a *appRunner) eventLoop(model Model, cmds chan Cmd) (Model, error) {
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
			case WindowSizeMsg:
				if err := a.terminal.Resize(NewRect(msg.Width, msg.Height)); err != nil {
					return model, fmt.Errorf("resize: %w", err)
				}
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

func (a *appRunner) shutdown() error {
	a.cancelCtx()

	if a.onShutdown != nil {
		_ = a.onShutdown()
	}

	a.handlers.shutdown()

	return a.restore()
}

func (a *appRunner) draw(draw func(*Frame)) {
	_, err := a.terminal.Draw(draw)
	if err != nil {
		a.errs <- err
	}
}

func (a *appRunner) init() error {
	if err := a.initTerminal(); err != nil {
		return fmt.Errorf("init terminal: %w", err)
	}

	// if err := a.terminal.EnableRawMode(); err != nil {
	// 	return fmt.Errorf("enable raw mode: %w", err)
	// }

	if err := a.terminal.EnableAlternateScreen(); err != nil {
		return fmt.Errorf("enable alt screen buffer: %w", err)
	}

	return nil
}

func (a *appRunner) restore() error {
	a.restoreTerminalState()

	if err := a.terminal.LeaveAlternateScreen(); err != nil {
		return fmt.Errorf("leave alt screen buffer: %w", err)
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
