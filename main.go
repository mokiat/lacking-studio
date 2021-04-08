package main

import (
	"log"

	"github.com/mokiat/lacking-studio/internal/studio"
	"github.com/mokiat/lacking/app"
	glfwapp "github.com/mokiat/lacking/glfw/app"
	openglui "github.com/mokiat/lacking/opengl/ui"
	"github.com/mokiat/lacking/ui"
	_ "github.com/mokiat/lacking/ui/standard"
)

func main() {
	cfg := glfwapp.NewConfig("Lacking Studio", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetMaximized(true)
	cfg.SetIcon("resources/studio/icon.png")

	gfxCanvas := openglui.NewCanvas()
	controller := app.NewLayeredController(
		ui.NewController(gfxCanvas, studio.Config{}),
	)

	log.Println("running application")
	if err := glfwapp.Run(cfg, controller); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}
