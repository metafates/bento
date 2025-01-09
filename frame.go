package bento

type CompletedFrame struct {
	buffer *Buffer
	area   Rect
	count  int
}

type Frame struct {
	cursorPosition *Position
	viewportArea   Rect
	buffer         *Buffer
	count          int
}

func (f *Frame) RenderWidget(widget Widget, area Rect) {
	widget.Render(area, f.buffer)
}

func RenderStatefulWidget[S any](frame *Frame, widget StatefulWidget[S], area Rect, state S) {
	widget.RenderStateful(area, frame.buffer, state)
}

func (f *Frame) Area() Rect {
	return f.viewportArea
}
