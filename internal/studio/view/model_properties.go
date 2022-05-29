package view

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var ModelProperties = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	// editor := props.Data().(model.ModelEditor)

	return co.New(mat.Element, func() {})

	// return co.New(mat.Container, func() {
	// 	co.WithData(mat.ContainerData{
	// 		Padding: ui.Spacing{
	// 			Left:   5,
	// 			Right:  5,
	// 			Top:    5,
	// 			Bottom: 5,
	// 		},
	// 		Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
	// 			ContentAlignment: mat.AlignmentLeft,
	// 			ContentSpacing:   5,
	// 		}),
	// 	})
	// 	co.WithLayoutData(props.LayoutData())

	// 	co.WithChild("asset", co.New(AssetAccordion, func() {
	// 		co.WithData(AssetAccordionData{
	// 			AssetID:   editor.ID(),
	// 			AssetName: editor.Name(),
	// 			AssetType: "Model",
	// 			Expanded:  editor.IsAssetAccordionExpanded(),
	// 		})
	// 		co.WithLayoutData(mat.LayoutData{
	// 			GrowHorizontally: true,
	// 		})
	// 		co.WithCallbackData(AssetAccordionCallbackData{
	// 			OnToggleExpanded: func() {
	// 				editor.SetAssetAccordionExpanded(!editor.IsAssetAccordionExpanded())
	// 			},
	// 		})
	// 	}))

	// 	// co.WithChild("config", co.New(TwoDTextureConfig, func() {
	// 	// 	co.WithData(editor)
	// 	// 	co.WithLayoutData(mat.LayoutData{
	// 	// 		GrowHorizontally: true,
	// 	// 	})
	// 	// }))
	// })
}))
