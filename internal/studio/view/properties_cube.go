package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

type CubeTexturePropertiesData struct {
	Model            *model.CubeTextureEditorProperties
	ResourceModel    *model.Resource
	TextureModel     *model.CubeTexture
	StudioController StudioController
	EditorController EditorController
}

var CubeTextureProperties = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	data := co.GetData[CubeTexturePropertiesData](props)
	properties := data.Model

	mvc.UseBinding(properties, func(change mvc.Change) bool {
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
					properties.SetAssetAccordionExpanded(!properties.IsAssetAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(AssetPropertiesSection, func() {
				co.WithData(AssetPropertiesSectionData{
					Model:            data.ResourceModel,
					StudioController: data.StudioController,
					EditorController: data.EditorController,
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
					properties.SetConfigAccordionExpanded(!properties.IsConfigAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(CubeTextureConfigPropertiesSection, func() {
				co.WithData(CubeTextureConfigPropertiesSectionData{
					Texture: data.TextureModel,
				})
			}))
		}))
	})
})
