package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	glapp "github.com/mokiat/lacking-native/app"
	glui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking-studio/internal"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/debug/log"
	asset "github.com/mokiat/lacking/game/newasset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
)

func main() {
	flag.Parse()
	projectDir := "."
	if flag.NArg() > 0 {
		projectDir = flag.Arg(0)
	}

	log.Info("Started")
	if err := runApplication(projectDir); err != nil {
		log.Error("Crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Stopped")
}

func runApplication(projectDir string) error {
	storage, err := asset.NewFSStorage(filepath.Join(projectDir, "assets"))
	if err != nil {
		return fmt.Errorf("error creating registry storage: %w", err)
	}
	formatter := asset.NewJSONFormatter()
	registry, err := asset.NewRegistry(storage, formatter)
	if err != nil {
		return fmt.Errorf("error creating registry: %w", err)
	}

	locator := ui.WrappedLocator(resource.NewFSLocator(resources.FS))
	uiController := ui.NewController(locator, glui.NewShaderCollection(), func(window *ui.Window) {
		internal.BootstrapApplication(window, registry)
	})

	cfg := glapp.NewConfig("Lacking Studio", 1280, 800)
	cfg.SetMaximized(false)
	cfg.SetMinSize(1024, 768)
	cfg.SetVSync(true)
	cfg.SetIcon("icons/favicon.png")
	cfg.SetLocator(locator)
	return glapp.Run(cfg, uiController)
}
