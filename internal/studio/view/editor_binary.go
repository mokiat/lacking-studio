package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type BinaryEditorController interface {
	EditorController
	OnChangeContentFromPath(path string)
}

type BinaryEditorData struct {
	ResourceModel    *model.Resource
	BinaryModel      *model.Binary
	EditorModel      *model.BinaryEditor
	StudioController StudioController
	EditorController BinaryEditorController
}

var BinaryEditor = co.Define(&binaryEditorComponent{})

type binaryEditorComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	resourceModel    *model.Resource
	binaryModel      *model.Binary
	editorModel      *model.BinaryEditor
	studioController StudioController
	controller       BinaryEditorController
}

func (c *binaryEditorComponent) OnUpsert() {
	data := co.GetData[BinaryEditorData](c.Properties)
	c.resourceModel = data.ResourceModel
	c.binaryModel = data.BinaryModel
	c.editorModel = data.EditorModel
	c.studioController = data.StudioController
	c.controller = data.EditorController

	mvc.UseBinding(c.editorModel, func(change mvc.Change) bool {
		return true
	})
}

func (c *binaryEditorComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties.LayoutData())
		co.WithData(std.ElementData{
			Layout: layout.Frame(),
		})

		co.WithChild("center", co.New(std.DropZone, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment:   layout.VerticalAlignmentCenter,
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
			})
			co.WithCallbackData(std.DropZoneCallbackData{
				OnDrop: func(paths []string) bool {
					c.controller.OnChangeContentFromPath(paths[0])
					return true
				},
			})

			co.WithChild("panel", co.New(std.Container, func() {
				co.WithData(std.ContainerData{
					BackgroundColor: opt.V(std.BackgroundColor),
					Layout:          layout.Anchor(),
				})

				co.WithChild("icon", co.New(std.Picture, func() {
					co.WithLayoutData(layout.Data{
						Width:            opt.V(48),
						Height:           opt.V(48),
						HorizontalCenter: opt.V(48),
						VerticalCenter:   opt.V(48),
					})
					co.WithData(std.PictureData{
						Image:      co.OpenImage(c.Scope, "icons/upload.png"),
						ImageColor: opt.V(std.SurfaceColor),
						Mode:       std.ImageModeStretch,
					})
				}))
			}))
		}))

		if c.editorModel.IsPropertiesVisible() {
			co.WithChild("right", co.New(BinaryProperties, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment:   layout.VerticalAlignmentCenter,
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(500),
				})
				co.WithData(BinaryPropertiesData{
					Model:            c.editorModel.Properties(),
					ResourceModel:    c.resourceModel,
					BinaryModel:      c.binaryModel,
					StudioController: c.studioController,
					EditorController: c.controller,
				})
			}))
		}
	})
}
