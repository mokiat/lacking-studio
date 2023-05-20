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

type TwoDTextureConfigPropertiesSectionData struct {
	Texture *model.TwoDTexture
}

var TwoDTextureConfigPropertiesSection = co.Define(&twoDTextureConfigPropertiesSectionComponent{})

type twoDTextureConfigPropertiesSectionComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	texture *model.TwoDTexture
}

func (c *twoDTextureConfigPropertiesSectionComponent) OnUpsert() {
	data := co.GetData[TwoDTextureConfigPropertiesSectionData](c.Properties)
	c.texture = data.Texture

	mvc.UseBinding(c.texture, func(change mvc.Change) bool {
		return true // TODO
	})
}

func (c *twoDTextureConfigPropertiesSectionComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})
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

		co.WithChild("wrapping-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope, "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Wrapping:",
			})
		}))

		co.WithChild("wrapping-dropdown", co.New(std.Dropdown, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.DropdownData{
				Items: []std.DropdownItem{
					{Key: asset.WrapModeClampToEdge, Label: "Clamp To Edge"},
					{Key: asset.WrapModeRepeat, Label: "Repeat"},
					{Key: asset.WrapModeMirroredRepeat, Label: "Mirrored Repeat"},
				},
				SelectedKey: c.texture.Wrapping(),
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(c.Scope, action.ChangeTwoDTextureWrapping{
						Texture:  c.texture,
						Wrapping: key.(asset.WrapMode),
					})
				},
			})
		}))

		co.WithChild("filtering-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope, "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Filtering:",
			})
		}))

		co.WithChild("filtering-dropdown", co.New(std.Dropdown, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.DropdownData{
				Items: []std.DropdownItem{
					{Key: asset.FilterModeNearest, Label: "Nearest"},
					{Key: asset.FilterModeLinear, Label: "Linear"},
					{Key: asset.FilterModeAnisotropic, Label: "Anisotropic"},
				},
				SelectedKey: c.texture.Filtering(),
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(c.Scope, action.ChangeTwoDTextureFiltering{
						Texture:   c.texture,
						Filtering: key.(asset.FilterMode),
					})
				},
			})
		}))

		co.WithChild("data-format-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope, "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Data Format:",
			})
		}))

		co.WithChild("data-format-dropdown", co.New(std.Dropdown, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.DropdownData{
				Items: []std.DropdownItem{
					{Key: asset.TexelFormatRGBA8, Label: "RGBA8"},
					{Key: asset.TexelFormatRGBA16F, Label: "RGBA16F"},
					{Key: asset.TexelFormatRGBA32F, Label: "RGBA32F"},
				},
				SelectedKey: c.texture.Format(),
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(c.Scope, action.ChangeTwoDTextureFormat{
						Texture: c.texture,
						Format:  key.(asset.TexelFormat),
					})
				},
			})
		}))
	})
}
