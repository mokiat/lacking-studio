package studio

import (
	"github.com/mokiat/lacking-studio/internal/studio/controller"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
)

func BootstrapApplication(globalCtx global.Context) {
	co.RegisterContext(globalCtx)
	co.Initialize(globalCtx.Window, co.New(Bootstrap, nil))
}

var Bootstrap = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		globalCtx = co.GetContext[global.Context]()
	)
	studioModel := co.UseState(func() *model.Studio {
		return model.NewStudio()
	})
	studioController := co.UseState(func() *controller.Studio {
		return controller.NewStudio(globalCtx, studioModel.Get())
	})
	return studioController.Get().Render(scope)
})
