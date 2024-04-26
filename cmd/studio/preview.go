package main

import (
	"cmp"
	"fmt"
	"path/filepath"

	nativeapp "github.com/mokiat/lacking-native/app"
	nativegame "github.com/mokiat/lacking-native/game"
	nativeui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking-studio/internal"
	"github.com/mokiat/lacking-studio/internal/global"
	"github.com/mokiat/lacking-studio/internal/preview/view"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
	"github.com/urfave/cli/v2"
)

func runPreviewApplication(ctx *cli.Context) error {
	projectDir := cmp.Or(ctx.Args().First(), ".")

	storage, err := asset.NewFSStorage(filepath.Join(projectDir, "assets"))
	if err != nil {
		return fmt.Errorf("error creating registry storage: %w", err)
	}
	formatter := asset.NewBlobFormatter() // TODO: Make this configurable
	registry, err := asset.NewRegistry(storage, formatter)
	if err != nil {
		return fmt.Errorf("error creating registry: %w", err)
	}

	globalController := global.NewController(
		game.NewController(
			registry,
			nativegame.NewShaderCollection(),
			nativegame.NewShaderBuilder(),
		),
	)

	locator := ui.WrappedLocator(resource.NewFSLocator(resources.FS))
	uiController := ui.NewController(locator, nativeui.NewShaderCollection(), func(window *ui.Window) {
		internal.BootstrapApplication(window, globalController, view.Root)
	})

	cfg := nativeapp.NewConfig("Lacking Studio [Preview Mode]", 1280, 800)
	cfg.SetMaximized(true)
	cfg.SetMinSize(1024, 768)
	cfg.SetVSync(true)
	cfg.SetIcon("icons/favicon.png")
	cfg.SetLocator(locator)
	return nativeapp.Run(cfg, app.NewLayeredController(
		globalController,
		uiController,
	))
}
