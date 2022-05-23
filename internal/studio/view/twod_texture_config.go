package view

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

var TwoDTextureConfig = co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.TwoDTextureEditor)

	WithNotifications(editor.Target(), func(change observer.Change) bool {
		return true // TODO
	})

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

			co.WithChild("wrapping-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
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
					SelectedKey: editor.Wrapping(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeWrapping(key.(asset.WrapMode))
					},
				})
			}))

			co.WithChild("filtering-label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "bold"),
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
					SelectedKey: editor.Filtering(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
				co.WithCallbackData(mat.DropdownCallbackData{
					OnItemSelected: func(key interface{}) {
						editor.ChangeFiltering(key.(asset.FilterMode))
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
})
