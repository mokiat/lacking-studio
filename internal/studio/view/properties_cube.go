package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type CubeTexturePropertiesData struct {
	Model            *model.CubeTextureEditorProperties
	ResourceModel    *model.Resource
	TextureModel     *model.CubeTexture
	StudioController StudioController
	EditorController EditorController
}

var CubeTextureProperties = mvc.Wrap(co.Define(&cubeTexturePropertiesComponent{}))

type cubeTexturePropertiesComponent struct {
	co.BaseComponent

	properties       *model.CubeTextureEditorProperties
	resourceModel    *model.Resource
	textureModel     *model.CubeTexture
	studioController StudioController
	editorController EditorController
}

func (c *cubeTexturePropertiesComponent) OnUpsert() {
	data := co.GetData[CubeTexturePropertiesData](c.Properties())
	c.properties = data.Model
	c.resourceModel = data.ResourceModel
	c.textureModel = data.TextureModel
	c.studioController = data.StudioController
	c.editorController = data.EditorController

	mvc.UseBinding(c.Scope(), c.properties, func(change mvc.Change) bool {
		return true
	})
}

func (c *cubeTexturePropertiesComponent) Render() co.Instance {
	return co.New(std.Element, func() {
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
		co.WithLayoutData(c.Properties().LayoutData())

		co.WithChild("asset", co.New(std.Accordion, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.AccordionData{
				Title:    "Asset",
				Expanded: c.properties.IsAssetAccordionExpanded(),
			})
			co.WithCallbackData(std.AccordionCallbackData{
				OnToggle: func() {
					c.properties.SetAssetAccordionExpanded(!c.properties.IsAssetAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(AssetPropertiesSection, func() {
				co.WithData(AssetPropertiesSectionData{
					Model:            c.resourceModel,
					StudioController: c.studioController,
					EditorController: c.editorController,
				})
			}))
		}))

		co.WithChild("config", co.New(std.Accordion, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.AccordionData{
				Title:    "Config",
				Expanded: c.properties.IsConfigAccordionExpanded(),
			})
			co.WithCallbackData(std.AccordionCallbackData{
				OnToggle: func() {
					c.properties.SetConfigAccordionExpanded(!c.properties.IsConfigAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(CubeTextureConfigPropertiesSection, func() {
				co.WithData(CubeTextureConfigPropertiesSectionData{
					Texture: c.textureModel,
				})
			}))
		}))
	})
}
