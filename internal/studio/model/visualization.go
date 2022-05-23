package model

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
)

type Visualization interface {
	OnViewportRender(framebuffer render.Framebuffer, size ui.Size)
	OnViewportMouseEvent(event mat.ViewportMouseEvent) bool
}
