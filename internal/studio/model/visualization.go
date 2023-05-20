package model

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/std"
)

type Visualization interface {
	OnViewportRender(framebuffer render.Framebuffer, size ui.Size)
	OnViewportMouseEvent(event std.ViewportMouseEvent) bool
}
