package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/optional"
)

type TwoDTextureConfigPropertiesSectionData struct {
	Texture *model.TwoDTexture
}

var TwoDTextureConfigPropertiesSection = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data    = co.GetData[TwoDTextureConfigPropertiesSectionData](props)
		texture = data.Texture
	)

	mvc.UseBinding(texture, func(change mvc.Change) bool {
		return true // TODO
	})

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    5,
				Bottom: 5,
			},
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
				ContentSpacing:   5,
			}),
		})
		co.WithLayoutData(mat.LayoutData{
			GrowHorizontally: true,
		})

		co.WithChild("wrapping-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      "Wrapping:",
			})
		}))

		co.WithChild("wrapping-dropdown", co.New(mat.Dropdown, func() {
			co.WithData(mat.DropdownData{
				Items: []mat.DropdownItem{
					{Key: asset.WrapModeClampToEdge, Label: "Clamp To Edge"},
					{Key: asset.WrapModeRepeat, Label: "Repeat"},
					{Key: asset.WrapModeMirroredRepeat, Label: "Mirrored Repeat"},
				},
				SelectedKey: texture.Wrapping(),
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithCallbackData(mat.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(scope, action.ChangeTwoDTextureWrapping{
						Texture:  texture,
						Wrapping: key.(asset.WrapMode),
					})
				},
			})
		}))

		co.WithChild("filtering-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      "Filtering:",
			})
		}))

		co.WithChild("filtering-dropdown", co.New(mat.Dropdown, func() {
			co.WithData(mat.DropdownData{
				Items: []mat.DropdownItem{
					{Key: asset.FilterModeNearest, Label: "Nearest"},
					{Key: asset.FilterModeLinear, Label: "Linear"},
					{Key: asset.FilterModeAnisotropic, Label: "Anisotropic"},
				},
				SelectedKey: texture.Filtering(),
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithCallbackData(mat.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(scope, action.ChangeTwoDTextureFiltering{
						Texture:   texture,
						Filtering: key.(asset.FilterMode),
					})
				},
			})
		}))

		co.WithChild("data-format-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      "Data Format:",
			})
		}))

		co.WithChild("data-format-dropdown", co.New(mat.Dropdown, func() {
			co.WithData(mat.DropdownData{
				Items: []mat.DropdownItem{
					{Key: asset.TexelFormatRGBA8, Label: "RGBA8"},
					{Key: asset.TexelFormatRGBA16F, Label: "RGBA16F"},
					{Key: asset.TexelFormatRGBA32F, Label: "RGBA32F"},
				},
				SelectedKey: texture.Format(),
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithCallbackData(mat.DropdownCallbackData{
				OnItemSelected: func(key interface{}) {
					mvc.Dispatch(scope, action.ChangeTwoDTextureFormat{
						Texture: texture,
						Format:  key.(asset.TexelFormat),
					})
				},
			})
		}))
	})
})
