package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var RegistryItem = co.Define(&itemComponent{})

type RegistryItemData struct {
	Resource *asset.Resource
}

type RegistryItemCallbackData struct {
	OnSelected func(resource *asset.Resource)
}

type itemComponent struct {
	co.BaseComponent

	resource *asset.Resource

	onSelected func(resource *asset.Resource)
}

func (c *itemComponent) OnUpsert() {
	data := co.GetData[RegistryItemData](c.Properties())
	c.resource = data.Resource

	callbackData := co.GetCallbackData[RegistryItemCallbackData](c.Properties())
	c.onSelected = callbackData.OnSelected
}

func (c *itemComponent) Render() co.Instance {
	return co.New(std.ListItem, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ListItemData{
			Selected: false,
		})
		co.WithCallbackData(std.ListItemCallbackData{
			OnSelected: c.handleSelected,
		})

		co.WithChild("item", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Padding: ui.UniformSpacing(10),
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("info", co.New(std.Element, func() {
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentAlignment: layout.HorizontalAlignmentLeft,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("name", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      c.resource.Name(),
					})
				}))

				co.WithChild("id", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      c.resource.ID(),
					})
				}))
			}))
		}))
	})
}

func (c *itemComponent) handleSelected() {
	c.onSelected(c.resource)
}
