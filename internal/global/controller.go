package global

import (
	"github.com/mokiat/lacking-studio/internal/viewport"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
)

func NewController(gameController *game.Controller) *Controller {
	return &Controller{
		Controller: gameController,
	}
}

var _ app.Controller = (*Controller)(nil)

type Controller struct {
	*game.Controller

	commonData *viewport.CommonData
}

func (c *Controller) OnCreate(window app.Window) {
	c.Controller.OnCreate(window)

	gameEngine := c.Controller.Engine()
	gfxEngine := gameEngine.Graphics()

	c.commonData = viewport.NewCommonData(gfxEngine)
	c.commonData.Create()
}

func (c *Controller) OnDestroy(window app.Window) {
	c.commonData.Delete()

	c.Controller.OnDestroy(window)
}

func (c *Controller) CommonData() *viewport.CommonData {
	return c.commonData
}
