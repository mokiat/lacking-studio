package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

type CubeTextureEditorData struct {
	ResourceModel    *model.Resource
	TextureModel     *model.CubeTexture
	EditorModel      *model.CubeTextureEditor
	StudioController StudioController
	EditorController EditorController
	Visualization    model.Visualization
}

var CubeTextureEditor = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data        = co.GetData[CubeTextureEditorData](props)
		editorModel = data.EditorModel
		viz         = data.Visualization
	)

	mvc.UseBinding(editorModel, func(change mvc.Change) bool {
		return true
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: opt.V(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("center", co.New(mat.DropZone, func() {
			co.WithCallbackData(mat.DropZoneCallbackData{
				OnDrop: func(paths []string) bool {
					mvc.Dispatch(scope, action.ChangeCubeTextureContentFromPath{
						Texture: data.TextureModel,
						Path:    paths[0],
					})
					return true
				},
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})

			co.WithChild("viewport", co.New(mat.Viewport, func() {
				co.WithData(mat.ViewportData{
					API: co.GetContext[global.Context]().API,
				})
				co.WithCallbackData(mat.ViewportCallbackData{
					OnMouseEvent: viz.OnViewportMouseEvent,
					OnRender:     viz.OnViewportRender,
				})
			}))
		}))

		if editorModel.IsPropertiesVisible() {
			co.WithChild("right", co.New(CubeTextureProperties, func() {
				co.WithData(CubeTexturePropertiesData{
					Model:            editorModel.Properties(),
					ResourceModel:    data.ResourceModel,
					TextureModel:     data.TextureModel,
					StudioController: data.StudioController,
					EditorController: data.EditorController,
				})
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     opt.V(500),
				})
			}))
		}
	})
})
