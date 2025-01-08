package paragraphwidget

type Wrap struct {
	Trim bool
}

func NewWrap() Wrap {
	return Wrap{}
}

func (w Wrap) WithTrim(trim bool) Wrap {
	w.Trim = trim
	return w
}
