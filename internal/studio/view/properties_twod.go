package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type TwoDTexturePropertiesData struct {
	Model            *model.TwoDTextureEditorProperties
	ResourceModel    *model.Resource
	TextureModel     *model.TwoDTexture
	StudioController StudioController
	EditorController EditorController
}

var TwoDTextureProperties = mvc.Wrap(co.Define(&twoDTexturePropertiesComponent{}))

type twoDTexturePropertiesComponent struct {
	co.BaseComponent

	properties       *model.TwoDTextureEditorProperties
	resourceModel    *model.Resource
	textureModel     *model.TwoDTexture
	studioController StudioController
	editorController EditorController
}

func (c *twoDTexturePropertiesComponent) OnUpsert() {
	data := co.GetData[TwoDTexturePropertiesData](c.Properties())
	c.properties = data.Model
	c.resourceModel = data.ResourceModel
	c.textureModel = data.TextureModel
	c.studioController = data.StudioController
	c.editorController = data.EditorController
	mvc.UseBinding(c.Scope(), c.properties, func(change mvc.Change) bool {
		return true
	})
}

func (c *twoDTexturePropertiesComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
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

			co.WithChild("content", co.New(TwoDTextureConfigPropertiesSection, func() {
				co.WithData(TwoDTextureConfigPropertiesSectionData{
					Texture: c.textureModel,
				})
			}))
		}))
	})
}
