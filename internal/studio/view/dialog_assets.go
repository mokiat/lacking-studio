package view

import (
	"fmt"
	"image"

	"github.com/mokiat/gog/filter"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

type AssetDialogData struct {
	Registry   *model.Registry
	Controller StudioController
}

type AssetDialogCallbackData struct {
	OnOpen  func(id string)
	OnClose func()
}

var defaultAssetDialogCallbackData = AssetDialogCallbackData{
	OnOpen:  func(string) {},
	OnClose: func() {},
}

var AssetDialog = co.Define(&assetDialogComponent{})

type assetDialogComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`
	Invalidate func()        `co:"invalidate"`

	controller       StudioController
	registry         *model.Registry
	onOpen           func(id string)
	onClose          func()
	searchText       string
	selectedKind     model.ResourceKind
	selectedResource *model.Resource
}

func (c *assetDialogComponent) OnCreate() {
	c.selectedKind = model.ResourceKindTwoDTexture
	c.selectedResource = nil
	c.searchText = ""
}

func (c *assetDialogComponent) OnUpsert() {
	data := co.GetData[AssetDialogData](c.Properties)
	c.registry = data.Registry
	c.controller = data.Controller

	callbackData := co.GetOptionalCallbackData(c.Properties, defaultAssetDialogCallbackData)
	c.onOpen = callbackData.OnOpen
	c.onClose = callbackData.OnClose
}

func (c *assetDialogComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(600),
			Height:           opt.V(600),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("header", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(std.ElementData{
				Padding: ui.Spacing{
					Bottom: 10,
				},
				Layout: layout.Vertical(layout.VerticalSettings{
					ContentSpacing:   10,
					ContentAlignment: layout.HorizontalAlignmentLeft,
				}),
			})

			co.WithChild("toolbar", co.New(std.Toolbar, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})

				co.WithChild("twod_texture", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope, "icons/texture.png"),
						Text:     "2D Texture",
						Selected: c.selectedKind == model.ResourceKindTwoDTexture,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setSelectedKind(model.ResourceKindTwoDTexture)
						},
					})
				}))

				co.WithChild("cube_texture", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope, "icons/texture.png"),
						Text:     "Cube Texture",
						Selected: c.selectedKind == model.ResourceKindCubeTexture,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setSelectedKind(model.ResourceKindCubeTexture)
						},
					})
				}))

				co.WithChild("model", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope, "icons/model.png"),
						Text:     "Model",
						Selected: c.selectedKind == model.ResourceKindModel,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setSelectedKind(model.ResourceKindModel)
						},
					})
				}))

				co.WithChild("scene", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text:     "Scene",
						Icon:     co.OpenImage(c.Scope, "icons/scene.png"),
						Selected: c.selectedKind == model.ResourceKindScene,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setSelectedKind(model.ResourceKindScene)
						},
					})
				}))

				co.WithChild("binary", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text:     "Binary",
						Icon:     co.OpenImage(c.Scope, "icons/broken-image.png"),
						Selected: c.selectedKind == model.ResourceKindBinary,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setSelectedKind(model.ResourceKindBinary)
						},
					})
				}))
			}))

			co.WithChild("search", co.New(std.Element, func() {
				co.WithData(std.ElementData{
					Layout: layout.Horizontal(layout.HorizontalSettings{
						ContentAlignment: layout.VerticalAlignmentCenter,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("label", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope, "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(18)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      "Search:",
					})
				}))

				co.WithChild("editbox", co.New(std.Editbox, func() {
					co.WithData(std.EditboxData{
						Text: c.searchText,
					})

					co.WithLayoutData(layout.Data{
						Width: opt.V(200),
					})

					co.WithCallbackData(std.EditboxCallbackData{
						OnChanged: func(text string) {
							c.setSearchText(text)
						},
					})
				}))

				co.WithChild("clear", co.New(std.Button, func() {
					co.WithData(std.ButtonData{
						Text: "Clear",
					})

					co.WithCallbackData(std.ButtonCallbackData{
						OnClick: func() {
							c.setSearchText("")
						},
					})
				}))
			}))
		}))

		co.WithChild(fmt.Sprintf("content-%s", c.selectedKind), co.New(std.ScrollPane, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ScrollPaneData{
				DisableHorizontal: true,
			})

			co.WithChild("content", co.New(std.List, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})

				c.eachResource(func(resource *model.Resource) {
					previewImage := resource.PreviewImage()
					co.WithChild(resource.ID(), co.New(AssetItem, func() {
						co.WithData(AssetItemData{
							PreviewImage: previewImage,
							ID:           resource.ID(),
							Kind:         resource.Kind(),
							Name:         resource.Name(),
							Selected:     resource == c.selectedResource,
						})
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithCallbackData(AssetItemCallbackData{
							OnSelected: func(id string) {
								c.setSelectedResource(resource)
							},
						})
					}))
				})
			}))
		}))

		co.WithChild("footer", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentBottom,
			})
			co.WithData(std.ElementData{
				Padding: ui.Spacing{
					Top: 10,
				},
				Layout: layout.Vertical(layout.VerticalSettings{
					ContentSpacing:   10,
					ContentAlignment: layout.HorizontalAlignmentLeft,
				}),
			})

			co.WithChild("actions", co.New(std.Element, func() {
				co.WithData(std.ElementData{
					Layout: layout.Horizontal(layout.HorizontalSettings{
						ContentSpacing:   5,
						ContentAlignment: layout.VerticalAlignmentCenter,
					}),
				})

				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})

				co.WithChild("delete", co.New(std.Button, func() {
					co.WithData(std.ButtonData{
						Icon:    co.OpenImage(c.Scope, "icons/delete.png"),
						Text:    "Delete",
						Enabled: opt.V(c.selectedResource != nil),
					})

					co.WithCallbackData(std.ButtonCallbackData{
						OnClick: func() {
							// TODO: Open confirm dialog
							c.handleDelete()
						},
					})
				}))

				co.WithChild("clone", co.New(std.Button, func() {
					co.WithData(std.ButtonData{
						Icon:    co.OpenImage(c.Scope, "icons/file-copy.png"),
						Text:    "Clone",
						Enabled: opt.V(c.selectedResource != nil),
					})

					co.WithCallbackData(std.ButtonCallbackData{
						OnClick: func() {
							c.handleClone()
						},
					})
				}))

				co.WithChild("spacing", co.New(std.Spacing, func() {
					co.WithData(std.SpacingData{
						Size: ui.NewSize(20, 0),
					})
				}))

				co.WithChild("new", co.New(std.Button, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithData(std.ButtonData{
						Icon: co.OpenImage(c.Scope, "icons/file-add.png"),
						Text: "New",
					})
					co.WithCallbackData(std.ButtonCallbackData{
						OnClick: func() {
							c.handleNew()
						},
					})
				}))
			}))

			co.WithChild("toolbar", co.New(std.Toolbar, func() {
				co.WithData(std.ToolbarData{
					Positioning: std.ToolbarPositioningBottom,
				})
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})

				co.WithChild("open", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text:    "Open",
						Enabled: opt.V(c.selectedResource != nil),
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.handleOpen(c.selectedResource)
						},
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
		}))
	})
}

func (c *assetDialogComponent) handleCancel() {
	c.onClose()
}

func (c *assetDialogComponent) handleOpen(resource *model.Resource) {
	c.onOpen(resource.ID())
	c.onClose()
}

func (c *assetDialogComponent) handleNew() {
	resource := c.controller.OnCreateResource(c.selectedKind)
	if resource != nil {
		c.searchText = resource.Name()
		c.selectedResource = resource
		c.Invalidate()
	}
}

func (c *assetDialogComponent) handleClone() {
	resource := c.controller.OnCloneResource(c.selectedResource.ID())
	if resource != nil {
		c.searchText = resource.Name()
		c.selectedResource = resource
		c.Invalidate()
	}
}

func (c *assetDialogComponent) handleDelete() {
	c.controller.OnDeleteResource(c.selectedResource.ID())
	c.selectedResource = nil
	c.Invalidate()
}

func (c *assetDialogComponent) setSelectedKind(kind model.ResourceKind) {
	c.selectedKind = kind
	c.selectedResource = nil
	c.searchText = ""
	c.Invalidate()
}

func (c *assetDialogComponent) setSearchText(text string) {
	c.searchText = text
	c.selectedResource = nil
	c.Invalidate()
}

func (c *assetDialogComponent) setSelectedResource(resource *model.Resource) {
	c.selectedResource = resource
	c.Invalidate()
}

func (c *assetDialogComponent) eachResource(fn func(*model.Resource)) {
	fltrs := []filter.Func[*model.Resource]{
		model.ResourcesWithKind(c.selectedKind),
	}
	if c.searchText != "" {
		fltrs = append(fltrs, model.ResourcesWithSimilarName(c.searchText))
	}
	c.registry.IterateResources(fn, fltrs...)
}

type AssetItemData struct {
	Selected     bool
	PreviewImage image.Image
	ID           string
	Kind         model.ResourceKind
	Name         string
}

type AssetItemCallbackData struct {
	OnSelected func(id string)
}

var defaultAssetItemCallbackData = AssetItemCallbackData{
	OnSelected: func(id string) {},
}

var AssetItem = co.Define(&assetItemComponent{})

type assetItemComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`
	Invalidate func()        `co:"invalidate"`

	previewImage        *ui.Image
	defaultPreviewImage *ui.Image
	assetID             string
	assetKind           model.ResourceKind
	assetName           string
	selected            bool
	onSelected          func(id string)
}

