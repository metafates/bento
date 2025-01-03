package bento

type Widget interface {
	Render(area Rect, buffer *Buffer)
}
