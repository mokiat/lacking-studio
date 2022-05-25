package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

var (
	projectDirFlag string
)

func init() {
	flag.StringVar(&projectDirFlag, "project", ".", "project directory")
}

func main() {
	flag.Parse()

	log.Info("Starting studio")
	if err := runApplication(); err != nil {
		log.Error("Studio crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Studio closed")
}

func runApplication() error {
	projectDir, err := evalProjectDir()
	if err != nil {
		return fmt.Errorf("failed to evaluate project dir: %w", err)
	}
	log.Debug("Using project directory %q", projectDir)

	registry, err := asset.NewDirRegistry(projectDir)
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	locator := resource.NewFSLocator(resources.FS)

	cfg := glapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetLocator(locator)
	cfg.SetIcon("icons/favicon.png")

	physicsEngine := physics.NewEngine()
	ecsEngine := ecs.NewEngine()

	renderAPI := glrender.NewAPI()
	graphicsEngine := graphics.NewEngine(renderAPI, glgame.NewShaderCollection())
	resourceLocator := mat.WrappedResourceLocator(locator)

	uiCfg := ui.NewConfig(resourceLocator, renderAPI, glui.NewShaderCollection())
	controller := app.NewLayeredController(
		studio.NewController(graphicsEngine),
		ui.NewController(uiCfg, func(w *ui.Window) {
			studio.BootstrapApplication(projectDir, w, renderAPI, registry, graphicsEngine, physicsEngine, ecsEngine)
		}),
	)

	return glapp.Run(cfg, controller)
}

func evalProjectDir() (string, error) {
	absWorkDir, err := filepath.Abs(projectDirFlag)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute project directory: %w", err)
	}
	return absWorkDir, nil
}
