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
	scope := co.RootScope(globalCtx.Window)
	scope = co.TypedValueScope(scope, globalCtx)
	co.Initialize(scope, co.New(Bootstrap, func() {
		co.WithData(registryModel)
	}))
	return nil
}

var Bootstrap = co.Define(&bootstrapComponent{})

type bootstrapComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	studioController *controller.Studio
}

func (c *bootstrapComponent) OnCreate() {
	globalCtx := co.TypedValue[global.Context](c.Scope)
	registryModel := co.GetData[*model.Registry](c.Properties)
	studioModel := model.NewStudio(registryModel)
	c.studioController = controller.NewStudio(globalCtx, studioModel)
}

func (c *bootstrapComponent) Render() co.Instance {
	return c.studioController.Render(c.Scope)
}
