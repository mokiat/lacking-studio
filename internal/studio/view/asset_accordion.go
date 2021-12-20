package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type AssetAccordionData struct {
	Expanded  bool
	AssetID   string
	AssetName string
	AssetType string
}

type AssetAccordionCallbackData struct {
	OnToggleExpanded func()
	OnNameChanged    func(name string)
}

var AssetAccordion = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data AssetAccordionData
	props.InjectData(&data)

	var callbackData AssetAccordionCallbackData
	props.InjectCallbackData(&callbackData)

	return co.New(widget.Accordion, func() {
		co.WithData(widget.AccordionData{
			Title:    "Asset",
			Expanded: data.Expanded,
		})
		co.WithLayoutData(props.LayoutData())
		co.WithCallbackData(widget.AccordionCallbackData{
			OnToggle: callbackData.OnToggleExpanded,
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

			co.WithChild("id", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   10,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "bold"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      "ID:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      data.AssetID,
					})
				}))
			}))

			co.WithChild("type", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   10,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "bold"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      "Type:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      data.AssetType,
					})
				}))
			}))

			co.WithChild("name", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   10,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "bold"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      "Name:",
					})
				}))

				co.WithChild("value", co.New(widget.Editbox, func() {
					co.WithData(widget.EditboxData{
						Text: data.AssetName,
					})
					co.WithLayoutData(mat.LayoutData{
						Height: optional.NewInt(18),
					})
					co.WithCallbackData(widget.EditboxCallbackData{
						OnChanged: callbackData.OnNameChanged,
					})
				}))
			}))
		}))
	})
}))
