package bento

type Viewport interface {
	isViewport()
}

type (
	ViewportFullscreen struct{}
	ViewportInline     int
	ViewportFixed      Rect
)

func (ViewportFullscreen) isViewport() {}
func (ViewportInline) isViewport()     {}
func (ViewportFixed) isViewport()      {}
