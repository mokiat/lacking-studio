package studio

import (
	"cmp"
	"fmt"

	nativeapp "github.com/mokiat/lacking-native/app"
	nativegame "github.com/mokiat/lacking-native/game"
	nativeui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking-studio/internal"
	"github.com/mokiat/lacking-studio/internal/global"
	"github.com/mokiat/lacking-studio/internal/preview/view"
	"github.com/mokiat/lacking-studio/resources"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
	"github.com/urfave/cli/v2"
)

func runEditorApplication(ctx *cli.Context) error {
	projectDir := cmp.Or(ctx.Args().First(), ".")

	registry, err := createRegistry(projectDir)
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

	cfg := nativeapp.NewConfig("Lacking Studio [Editor Mode]", 1280, 800)
	cfg.SetMaximized(true)
	cfg.SetMinSize(1024, 768)
	cfg.SetVSync(true)
	cfg.SetIcon("icons/favicon.png")
	cfg.SetLocator(locator)
	cfg.SetAudioEnabled(false)
	return nativeapp.Run(cfg, app.NewLayeredController(
		globalController,
		uiController,
	))
}
