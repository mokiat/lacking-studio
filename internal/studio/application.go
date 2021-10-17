package studio

import (
	"github.com/mokiat/lacking-studio/internal/studio/controller"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
)

func BootstrapApplication(window *ui.Window, gfxEngine graphics.Engine) {
	studio := controller.NewStudio(window, gfxEngine)

	co.Initialize(window, co.New(co.StoreProvider, func() {
		co.WithData(co.StoreProviderData{
			Entries: []co.StoreProviderEntry{
				co.NewStoreProviderEntry(global.Reducer()),
			},
		})

		co.WithChild("app", studio.Render())
	}))
}
