package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type VisualizationController interface {
	OnViewportKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool
	OnViewportMouseEvent(element *ui.Element, event ui.MouseEvent) bool
	OnViewportRender(framebuffer render.Framebuffer, size ui.Size)
}

type VisualizationData struct {
	Controller VisualizationController
}

var Visualization = co.Define(&visualizationComponent{})

type visualizationComponent struct {
	co.BaseComponent

	api        render.API
	controller VisualizationController
}

func (c *visualizationComponent) OnUpsert() {
	context := co.TypedValue[global.Context](c.Scope())
	c.api = context.API

	data := co.GetData[VisualizationData](c.Properties())
	c.controller = data.Controller
}

func (c *visualizationComponent) Render() co.Instance {
	return co.New(std.Viewport, func() {
		co.WithData(std.ViewportData{
			API: c.api,
		})
		co.WithCallbackData(std.ViewportCallbackData{
			OnMouseEvent:    c.controller.OnViewportMouseEvent,
			OnKeyboardEvent: c.controller.OnViewportKeyboardEvent,
			OnRender:        c.controller.OnViewportRender,
		})
	})
}