func (c *assetItemComponent) OnUpsert() {
	data := co.GetData[AssetItemData](c.Properties)
	callbackData := co.GetOptionalCallbackData(c.Properties, defaultAssetItemCallbackData)

	if c.previewImage != nil {
		c.previewImage.Destroy()
	}
	if data.PreviewImage != nil {
		c.previewImage = co.CreateImage(c.Scope, data.PreviewImage)
	}
	c.defaultPreviewImage = co.OpenImage(c.Scope, "icons/broken-image.png")
	c.assetID = data.ID
	c.assetKind = data.Kind
	c.assetName = data.Name
	c.selected = data.Selected
	c.onSelected = callbackData.OnSelected
}

func (c *assetItemComponent) OnDelete() {
	if c.previewImage != nil {
		c.previewImage.Destroy()
	}
	c.previewImage = nil
}

func (c *assetItemComponent) Render() co.Instance {
	previewImage := c.previewImage
	if previewImage == nil {
		previewImage = c.defaultPreviewImage
	}

	return co.New(std.ListItem, func() {
		co.WithLayoutData(c.Properties.LayoutData())
		co.WithData(std.ListItemData{
			Selected: c.selected,
		})
		co.WithCallbackData(std.ListItemCallbackData{
			OnSelected: func() {
				c.onSelected(c.assetID)
			},
		})

		co.WithChild("item", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("preview", co.New(std.Picture, func() {
				co.WithData(std.PictureData{
					Image:           previewImage,
					BackgroundColor: opt.V(ui.Black()),
					ImageColor:      opt.V(ui.White()),
					Mode:            std.ImageModeFit,
				})
				co.WithLayoutData(layout.Data{
					Width:  opt.V(64),
					Height: opt.V(64),
				})
			}))

			co.WithChild("info", co.New(std.Element, func() {
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentAlignment: layout.HorizontalAlignmentLeft,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("name", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope, "ui:///roboto-bold.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      c.assetName,
					})
				}))

				co.WithChild("id", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope, "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      c.assetID,
					})
				}))
			}))
		}))
	})
}
