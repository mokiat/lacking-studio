package registry

import (
	"strings"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var CreateAssetModal = co.Define(&createAssetModalComponent{})

type CreateAssetModalCallbackData struct {
	OnApply func(name string)
}

type createAssetModalComponent struct {
	co.BaseComponent

	name string

	onApply func(name string)
}

func (c *createAssetModalComponent) OnCreate() {
	callbackData := co.GetCallbackData[CreateAssetModalCallbackData](c.Properties())
	c.onApply = callbackData.OnApply
}

func (c *createAssetModalComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(400),
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
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentSpacing: 30,
					}),
				})

				co.WithChild("info", co.New(std.Label, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(18)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      "Specify a name for your new scene.",
					})
				}))

				co.WithChild("settings", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ElementData{
						Padding: ui.UniformSpacing(10),
						Layout: layout.Frame(layout.FrameSettings{
							ContentSpacing: ui.Spacing{
								Left: 10,
							},
						}),
					})

					co.WithChild("label", co.New(std.Label, func() {
						co.WithLayoutData(layout.Data{
							HorizontalAlignment: layout.HorizontalAlignmentLeft,
						})
						co.WithData(std.LabelData{
							Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
							FontSize:  opt.V(float32(18)),
							FontColor: opt.V(std.OnSurfaceColor),
							Text:      "Name:",
						})
					}))

					co.WithChild("editbox", co.New(std.EditBox, func() {
						co.WithLayoutData(layout.Data{
							HorizontalAlignment: layout.HorizontalAlignmentCenter,
						})
						co.WithData(std.EditBoxData{
							Text: c.name,
						})
						co.WithCallbackData(std.EditBoxCallbackData{
							OnChange: func(text string) {
								c.setName(text)
							},
							OnSubmit: func(text string) {
								c.setName(text)
								c.onCreate()
							},
						})
					}))
				}))

			}))

			co.WithChild("footer", co.New(std.Toolbar, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment: layout.VerticalAlignmentBottom,
				})
				co.WithData(std.ToolbarData{
					Positioning: std.ToolbarPositioningBottom,
				})

				co.WithChild("create", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text:    "Create",
						Enabled: opt.V(strings.TrimSpace(c.name) != ""),
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: c.onCreate,
					})
				}))

				co.WithChild("cancel", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text: "Cancel",
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: c.onCancel,
					})
				}))
			}))
		}))
	})
}

func (c *createAssetModalComponent) setName(name string) {
	c.name = name
	c.Invalidate()
}

func (c *createAssetModalComponent) onCreate() {
	c.onApply(strings.TrimSpace(c.name))
	co.CloseOverlay(c.Scope())
}

func (c *createAssetModalComponent) onCancel() {
	co.CloseOverlay(c.Scope())
}
