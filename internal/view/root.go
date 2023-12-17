package view

import (
	"github.com/mokiat/gog/opt"
	appview "github.com/mokiat/lacking-studio/internal/view/app"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Root = mvc.EventListener(co.Define(&rootComponent{}))

type rootComponent struct {
	co.BaseComponent

	eventBus *mvc.EventBus
}

func (c *rootComponent) OnCreate() {
	c.eventBus = co.TypedValue[*mvc.EventBus](c.Scope())

	co.Window(c.Scope()).SetCloseInterceptor(c.onCloseRequested)
}

func (c *rootComponent) OnDelete() {
	co.Window(c.Scope()).SetCloseInterceptor(nil)
}

func (c *rootComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout:          layout.Frame(),
		})

		co.WithChild("header", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(std.ContainerData{
				Layout:          layout.Vertical(),
				BackgroundColor: opt.V(ui.Red()),
			})

			co.WithChild("toolbar", co.New(appview.Toolbar, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})
				co.WithData(appview.ToolbarData{
					// TODO
				})
			}))

			co.WithChild("tabbar", co.New(appview.Tabbar, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})
				co.WithData(appview.TabbarData{
					// 	// TODO
				})
			}))
		}))

		co.WithChild("editors", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ElementData{
				Layout: layout.Fill(),
			})

			// for _, editor := range editors {
			co.WithChild("editor-0101", co.New(appview.Editor, func() {
				co.WithData(appview.EditorData{
					Visible: true,
				})
			}))
			// }
		}))
	})
}

func (c *rootComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	}
}

func (c *rootComponent) onCloseRequested() bool {
	return true
}
