package bento

import (
	"context"
	"sync"
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

	ctx    context.Context
	cancel context.CancelFunc

	msgs     chan Msg
	errs     chan error
	finished chan struct{}

	terminal *Terminal
}

func (a *App) Run() (Model, error) {
	// TODO: everything else

	cmds := make(chan Cmd)

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

	a.terminal.Draw(model.Draw)

	a.shutdown()

	return a.initialModel, nil
}

func (a *App) shutdown() {
	a.cancel()
	a.handlers.shutdown()
}

func (a *App) restoreTerminal() {
	a.terminal.ShowCursor()

	// TODO
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
