package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mokiat/lacking-studio/internal/studio"
	"github.com/mokiat/lacking/app"
	glfwapp "github.com/mokiat/lacking/framework/glfw/app"
	glgraphics "github.com/mokiat/lacking/framework/opengl/game/graphics"
	glui "github.com/mokiat/lacking/framework/opengl/ui"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
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

	cfg := glfwapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetIcon(filepath.Join(studioDir, "resources/icons/favicon.png"))

	graphicsEngine := glgraphics.NewEngine()
	physicsEngine := physics.NewEngine()
	ecsEngine := ecs.NewEngine()
	resourceLocator := ui.NewFileResourceLocator(studioDir)
	uiGLGraphics := glui.NewGraphics()

	controller := app.NewLayeredController(
		studio.NewController(graphicsEngine),
		ui.NewController(resourceLocator, uiGLGraphics, func(w *ui.Window) {
			studio.BootstrapApplication(projectDir, w, registry, graphicsEngine, physicsEngine, ecsEngine)
		}),
	)

	return glfwapp.Run(cfg, controller)
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
