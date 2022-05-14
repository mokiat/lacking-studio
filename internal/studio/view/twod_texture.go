package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

var TwoDTexture = co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.TwoDTextureEditor)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("center", co.New(mat.DropZone, func() {
			co.WithCallbackData(mat.DropZoneCallbackData{
				OnDrop: func(paths []string) bool {
					editor.ChangeSourcePath(paths[0])
					return true
				},
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})

			co.WithChild("viewport", co.New(mat.Viewport, func() {
				co.WithData(mat.ViewportData{
					API: editor.API(),
				})
				co.WithCallbackData(mat.ViewportCallbackData{
					OnMouseEvent: editor.OnViewportMouseEvent,
					OnRender: func(framebuffer render.Framebuffer, size ui.Size) {
						editor.Update()
						editor.Scene().RenderFramebuffer(framebuffer, graphics.Viewport{
							X:      0,
							Y:      0,
							Width:  size.Width,
							Height: size.Height,
						}, editor.Camera())
					},
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
