package studio

import (
	"github.com/mokiat/lacking-studio/internal/studio/controller"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
)

func BootstrapApplication(
	window *ui.Window,
	registry asset.Registry,
	gfxEngine graphics.Engine,
	physicsEngine *physics.Engine,
	ecsEngine *ecs.Engine,
) {
	studio := controller.NewStudio(
		window,
		registry,
		gfxEngine,
		physicsEngine,
		ecsEngine,
	)
	co.Initialize(window, studio.Render())
}
