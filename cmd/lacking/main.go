package main

import (
	"log"
	"os"

	"github.com/mokiat/lacking-studio/internal/studio"
	"github.com/mokiat/lacking/app"
	glfwapp "github.com/mokiat/lacking/framework/glfw/app"
	glgraphics "github.com/mokiat/lacking/framework/opengl/game/graphics"
	glui "github.com/mokiat/lacking/framework/opengl/ui"
	"github.com/mokiat/lacking/ui"
)

func main() {
	cfg := glfwapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetIcon("resources/icons/favicon.png")

	graphicsEngine := glgraphics.NewEngine()

	dir, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable dir: %v", err)
	}
	log.Printf("executable dir: %s", dir)

	uiGLGraphics := glui.NewGraphics()
	uiController := ui.NewController(ui.FileResourceLocator{}, uiGLGraphics, func(w *ui.Window) {
		studio.BootstrapApplication(w, graphicsEngine)
	})

	controller := app.NewLayeredController(uiController)

	log.Println("running application")
	if err := glfwapp.Run(cfg, controller); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}
