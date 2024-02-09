package model

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
)

type Visualization interface {
	OnViewportRender(framebuffer render.Framebuffer, size ui.Size)
	OnViewportMouseEvent(element *ui.Element, event ui.MouseEvent) bool
}
