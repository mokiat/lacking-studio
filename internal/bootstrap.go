package internal

import (
	"github.com/mokiat/lacking-studio/internal/global"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
)

func BootstrapApplication(window *ui.Window, globalController *global.Controller, component co.Component) {
	eventBus := mvc.NewEventBus()

	scope := co.RootScope(window)
	scope = co.TypedValueScope(scope, eventBus)
	scope = co.TypedValueScope(scope, &global.Context{
		EventBus:   eventBus,
		Registry:   globalController.Registry(),
		GameEngine: globalController.Engine(),
		CommonData: globalController.CommonData(),
	})
	co.Initialize(scope, co.New(component, nil))
}
