package throbberwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
)

var _ bento.StatefulWidget[State] = (*Throbber)(nil)

type Throbber struct {
	style bento.Style
	set   Set

	vertical, horizontal bento.Flex
	padding              bento.Padding
}

func New() Throbber {
	return Throbber{
		style:      bento.NewStyle(),
		set:        _lineSet,
		vertical:   bento.FlexCenter,
		horizontal: bento.FlexCenter,
		padding:    bento.NewPadding(),
	}
}

func (t Throbber) WithType(type_ Type) Throbber {
	t.set = type_.Set()
	return t
}

func (t Throbber) RenderStateful(area bento.Rect, buffer *bento.Buffer, state State) {
	frame := t.set.Frame(state.frame)

	width, height := frame.Dimensions()

	vertical := bento.NewLayout(bento.ConstraintLen(height)).Vertical().WithPadding(t.padding).WithFlex(t.vertical)
	horizontal := bento.NewLayout(bento.ConstraintLen(width)).Horizontal().WithPadding(t.padding).WithFlex(t.horizontal)

	area = vertical.Split(area).Unwrap()
	area = horizontal.Split(area).Unwrap()

	text := textwidget.NewTextStr(frame.String()).WithStyle(t.style)

	text.Render(area, buffer)
}
