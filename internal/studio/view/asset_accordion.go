package view

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
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
	data := co.GetData[AssetAccordionData](props)
	callbackData := co.GetCallbackData[AssetAccordionCallbackData](props)

	return co.New(mat.Accordion, func() {
		co.WithData(mat.AccordionData{
			Title:    "Asset",
			Expanded: data.Expanded,
		})
		co.WithLayoutData(props.LayoutData())
		co.WithCallbackData(mat.AccordionCallbackData{
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
						FontSize:  optional.Value(float32(18)),
						FontColor: optional.Value(ui.Black()),
						Text:      "ID:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.Value(float32(18)),
						FontColor: optional.Value(ui.Black()),
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
						FontSize:  optional.Value(float32(18)),
						FontColor: optional.Value(ui.Black()),
						Text:      "Type:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.Value(float32(18)),
						FontColor: optional.Value(ui.Black()),
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
						FontSize:  optional.Value(float32(18)),
						FontColor: optional.Value(ui.Black()),
						Text:      "Name:",
					})
				}))

				co.WithChild("value", co.New(mat.Editbox, func() {
					co.WithData(mat.EditboxData{
						Text: data.AssetName,
					})
					co.WithLayoutData(mat.LayoutData{
						Width: optional.Value(300),
					})
					co.WithCallbackData(mat.EditboxCallbackData{
						OnChanged: callbackData.OnNameChanged,
					})
				}))
			}))
		}))
	})
}))
