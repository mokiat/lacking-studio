package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	glapp "github.com/mokiat/lacking-native/app"
	glgame "github.com/mokiat/lacking-native/game"
	glrender "github.com/mokiat/lacking-native/render"
	glui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking-studio/internal/studio"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
)

func main() {
	flag.Parse()
	projectDir := "."
	if flag.NArg() > 0 {
		projectDir = flag.Arg(0)
	}

	log.Info("Starting studio")
	if err := runApplication(projectDir); err != nil {
		log.Error("Studio crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Studio closed")
}

func runApplication(projectDir string) error {
	registry, err := asset.NewDirRegistry(projectDir)
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	locator := resource.NewFSLocator(resources.FS)

	physicsEngine := physics.NewEngine(16 * time.Millisecond)
	ecsEngine := ecs.NewEngine()
	renderAPI := glrender.NewAPI()
	graphicsEngine := graphics.NewEngine(renderAPI, glgame.NewShaderCollection())

	controller := app.NewLayeredController(
		studio.NewController(graphicsEngine),
		ui.NewController(ui.WrappedLocator(locator), glui.NewShaderCollection(), func(w *ui.Window) {
			ctx := global.Context{
				Window:         w,
				API:            renderAPI,
				Registry:       registry,
				GraphicsEngine: graphicsEngine,
				PhysicsEngine:  physicsEngine,
				ECSEngine:      ecsEngine,
			}
			if err := studio.BootstrapApplication(ctx); err != nil {
				log.Error("Error bootstrapping application: %v", err)
				w.Close()
			}
		}),
	)

	cfg := glapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetLocator(locator)
	cfg.SetIcon("icons/favicon.png")
	return glapp.Run(cfg, controller)
}
