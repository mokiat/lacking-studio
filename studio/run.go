package studio

import (
	"fmt"
	"os"
	"path/filepath"

	glapp "github.com/mokiat/lacking-native/app"
	glui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking-studio/internal"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
)

type StudioOption func(c *config)

func WithProjectDir(dir string) StudioOption {
	return func(c *config) {
		c.projectDir = dir
	}
}

type config struct {
	projectDir string
}

func Run(opts ...StudioOption) {
	cfg := &config{
		projectDir: ".",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	log.Info("Studio started")
	if err := runApplication(cfg.projectDir); err != nil {
		log.Error("Studio crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Studio stopped")
}

func runApplication(projectDir string) error {
	storage, err := asset.NewFSStorage(filepath.Join(projectDir, "assets"))
	if err != nil {
		return fmt.Errorf("error creating registry storage: %w", err)
	}
	formatter := asset.NewBlobFormatter() // TODO: Make this configurable
	registry, err := asset.NewRegistry(storage, formatter)
	if err != nil {
		return fmt.Errorf("error creating registry: %w", err)
	}

	locator := ui.WrappedLocator(resource.NewFSLocator(resources.FS))
	uiController := ui.NewController(locator, glui.NewShaderCollection(), func(window *ui.Window) {
		internal.BootstrapApplication(window, registry)
	})

	cfg := glapp.NewConfig("Lacking Studio", 1280, 800)
	cfg.SetMaximized(true)
	cfg.SetMinSize(1024, 768)
	cfg.SetVSync(true)
	cfg.SetIcon("icons/favicon.png")
	cfg.SetLocator(locator)
	return glapp.Run(cfg, uiController)
}
