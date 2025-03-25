package paragraphwidget

// Wrap describes how to wrap text across lines.
type Wrap struct {
	// Trim leading whitespace
	Trim bool
}

func NewWrap() Wrap {
	return Wrap{}
}

func (w Wrap) WithTrim(trim bool) Wrap {
	w.Trim = trim
	return w
}
