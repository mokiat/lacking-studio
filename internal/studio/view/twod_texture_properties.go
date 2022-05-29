package view

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type TwoDTexturePropertiesData struct {
	Model         *model.TwoDTextureEditorProperties
	ResourceModel *model.Resource
	TextureModel  *model.TwoDTexture
	Controller    Controller
}

// TODO: Move to controller
// func (r *Resource) ChangeName(name string) {
// 	ch := change.Name(r,
// 		change.NameState{
// 			Value: r.Name(),
// 		},
// 		change.NameState{
// 			Value: name,
// 		},
// 	)
// 	if err := r.history.Add(ch); err != nil {
// 		// TODO: Display UI message
// 		panic(fmt.Errorf("error applying change: %w", err))
// 	}
// }

var TwoDTextureProperties = co.Define(func(props co.Properties) co.Instance {
	data := co.GetData[TwoDTexturePropertiesData](props)
	properties := data.Model

	WithBinding(properties, func(change observer.Change) bool {
		return true
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
		co.WithLayoutData(props.LayoutData())

		co.WithChild("asset", co.New(mat.Accordion, func() {
			co.WithData(mat.AccordionData{
				Title:    "Asset",
				Expanded: properties.IsAssetAccordionExpanded(),
			})
			co.WithLayoutData(props.LayoutData())
			co.WithCallbackData(mat.AccordionCallbackData{
				OnToggle: func() {
					// TODO: Go through an action
					properties.SetAssetAccordionExpanded(!properties.IsAssetAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(AssetPropContent, func() {
				co.WithData(AssetPropContentData{
					Model:      data.ResourceModel,
					Controller: data.Controller,
				})
			}))
		}))

		co.WithChild("config", co.New(mat.Accordion, func() {
			co.WithData(mat.AccordionData{
				Title:    "Config",
				Expanded: properties.IsConfigAccordionExpanded(),
			})
			co.WithLayoutData(props.LayoutData())
			co.WithCallbackData(mat.AccordionCallbackData{
				OnToggle: func() {
					// TODO: Go through an action
					properties.SetConfigAccordionExpanded(!properties.IsConfigAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(TwoDTextureConfigPropContent, func() {
				co.WithData(TwoDTextureConfigPropContentData{
					Texture:    data.TextureModel,
					Controller: data.Controller,
				})
			}))
		}))
	})
})
