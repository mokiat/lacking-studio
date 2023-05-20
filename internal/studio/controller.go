package studio

import (
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/graphics"
)

func NewController(gfxEngine *graphics.Engine) *Controller {
	return &Controller{
		gfxEngine: gfxEngine,
	}
}

type Controller struct {
	app.NopController
	gfxEngine *graphics.Engine
}

func (c *Controller) OnCreate(window app.Window) {
	c.gfxEngine.Create()
}

func (c *Controller) OnRender(window app.Window) {
	window.Invalidate() // force redraw
}

func (c *Controller) OnDestroy(window app.Window) {
	c.gfxEngine.Destroy()
}
