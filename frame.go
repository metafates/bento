package bento

type CompletedFrame struct {
	buffer *Buffer
	area   Rect
	count  int
}

type Frame struct {
	cursorPosition *Position
	viewport       Rect
	buffer         *Buffer
	count          int
}

func (f *Frame) RenderWidget(widget Widget, area Rect) {
	widget.Render(area, f.buffer)
}
