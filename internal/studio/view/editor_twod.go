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

type TwoDTextureEditorData struct {
	ResourceModel    *model.Resource
	TextureModel     *model.TwoDTexture
	EditorModel      *model.TwoDTextureEditor
	Visualization    model.Visualization
	StudioController StudioController
	EditorController EditorController
}

var TwoDTextureEditor = co.ContextScoped(co.Define(&twoDTextureEditorComponent{}))

type twoDTextureEditorComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	resourceModel    *model.Resource
	textureModel     *model.TwoDTexture
	editorModel      *model.TwoDTextureEditor
	studioController StudioController
	editorController EditorController
	viz              model.Visualization
}

func (c *twoDTextureEditorComponent) OnUpsert() {
	data := co.GetData[TwoDTextureEditorData](c.Properties)
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

func (c *twoDTextureEditorComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Layout: layout.Frame(),
		})
		co.WithLayoutData(c.Properties.LayoutData())

		co.WithChild("center", co.New(std.DropZone, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment:   layout.VerticalAlignmentCenter,
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
			})
			co.WithCallbackData(std.DropZoneCallbackData{
				OnDrop: func(paths []string) bool {
					mvc.Dispatch(c.Scope, action.ChangeTwoDTextureContentFromPath{
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
			co.WithChild("right", co.New(TwoDTextureProperties, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment:   layout.VerticalAlignmentCenter,
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(500),
				})
				co.WithData(TwoDTexturePropertiesData{
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
