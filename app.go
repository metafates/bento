package bento

import (
	"fmt"
	"sync"
	"time"
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

	// handlers is a list of channels that need to be waited on before the
	// program can exit.
	handlers channelHandlers

	msgs     chan Msg
	errs     chan error
	finished chan struct{}

	terminal *Terminal
}

func NewApp(initialModel Model) (App, error) {
	backend := NewDefaultBackend()
	terminal, err := NewTerminal(&backend, ViewportFullscreen{})
	if err != nil {
		return App{}, fmt.Errorf("new terminal: %w", err)
	}

	return App{
		initialModel: initialModel,
		handlers:     []chan struct{}{},
		msgs:         make(chan Msg),
		errs:         make(chan error),
		finished:     make(chan struct{}),
		terminal:     terminal,
	}, nil
}

func (a *App) Run() (Model, error) {
	// TODO: everything else
	if err := a.init(); err != nil {
		return a.initialModel, fmt.Errorf("init: %w", err)
	}

	model := a.initialModel

	if initCmd := model.Init(); initCmd != nil {
		_ = initCmd
	}

	a.terminal.Draw(model.Draw)

	time.Sleep(5 * time.Second)

	if err := a.shutdown(); err != nil {
		return a.initialModel, fmt.Errorf("shutdown: %w", err)
	}

	return a.initialModel, nil
}

func (a *App) shutdown() error {
	return a.restore()
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
