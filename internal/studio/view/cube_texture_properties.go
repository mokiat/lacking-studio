package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var CubeTextureProperties = co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.CubeTextureEditor)

	// WithNotifications(editor.Target(), func(change observer.Change) bool {
	// 	return true // TODO
	// })

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
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
		co.WithLayoutData(props.LayoutData())

		co.WithChild("asset", co.New(mat.Accordion, func() {
			co.WithData(mat.AccordionData{
				Title:    "Asset",
				Expanded: editor.IsAssetAccordionExpanded(),
			})
			co.WithLayoutData(props.LayoutData())
			co.WithCallbackData(mat.AccordionCallbackData{
				OnToggle: func() {
					editor.SetAssetAccordionExpanded(!editor.IsAssetAccordionExpanded())
				},
			})

			// co.WithChild("content", co.New(AssetPropContent, func() {
			// 	co.WithData(AssetPropContentData{
			// 		Model:      data.ResourceModel,
			// 		Controller: data.Controller,
			// 	})
			// }))
		}))

		co.WithChild("config", co.New(CubeTextureConfig, func() {
			co.WithData(editor)
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))
	})
})
