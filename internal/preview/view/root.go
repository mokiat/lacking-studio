package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/preview/model"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

// TODO: Get dark theme working

var Root = mvc.EventListener(co.Define(&rootComponent{}))

type rootComponent struct {
	co.BaseComponent

	registry *asset.Registry
	appModel *model.AppModel
}

func (c *rootComponent) OnCreate() {
	// TODO: Initialize gfx engine and common data

	c.registry = co.TypedValue[*asset.Registry](c.Scope())

	eventBus := co.TypedValue[*mvc.EventBus](c.Scope())
	c.appModel = model.NewAppModel(eventBus, c.registry)

	co.Window(c.Scope()).SetCloseInterceptor(c.onCloseRequested)
}

func (c *rootComponent) OnDelete() {
	co.Window(c.Scope()).SetCloseInterceptor(nil)

	// TODO: release common data and gfx engine
}

func (c *rootComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout:          layout.Frame(),
		})

		co.WithChild("toolbar", co.New(Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(ToolbarData{
				AppModel: c.appModel,
			})
		}))

		if c.appModel.SelectedResource() == nil {
			co.WithChild("registry", co.New(Registry, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(RegistryData{
					AppModel: c.appModel,
				})
			}))
		} else {
			co.WithChild("viewport", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ContainerData{
					Layout:          layout.Fill(),
					BackgroundColor: opt.V(ui.Black()),
				})
			}))
		}
	})
}

func (c *rootComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case model.SelectedResourceChangedEvent:
		c.Invalidate()
	case model.RefreshErrorEvent:
		// TODO: Open error dialog
		log.Error("Refresh error: %v", event.Err)
	}
}

func (c *rootComponent) onCloseRequested() bool {
	return true
}
