package main

import (
	"os"

	glapp "github.com/mokiat/lacking-native/app"
	glui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking-studio/internal"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
)

func main() {
	log.Info("Started")
	if err := runApplication(); err != nil {
		log.Error("Crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Stopped")
}

func runApplication() error {
	locator := ui.WrappedLocator(resource.NewFSLocator(resources.FS))

	uiController := ui.NewController(
		locator,
		glui.NewShaderCollection(),
		internal.BootstrapApplication,
	)

	cfg := glapp.NewConfig("Lacking Studio", 1280, 800)
	cfg.SetMaximized(false)
	cfg.SetMinSize(1024, 768)
	cfg.SetVSync(true)
	cfg.SetIcon("icons/favicon.png")
	cfg.SetLocator(locator)
	return glapp.Run(cfg, uiController)
}
