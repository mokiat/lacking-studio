package main

import (
	"flag"
	"fmt"
	"os"

	glapp "github.com/mokiat/lacking-gl/app"
	glgame "github.com/mokiat/lacking-gl/game"
	glrender "github.com/mokiat/lacking-gl/render"
	glui "github.com/mokiat/lacking-gl/ui"
	"github.com/mokiat/lacking-studio/internal/studio"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
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

	physicsEngine := physics.NewEngine()
	ecsEngine := ecs.NewEngine()
	renderAPI := glrender.NewAPI()
	graphicsEngine := graphics.NewEngine(renderAPI, glgame.NewShaderCollection())

	uiCfg := ui.NewConfig(mat.WrappedResourceLocator(locator), renderAPI, glui.NewShaderCollection())
	controller := app.NewLayeredController(
		studio.NewController(graphicsEngine),
		ui.NewController(uiCfg, func(w *ui.Window) {
			if err := studio.BootstrapApplication(w, renderAPI, registry, graphicsEngine, physicsEngine, ecsEngine); err != nil {
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
