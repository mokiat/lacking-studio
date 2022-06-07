package studio

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/controller"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
)

func BootstrapApplication(globalCtx global.Context) error {
	registryModel, err := model.NewRegistry(globalCtx.Registry)
	if err != nil {
		return fmt.Errorf("error creating registry model: %w", err)
	}

	co.RegisterContext(globalCtx)
	co.Initialize(globalCtx.Window, co.New(Bootstrap, func() {
		co.WithData(registryModel)
	}))
	return nil
}

var Bootstrap = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		globalCtx     = co.GetContext[global.Context]()
		registryModel = co.GetData[*model.Registry](props)
	)
	studioModel := co.UseState(func() *model.Studio {
		return model.NewStudio(registryModel)
	})
	studioController := co.UseState(func() *controller.Studio {
		return controller.NewStudio(globalCtx, studioModel.Get())
	})
	return studioController.Get().Render(scope)
})
