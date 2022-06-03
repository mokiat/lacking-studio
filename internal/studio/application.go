package studio

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/controller"
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
)

func BootstrapApplication(
	window *ui.Window,
	api render.API,
	registry asset.Registry,
	gfxEngine *graphics.Engine,
	physicsEngine *physics.Engine,
	ecsEngine *ecs.Engine,
) error {

	studioRegistry := data.NewRegistry(registry)
	if err := studioRegistry.Init(); err != nil {
		return fmt.Errorf("error initializing registry: %w", err)
	}

	co.RegisterContext(global.Context{
		API:      api,
		Registry: studioRegistry,
	})

	studio := controller.NewStudio(
		window,
		api,
		studioRegistry,
		gfxEngine,
		physicsEngine,
		ecsEngine,
	)
	co.Initialize(window, studio.Render())
	return nil
}
