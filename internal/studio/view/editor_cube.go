package view

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type CubeTextureEditorData struct {
	ResourceModel *model.Resource
	TextureModel  *model.CubeTexture
	EditorModel   *model.CubeTextureEditor
	Visualization model.Visualization
	Controller    Controller
}

var CubeTextureEditor = co.Define(func(props co.Properties) co.Instance {
	var (
		data        = co.GetData[CubeTextureEditorData](props)
		editorModel = data.EditorModel
		viz         = data.Visualization
		controller  = data.Controller
	)

	WithBinding(editorModel, func(change observer.Change) bool {
		return true
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("center", co.New(mat.DropZone, func() {
			co.WithCallbackData(mat.DropZoneCallbackData{
				OnDrop: func(paths []string) bool {
					controller.Dispatch(action.ChangeCubeTextureContentFromPath{
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
					Model:         editorModel.Properties(),
					ResourceModel: data.ResourceModel,
					TextureModel:  data.TextureModel,
					Controller:    data.Controller,
				})
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     optional.Value(500),
				})
			}))
		}
	})
})
