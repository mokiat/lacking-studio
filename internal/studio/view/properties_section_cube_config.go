package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type CubeTextureConfigPropertiesSectionData struct {
	Texture *model.CubeTexture
}

var CubeTextureConfigPropertiesSection = mvc.Wrap(co.Define(&cubeTextureConfigPropertiesSectionComponent{}))

type cubeTextureConfigPropertiesSectionComponent struct {
	co.BaseComponent

	texture *model.CubeTexture
}

func (c *cubeTextureConfigPropertiesSectionComponent) OnUpsert() {
	data := co.GetData[CubeTextureConfigPropertiesSectionData](c.Properties())
	c.texture = data.Texture

	mvc.UseBinding(c.Scope(), c.texture, func(change mvc.Change) bool {
		return true // TODO
	})
}

func (c *cubeTextureConfigPropertiesSectionComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    5,
				Bottom: 5,
			},
			Layout: layout.Vertical(layout.VerticalSettings{
				ContentAlignment: layout.HorizontalAlignmentLeft,
				ContentSpacing:   5,
			}),
		})
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})

		co.WithChild("filtering-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Filtering:",
			})
		}))

		co.WithChild("filtering-dropdown", co.New(std.Dropdown, func() {
			co.WithData(std.DropdownData{
				Items: []std.DropdownItem{
					{Key: asset.FilterModeNearest, Label: "Nearest"},
					{Key: asset.FilterModeLinear, Label: "Linear"},
					{Key: asset.FilterModeAnisotropic, Label: "Anisotropic"},
				},
				SelectedKey: c.texture.Filtering(),
			})
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(c.Scope(), action.ChangeCubeTextureFiltering{
						Texture:   c.texture,
						Filtering: key.(asset.FilterMode),
					})
				},
			})
		}))

		co.WithChild("data-format-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Data Format:",
			})
		}))

		co.WithChild("data-format-dropdown", co.New(std.Dropdown, func() {
			co.WithData(std.DropdownData{
				Items: []std.DropdownItem{
					{Key: asset.TexelFormatRGBA8, Label: "RGBA8"},
					{Key: asset.TexelFormatRGBA16F, Label: "RGBA16F"},
					{Key: asset.TexelFormatRGBA32F, Label: "RGBA32F"},
				},
				SelectedKey: c.texture.Format(),
			})
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(c.Scope(), action.ChangeCubeTextureFormat{
						Texture: c.texture,
						Format:  key.(asset.TexelFormat),
					})
				},
			})
		}))
	})
}
