package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type BinaryPropertiesData struct {
	Model            *model.BinaryEditorProperties
	ResourceModel    *model.Resource
	BinaryModel      *model.Binary
	StudioController StudioController
	EditorController BinaryEditorController
}

var BinaryProperties = co.Define(&binaryPropertiesComponent{})

type binaryPropertiesComponent struct {
	Properties co.Properties `co:"properties"`

	properties       *model.BinaryEditorProperties
	resourceModel    *model.Resource
	binaryModel      *model.Binary
	studioController StudioController
	editorController BinaryEditorController
}

func (c *binaryPropertiesComponent) OnUpsert() {
	data := co.GetData[BinaryPropertiesData](c.Properties)
	c.properties = data.Model
	c.resourceModel = data.ResourceModel
	c.binaryModel = data.BinaryModel
	c.studioController = data.StudioController
	c.editorController = data.EditorController

	mvc.UseBinding(c.properties, func(change mvc.Change) bool {
		return true
	})
}

func (c *binaryPropertiesComponent) Render() co.Instance {
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
		co.WithLayoutData(c.Properties.LayoutData())

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

		co.WithChild("info", co.New(std.Accordion, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.AccordionData{
				Title:    "Info",
				Expanded: c.properties.IsInfoAccordionExpanded(),
			})
			co.WithCallbackData(std.AccordionCallbackData{
				OnToggle: func() {
					c.properties.SetInfoAccordionExpanded(!c.properties.IsInfoAccordionExpanded())
				},
			})

			co.WithChild("content", co.New(BinaryInfoPropertiesSection, func() {
				co.WithData(BinaryInfoPropertiesSectionData{
					Binary: c.binaryModel,
				})
			}))
		}))
	})
}
