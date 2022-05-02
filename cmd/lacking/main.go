package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	glapp "github.com/mokiat/lacking-gl/app"
	glgame "github.com/mokiat/lacking-gl/game"
	glrender "github.com/mokiat/lacking-gl/render"
	glui "github.com/mokiat/lacking-gl/ui"
	"github.com/mokiat/lacking-studio/internal/studio"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/ui"
)

var (
	projectDirFlag string
	studioDirFlag  string
)

func init() {
	flag.StringVar(&projectDirFlag, "project", ".", "project directory")
	flag.StringVar(&studioDirFlag, "studio", "", "studio directory")
}

func main() {
	flag.Parse()

	log.Println("running application")
	if err := runApplication(); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}

func runApplication() error {
	studioDir, err := evalStudioDir()
	if err != nil {
		return fmt.Errorf("failed to evaluate studio dir: %w", err)
	}
	log.Printf("studio directory: %s", studioDir)

	projectDir, err := evalProjectDir()
	if err != nil {
		return fmt.Errorf("failed to evaluate project dir: %w", err)
	}
	log.Printf("project directory: %s", projectDir)

	registry, err := asset.NewDirRegistry(projectDir)
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	cfg := glapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetIcon(filepath.Join(studioDir, "resources/icons/favicon.png"))

	physicsEngine := physics.NewEngine()
	ecsEngine := ecs.NewEngine()

	renderAPI := glrender.NewAPI()
	graphicsEngine := graphics.NewEngine(renderAPI, glgame.NewShaderCollection())
	resourceLocator := ui.NewFileResourceLocator(studioDir)

	uiCfg := ui.NewConfig(resourceLocator, renderAPI, glui.NewShaderCollection())
	controller := app.NewLayeredController(
		studio.NewController(graphicsEngine),
		ui.NewController(uiCfg, func(w *ui.Window) {
			studio.BootstrapApplication(projectDir, w, renderAPI, registry, graphicsEngine, physicsEngine, ecsEngine)
		}),
	)

	return glapp.Run(cfg, controller)
}

func evalStudioDir() (string, error) {
	studioDir := studioDirFlag
	if studioDir == "" {
		execPath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to get executable path: %w", err)
		}
		studioDir = filepath.Dir(execPath)
	}
	absStudioDir, err := filepath.Abs(studioDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return absStudioDir, nil
}

func evalProjectDir() (string, error) {
	absWorkDir, err := filepath.Abs(projectDirFlag)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute working directory: %w", err)
	}
	return absWorkDir, nil
}
