package paragraphwidget

type Wrap struct {
	Trim bool
}

func NewWrap() Wrap {
	return Wrap{}
}

func (w Wrap) Untrimmed() Wrap {
	return w.WithTrim(false)
}

func (w Wrap) Trimmed() Wrap {
	return w.WithTrim(true)
}

func (w Wrap) WithTrim(trim bool) Wrap {
	w.Trim = trim
	return w
}
