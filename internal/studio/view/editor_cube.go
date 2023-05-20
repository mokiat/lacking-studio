package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type CubeTextureEditorData struct {
	ResourceModel    *model.Resource
	TextureModel     *model.CubeTexture
	EditorModel      *model.CubeTextureEditor
	StudioController StudioController
	EditorController EditorController
	Visualization    model.Visualization
}

var CubeTextureEditor = co.ContextScoped(co.Define(&cubeTextureEditorComponent{}))

type cubeTextureEditorComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	resourceModel    *model.Resource
	textureModel     *model.CubeTexture
	editorModel      *model.CubeTextureEditor
	studioController StudioController
	editorController EditorController
	viz              model.Visualization
}

func (c *cubeTextureEditorComponent) OnUpsert() {
	data := co.GetData[CubeTextureEditorData](c.Properties)
	c.resourceModel = data.ResourceModel
	c.textureModel = data.TextureModel
	c.editorModel = data.EditorModel
	c.studioController = data.StudioController
	c.editorController = data.EditorController
	c.viz = data.Visualization

	mvc.UseBinding(c.editorModel, func(change mvc.Change) bool {
		return true
	})
}

func (c *cubeTextureEditorComponent) Render() co.Instance {
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
					mvc.Dispatch(c.Scope, action.ChangeCubeTextureContentFromPath{
						Texture: c.textureModel,
						Path:    paths[0],
					})
					return true
				},
			})

			co.WithChild("viewport", co.New(std.Viewport, func() {
				co.WithData(std.ViewportData{
					API: co.TypedValue[global.Context](c.Scope).API,
				})
				co.WithCallbackData(std.ViewportCallbackData{
					OnKeyboardEvent: func(event ui.KeyboardEvent) bool { return false },
					OnMouseEvent:    c.viz.OnViewportMouseEvent,
					OnRender:        c.viz.OnViewportRender,
				})
			}))
		}))

		if c.editorModel.IsPropertiesVisible() {
			co.WithChild("right", co.New(CubeTextureProperties, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment:   layout.VerticalAlignmentCenter,
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(500),
				})
				co.WithData(CubeTexturePropertiesData{
					Model:            c.editorModel.Properties(),
					ResourceModel:    c.resourceModel,
					TextureModel:     c.textureModel,
					StudioController: c.studioController,
					EditorController: c.editorController,
				})
			}))
		}
	})
}
