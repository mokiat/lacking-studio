package studio

import (
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
)

func NewController(gfxEngine *graphics.Engine) *Controller {
	return &Controller{
		gfxEngine:     gfxEngine,
		physicsEngine: physics.NewEngine(),
		ecsEngine:     ecs.NewEngine(),
	}
}

type Controller struct {
	app.NopController

	gfxEngine     *graphics.Engine
	physicsEngine *physics.Engine
	ecsEngine     *ecs.Engine
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
