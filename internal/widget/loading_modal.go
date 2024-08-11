package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var LoadingModal = co.Define(&loadingModalComponent{})

type loadingModalComponent struct {
	co.BaseComponent

	icon *ui.Image
}

func (c *loadingModalComponent) OnCreate() {
	c.icon = co.OpenImage(c.Scope(), "icons/info.png")
}

func (c *loadingModalComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(500),
			Height:           opt.V(200),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("dialog", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Frame(layout.FrameSettings{
					ContentSpacing: ui.SymmetricSpacing(0, 20),
				}),
			})

			co.WithChild("content", co.New(std.Element, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment: layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ElementData{
					Layout: layout.Frame(layout.FrameSettings{
						ContentSpacing: ui.Spacing{
							Left: 5,
						},
					}),
				})

				co.WithChild("icon", co.New(std.Picture, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentTop,
						Width:             opt.V(48),
						Height:            opt.V(48),
					})
					co.WithData(std.PictureData{
						Image:      c.icon,
						ImageColor: opt.V(ui.Black()),
						Mode:       std.ImageModeFit,
					})
				}))

				co.WithChild("text", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(20)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      "Loading...",
					})
				}))
			}))
		}))
	})
}
