package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/optional"
)

type TwoDTextureEditorData struct {
	ResourceModel *model.Resource
	TextureModel  *model.TwoDTexture
	EditorModel   *model.TwoDTextureEditor
	Visualization model.Visualization
}

var TwoDTextureEditor = co.ContextScoped(co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data        = co.GetData[TwoDTextureEditorData](props)
		editorModel = data.EditorModel
		viz         = data.Visualization
	)

	mvc.UseBinding(editorModel, func(change mvc.Change) bool {
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
					mvc.Dispatch(scope, action.ChangeTwoDTextureContentFromPath{
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
			co.WithChild("right", co.New(TwoDTextureProperties, func() {
				co.WithData(TwoDTexturePropertiesData{
					Model:         editorModel.Properties(),
					ResourceModel: data.ResourceModel,
					TextureModel:  data.TextureModel,
				})
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     optional.Value(500),
				})
			}))
		}
	})
}))
