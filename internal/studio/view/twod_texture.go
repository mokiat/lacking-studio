package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

var TwoDTexture = co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.TwoDTextureEditor)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(widget.BackgroundColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("center", co.New(widget.DropZone, func() {
			co.WithCallbackData(widget.DropZoneCallbackData{
				OnDrop: func(paths []string) {
					editor.ChangeSourcePath(paths[0])
				},
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})

			co.WithChild("viewport", co.New(widget.Viewport, func() {
				co.WithData(widget.ViewportData{
					API:    editor.API(),
					Scene:  editor.Scene(),
					Camera: editor.Camera(),
				})
				co.WithCallbackData(widget.ViewportCallbackData{
					OnUpdate:     editor.Update,
					OnMouseEvent: editor.OnViewportMouseEvent,
				})
			}))
		}))

		if editor.IsPropertiesVisible() {
			co.WithChild("right", co.New(TwoDTextureProperties, func() {
				co.WithData(editor)
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     optional.Value(500),
				})
			}))
		}
	})
})
