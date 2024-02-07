package registry

import (
	"strings"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var RenameAssetModal = co.Define(&renameAssetModalComponent{})

type RenameAssetModalData struct {
	OldName string
}

type RenameAssetModalCallbackData struct {
	OnApply func(name string)
}

type renameAssetModalComponent struct {
	co.BaseComponent

	oldName string
	newName string

	onApply func(name string)
}

func (c *renameAssetModalComponent) OnCreate() {
	data := co.GetData[RenameAssetModalData](c.Properties())
	c.oldName = data.OldName
	c.newName = data.OldName

	callbackData := co.GetCallbackData[RenameAssetModalCallbackData](c.Properties())
	c.onApply = callbackData.OnApply
}

func (c *renameAssetModalComponent) Render() co.Instance {
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
						Text:      "Specify a new name for your scene.",
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
							Text: c.newName,
						})
						co.WithCallbackData(std.EditBoxCallbackData{
							OnChange: func(text string) {
								c.setName(text)
							},
							OnSubmit: func(text string) {
								c.setName(text)
								c.onRename()
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

				co.WithChild("rename", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text:    "Rename",
						Enabled: opt.V(strings.TrimSpace(c.newName) != "" && c.newName != c.oldName),
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: c.onRename,
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

func (c *renameAssetModalComponent) setName(name string) {
	c.newName = name
	c.Invalidate()
}

func (c *renameAssetModalComponent) onRename() {
	c.onApply(strings.TrimSpace(c.newName))
	co.CloseOverlay(c.Scope())
}

func (c *renameAssetModalComponent) onCancel() {
	co.CloseOverlay(c.Scope())
}
