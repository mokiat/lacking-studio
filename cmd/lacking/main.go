package main

import (
	"log"
	"os"

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

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: studio <project_folder>")
	}
	registry, err := asset.NewDirRegistry(os.Args[1])
	if err != nil {
		log.Fatalf("failed to open registry: %v", err)
	}

	cfg := glfwapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetIcon("resources/icons/favicon.png")

	// baseDir, err := evalBaseDir()
	// if err != nil {
	// 	log.Fatalf("failed to evaluate executable dir: %v", err)
	// }
	// log.Printf("EXEC DIR: %s", baseDir)

	graphicsEngine := glgraphics.NewEngine()
	physicsEngine := physics.NewEngine()
	ecsEngine := ecs.NewEngine()

	resourceLocator := ui.NewFileResourceLocator(os.DirFS("."))
	uiGLGraphics := glui.NewGraphics()

	controller := app.NewLayeredController(
		studio.NewController(graphicsEngine),
		ui.NewController(resourceLocator, uiGLGraphics, func(w *ui.Window) {
			studio.BootstrapApplication(w, registry, graphicsEngine, physicsEngine, ecsEngine)
		}),
	)

	log.Println("running application")
	if err := glfwapp.Run(cfg, controller); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}

// func evalBaseDir() (fs.FS, error) {
// 	execPath, err := os.Executable()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get executable path: %w", err)
// 	}

// 	directExecPath, err := filepath.EvalSymlinks(execPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to evaluate symlinks to executable: %w", err)
// 	}

// 	return os.DirFS(filepath.Dir(directExecPath)), nil
// }
