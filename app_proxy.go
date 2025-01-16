package bento

type AppProxy struct {
	runner *appRunner
}

func (p AppProxy) Send(cmd Cmd) {
	p.runner.handleCmd(cmd)
}

func (p AppProxy) Quit() {
	p.Send(Quit)
}
