package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type VisualizationController interface {
	OnViewportKeyboardEvent(event ui.KeyboardEvent) bool
	OnViewportMouseEvent(event mat.ViewportMouseEvent) bool
	OnViewportRender(framebuffer render.Framebuffer, size ui.Size)
}

type VisualizationData struct {
	Controller VisualizationController
}

var Visualization = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	ctx := co.GetContext[global.Context]()
	data := co.GetData[VisualizationData](props)
	controller := data.Controller

	return co.New(mat.Viewport, func() {
		co.WithData(mat.ViewportData{
			API: ctx.API,
		})
		co.WithCallbackData(mat.ViewportCallbackData{
			OnMouseEvent:    controller.OnViewportMouseEvent,
			OnKeyboardEvent: controller.OnViewportKeyboardEvent,
			OnRender:        controller.OnViewportRender,
		})
	})
})
