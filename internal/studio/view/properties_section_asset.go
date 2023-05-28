package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type EditorController interface {
	OnRenameResource(name string)
}

type AssetPropertiesSectionData struct {
	Model            *model.Resource
	StudioController StudioController
	EditorController EditorController
}

var AssetPropertiesSection = mvc.Wrap(co.Define(&assetPropertiesSectionComponent{}))

type assetPropertiesSectionComponent struct {
	co.BaseComponent

	resource         *model.Resource
	controller       StudioController
	editorController EditorController
}

func (c *assetPropertiesSectionComponent) OnUpsert() {
	data := co.GetData[AssetPropertiesSectionData](c.Properties())
	c.resource = data.Model
	c.controller = data.StudioController
	c.editorController = data.EditorController

	mvc.UseBinding(c.Scope(), c.resource, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeResourceName)
	})
}

func (c *assetPropertiesSectionComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})
		co.WithData(std.ElementData{
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    5,
				Bottom: 5,
			},
			Layout: layout.Vertical(layout.VerticalSettings{
				ContentAlignment: layout.HorizontalAlignmentLeft,
				ContentSpacing:   5,
			}),
		})

		// TODO: Use GridLayout (2 columns)

		co.WithChild("id", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("label", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
					FontSize:  opt.V(float32(18)),
					FontColor: opt.V(std.OnSurfaceColor),
					Text:      "ID:",
				})
			}))

			co.WithChild("value", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
					FontSize:  opt.V(float32(18)),
					FontColor: opt.V(std.OnSurfaceColor),
					Text:      c.resource.ID(),
				})
			}))
		}))

		co.WithChild("type", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("label", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
					FontSize:  opt.V(float32(18)),
					FontColor: opt.V(std.OnSurfaceColor),
					Text:      "Type:",
				})
			}))

			co.WithChild("value", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
					FontSize:  opt.V(float32(18)),
					FontColor: opt.V(std.OnSurfaceColor),
					Text:      string(c.resource.Kind()),
				})
			}))
		}))

		co.WithChild("name", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("label", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
					FontSize:  opt.V(float32(18)),
					FontColor: opt.V(std.OnSurfaceColor),
					Text:      "Name:",
				})
			}))

			co.WithChild("value", co.New(std.Editbox, func() {
				co.WithData(std.EditboxData{
					Text: c.resource.Name(),
				})
				co.WithLayoutData(layout.Data{
					Width: opt.V(300),
				})
				co.WithCallbackData(std.EditboxCallbackData{
					OnChanged: func(text string) {
						c.editorController.OnRenameResource(text)
					},
				})
			}))
		}))

		co.WithChild("actions", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("delete", co.New(std.Button, func() {
				co.WithData(std.ButtonData{
					Icon: co.OpenImage(c.Scope(), "icons/delete.png"),
					Text: "Delete",
				})

				co.WithCallbackData(std.ButtonCallbackData{
					OnClick: func() {
						c.controller.OnDeleteResource(c.resource.ID())
					},
				})
			}))

			co.WithChild("clone", co.New(std.Button, func() {
				co.WithData(std.ButtonData{
					Icon: co.OpenImage(c.Scope(), "icons/file-copy.png"),
					Text: "Clone",
				})

				co.WithCallbackData(std.ButtonCallbackData{
					OnClick: func() {
						newResource := c.controller.OnCloneResource(c.resource.ID())
						if newResource != nil {
							c.controller.OnOpenResource(newResource.ID())
						}
					},
				})
			}))
		}))
	})
}
