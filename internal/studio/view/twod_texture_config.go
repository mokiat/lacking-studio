package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

var TwoDTextureConfig = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.TwoDTextureEditor)

	return co.New(mat.Accordion, func() {
		co.WithData(mat.AccordionData{
			Title:    "Config",
			Expanded: editor.IsConfigAccordionExpanded(),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithCallbackData(mat.AccordionCallbackData{
			OnToggle: func() {
				editor.SetConfigAccordionExpanded(!editor.IsConfigAccordionExpanded())
			},
		})

		co.WithChild("content", co.New(mat.Container, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(mat.ContainerData{
				Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
					ContentAlignment: mat.AlignmentLeft,
					ContentSpacing:   5,
				}),
				Padding: ui.Spacing{
					Left:   5,
					Right:  5,
					Top:    5,
					Bottom: 5,
				},
			})

			co.WithChild("wrap-s-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(ui.Black()),
					Text:      "Wrap S:",
				})
			}))

			co.WithChild("wrap-s-dropdown", co.New(mat.Dropdown, func() {
				co.WithData(mat.DropdownData{
					Items: []mat.DropdownItem{
						{Key: asset.WrapModeClampToEdge, Label: "Clamp To Edge"},
						{Key: asset.WrapModeMirroredClampToEdge, Label: "Mirrored Clamp To Edge"},
						{Key: asset.WrapModeRepeat, Label: "Repeat"},
						{Key: asset.WrapModeMirroredRepeat, Label: "Mirrored Repeat"},
					},
					SelectedKey: editor.WrapS(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeWrapS(key.(asset.WrapMode))
					},
				})
			}))

			co.WithChild("wrap-t-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(ui.Black()),
					Text:      "Wrap T:",
				})
			}))

			co.WithChild("wrap-t-dropdown", co.New(mat.Dropdown, func() {
				co.WithData(mat.DropdownData{
					Items: []mat.DropdownItem{
						{Key: asset.WrapModeClampToEdge, Label: "Clamp To Edge"},
						{Key: asset.WrapModeMirroredClampToEdge, Label: "Mirrored Clamp To Edge"},
						{Key: asset.WrapModeRepeat, Label: "Repeat"},
						{Key: asset.WrapModeMirroredRepeat, Label: "Mirrored Repeat"},
					},
					SelectedKey: editor.WrapT(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeWrapT(key.(asset.WrapMode))
					},
				})
			}))

			co.WithChild("min-filter-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(ui.Black()),
					Text:      "Minification Filter:",
				})
			}))

			co.WithChild("min-filter-dropdown", co.New(mat.Dropdown, func() {
				co.WithData(mat.DropdownData{
					Items: []mat.DropdownItem{
						{Key: asset.FilterModeNearest, Label: "Nearest"},
						{Key: asset.FilterModeLinear, Label: "Linear"},
						{Key: asset.FilterModeNearestMipmapNearest, Label: "Nearest Mipmap Nearest"},
						{Key: asset.FilterModeNearestMipmapLinear, Label: "Nearest Mipmap Linear"},
						{Key: asset.FilterModeLinearMipmapNearest, Label: "Linear Mipmap Nearest"},
						{Key: asset.FilterModeLinearMipmapLinear, Label: "Linear Mipmap Linear"},
					},
					SelectedKey: editor.MinFilter(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeMinFilter(key.(asset.FilterMode))
					},
				})
			}))

			co.WithChild("mag-filter-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(ui.Black()),
					Text:      "Magnification Filter:",
				})
			}))

			co.WithChild("mag-filter-dropdown", co.New(mat.Dropdown, func() {
				co.WithData(mat.DropdownData{
					Items: []mat.DropdownItem{
						{Key: asset.FilterModeNearest, Label: "Nearest"},
						{Key: asset.FilterModeLinear, Label: "Linear"},
					},
					SelectedKey: editor.MagFilter(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeMagFilter(key.(asset.FilterMode))
					},
				})
			}))

			co.WithChild("data-format-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(ui.Black()),
					Text:      "Data Format:",
				})
			}))

			co.WithChild("data-format-dropdown", co.New(mat.Dropdown, func() {
				co.WithData(mat.DropdownData{
					Items: []mat.DropdownItem{
						{Key: asset.TexelFormatRGBA8, Label: "RGBA8"},
						{Key: asset.TexelFormatRGBA32F, Label: "RGBA32F"},
					},
					SelectedKey: editor.DataFormat(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeDataFormat(key.(asset.TexelFormat))
					},
				})
			}))
		}))
	})
}))
