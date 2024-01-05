package editor

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var ModelImport = co.Define(&modelImportComponent{})

type ModelImportData struct {
	Model *pack.Model
}

type ModelImportCallbackData struct {
	OnImport func(model *pack.Model)
}

type modelImportComponent struct {
	co.BaseComponent

	model *pack.Model

	onImport func(model *pack.Model)
}

func (c *modelImportComponent) OnCreate() {
	data := co.GetData[ModelImportData](c.Properties())
	c.model = data.Model

	callbackData := co.GetCallbackData[ModelImportCallbackData](c.Properties())
	c.onImport = callbackData.OnImport
	if c.onImport == nil {
		c.onImport = func(model *pack.Model) {}
	}
}

func (c *modelImportComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(800),
			Height:           opt.V(600),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("header", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})

			co.WithChild("select-all", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Select All",
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleSelectAll,
				})
			}))

			co.WithChild("separator", co.New(std.ToolbarSeparator, nil))

			co.WithChild("deselect-all", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Deselect All",
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleDeselectAll,
				})
			}))
		}))

		co.WithChild("article", co.New(std.ScrollPane, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ScrollPaneData{
				DisableHorizontal: true,
			})

			co.WithChild("content", co.New(std.Element, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentSpacing: 2,
					}),
					Padding: ui.UniformSpacing(2),
				})

				co.WithChild("nodes", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Nodes",
						Expanded: true, // TODO
					})

					co.WithChild("first", co.New(std.Checkbox, func() {
						co.WithData(std.CheckboxData{
							Checked: true,
							Label:   "hello.png",
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: func(active bool) {
								log.Info("Toggle to %t", active)
							},
						})
					}))

					co.WithChild("second", co.New(std.Checkbox, func() {
						co.WithData(std.CheckboxData{
							Checked: false,
							Label:   "world.png",
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: func(active bool) {
								log.Info("Toggle to %t", active)
							},
						})
					}))
				}))

				co.WithChild("textures", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Textures",
						Expanded: true, // TODO
					})
				}))

				co.WithChild("materials", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Materials",
						Expanded: true, // TODO
					})
				}))

				co.WithChild("meshes", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Meshes",
						Expanded: true, // TODO
					})
				}))

				co.WithChild("animations", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Animations",
						Expanded: true, // TODO
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

			co.WithChild("import", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Import",
				})
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleImport,
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
					OnClick: c.handleCancel,
				})
			}))
		}))
	})
}

func (c *modelImportComponent) handleSelectAll() {

}

func (c *modelImportComponent) handleDeselectAll() {

}

func (c *modelImportComponent) handleImport() {
	co.CloseOverlay(c.Scope())
	c.onImport(c.model)
}

func (c *modelImportComponent) handleCancel() {
	co.CloseOverlay(c.Scope())
}
