package bento

type Widget interface {
	Render(area Rect, buffer *Buffer)
}

type StatefulWidget[S any] interface {
	RenderStateful(area Rect, buffer *Buffer, state S)
}
