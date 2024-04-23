package internal

import (
	glgame "github.com/mokiat/lacking-native/game"
	"github.com/mokiat/lacking-studio/internal/view"
	"github.com/mokiat/lacking-studio/internal/view/editor/viewport"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
)

func BootstrapApplication(window *ui.Window, registry *asset.Registry) {
	eventBus := mvc.NewEventBus()

	gfxEngine := graphics.NewEngine(window.RenderAPI(), glgame.NewShaderCollection(), glgame.NewShaderBuilder())
	commonData := viewport.NewCommonData(gfxEngine)

	scope := co.RootScope(window)
	scope = co.TypedValueScope(scope, eventBus)
	scope = co.TypedValueScope[*asset.Registry](scope, registry)
	scope = co.TypedValueScope[*graphics.Engine](scope, gfxEngine)
	scope = co.TypedValueScope[*viewport.CommonData](scope, commonData)
	co.Initialize(scope, co.New(Bootstrap, nil))
}

var Bootstrap = co.Define(&bootstrapComponent{})

type bootstrapComponent struct {
	co.BaseComponent
}

func (c *bootstrapComponent) Render() co.Instance {
	return co.New(view.Root, nil)
}
