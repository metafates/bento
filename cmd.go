package bento

type (
	QuitMsg       struct{}
	BatchMsg      []Cmd
	WindowSizeMsg Size

	// FocusMsg represents a terminal focus message.
	// This occurs when the terminal gains focus.
	FocusMsg struct{}

	// BlurMsg represents a terminal blur message.
	// This occurs when the terminal loses focus.
	BlurMsg struct{}
)

type (
	sequenceMsg []Cmd
)

// Quit is a special command that tells the Bento app to exit.
func Quit() Msg {
	return QuitMsg{}
}

// Sequence runs the given commands one at a time, in order. Contrast this with
// Batch, which runs commands concurrently.
func Sequence(cmds ...Cmd) Cmd {
	return func() Msg {
		return sequenceMsg(cmds)
	}
}

// Batch performs a bunch of commands concurrently with no ordering guarantees
// about the results. Use a Batch to return several commands.
func Batch(cmds ...Cmd) Cmd {
	validCmds := make([]Cmd, 0, len(cmds))
	for _, c := range cmds {
		if c == nil {
			continue
		}

		validCmds = append(validCmds, c)
	}

	switch len(validCmds) {
	case 0:
		return nil
	case 1:
		return validCmds[0]
	default:
		return func() Msg {
			return BatchMsg(validCmds)
		}
	}
}
