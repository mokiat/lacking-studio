package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
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

var BinaryEditor = co.ContextScoped(co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data        = co.GetData[BinaryEditorData](props)
		editorModel = data.EditorModel
		controller  = data.EditorController
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
					controller.OnChangeContentFromPath(paths[0])
					return true
				},
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})

			co.WithChild("panel", co.New(mat.Container, func() {
				co.WithData(mat.ContainerData{
					BackgroundColor: opt.V(mat.BackgroundColor),
					Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
				})

				co.WithChild("icon", co.New(mat.Picture, func() {
					co.WithData(mat.PictureData{
						Image:      co.OpenImage(scope, "icons/upload.png"),
						ImageColor: opt.V(mat.SurfaceColor),
						Mode:       mat.ImageModeStretch,
					})
					co.WithLayoutData(mat.LayoutData{
						Width:            opt.V(48),
						Height:           opt.V(48),
						HorizontalCenter: opt.V(48),
						VerticalCenter:   opt.V(48),
					})
				}))
			}))
		}))

		if editorModel.IsPropertiesVisible() {
			co.WithChild("right", co.New(BinaryProperties, func() {
				co.WithData(BinaryPropertiesData{
					Model:            editorModel.Properties(),
					ResourceModel:    data.ResourceModel,
					BinaryModel:      data.BinaryModel,
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
}))
