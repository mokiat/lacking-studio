package view

import (
	"fmt"
	"image"

	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type AssetDialogData struct {
	Registry *data.Registry
}

type AssetDialogCallbackData struct {
	OnAssetSelected func(id string)
	OnClose         func()
}

var defaultAssetDialogCallbackData = AssetDialogCallbackData{
	OnAssetSelected: func(id string) {},
	OnClose:         func() {},
}

var AssetDialog = co.Define(func(props co.Properties) co.Instance {
	lifecycle := co.UseLifecycle(func(handle co.LifecycleHandle) *assetDialogLifecycle {
		return &assetDialogLifecycle{
			Lifecycle: co.NewBaseLifecycle(),
			handle:    handle,
		}
	})

	return co.New(mat.Modal, func() {
		co.WithLayoutData(mat.LayoutData{
			Width:            optional.Value(600),
			Height:           optional.Value(600),
			HorizontalCenter: optional.Value(0),
			VerticalCenter:   optional.Value(0),
		})

		co.WithChild("header", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Padding: ui.Spacing{
					Bottom: 10,
				},
				Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
					ContentSpacing:   10,
					ContentAlignment: mat.AlignmentLeft,
				}),
			})

			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentTop,
			})

			co.WithChild("toolbar", co.New(mat.Toolbar, func() {
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})

				co.WithChild("twod_texture", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/texture.png"),
						Text:     "2D Texture",
						Selected: lifecycle.SelectedKind() == data.ResourceKindTwoDTexture,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(data.ResourceKindTwoDTexture)
						},
					})
				}))

				co.WithChild("cube_texture", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/texture.png"),
						Text:     "Cube Texture",
						Selected: lifecycle.SelectedKind() == data.ResourceKindCubeTexture,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(data.ResourceKindCubeTexture)
						},
					})
				}))

				co.WithChild("model", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/model.png"),
						Text:     "Model",
						Selected: lifecycle.SelectedKind() == data.ResourceKindModel,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(data.ResourceKindModel)
						},
					})
				}))

				co.WithChild("scene", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text:     "Scene",
						Icon:     co.OpenImage("resources/icons/scene.png"),
						Selected: lifecycle.SelectedKind() == data.ResourceKindScene,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(data.ResourceKindScene)
						},
					})
				}))
			}))

			co.WithChild("search", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.OpenFont("mat:///roboto-regular.ttf"),
						FontSize:  optional.Value(float32(18)),
						FontColor: optional.Value(mat.OnSurfaceColor),
						Text:      "Search:",
					})
				}))

				co.WithChild("editbox", co.New(mat.Editbox, func() {
					co.WithData(mat.EditboxData{
						Text: lifecycle.SearchText(),
					})

					co.WithLayoutData(mat.LayoutData{
						Width: optional.Value(200),
					})

					co.WithCallbackData(mat.EditboxCallbackData{
						OnChanged: func(text string) {
							lifecycle.SetSearchText(text)
						},
					})
				}))

				co.WithChild("clear", co.New(mat.Button, func() {
					co.WithData(mat.ButtonData{
						Text: "Clear",
					})

					co.WithCallbackData(mat.ButtonCallbackData{
						ClickListener: func() {
							lifecycle.SetSearchText("")
						},
					})
				}))
			}))
		}))

		co.WithChild(fmt.Sprintf("content-%s", lifecycle.SelectedKind()), co.New(mat.ScrollPane, func() {
			co.WithData(mat.ScrollPaneData{
				DisableHorizontal: true,
			})

			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})

			co.WithChild("content", co.New(mat.List, func() {
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})

				lifecycle.EachResource(func(resource *data.Resource) {
					previewImage, err := resource.LoadPreview()
					if err != nil {
						previewImage = nil
					}
					co.WithChild(resource.ID(), co.New(AssetItem, func() {
						co.WithData(AssetItemData{
							PreviewImage: previewImage,
							ID:           resource.ID(),
							Kind:         resource.Kind(),
							Name:         resource.Name(),
							Selected:     resource == lifecycle.SelectedResource(),
						})
						co.WithLayoutData(mat.LayoutData{
							GrowHorizontally: true,
						})
						co.WithCallbackData(AssetItemCallbackData{
							OnSelected: func(id string) {
								lifecycle.SetSelectedResource(resource)
							},
						})
					}))
				})
			}))
		}))

		co.WithChild("footer", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Padding: ui.Spacing{
					Top: 10,
				},
				Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
					ContentSpacing:   10,
					ContentAlignment: mat.AlignmentLeft,
				}),
			})

			co.WithLayoutData(mat.LayoutData{
				Alignment:        mat.AlignmentBottom,
				GrowHorizontally: true,
			})

			co.WithChild("actions", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentSpacing:   5,
						ContentAlignment: mat.AlignmentCenter,
					}),
				})

				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})

				co.WithChild("delete", co.New(mat.Button, func() {
					co.WithData(mat.ButtonData{
						Icon:    co.OpenImage("resources/icons/delete.png"),
						Text:    "Delete",
						Enabled: optional.Value(lifecycle.SelectedResource() != nil),
					})

					co.WithCallbackData(mat.ButtonCallbackData{
						ClickListener: func() {
							// TODO: Open confirm dialog
							lifecycle.OnDelete()
						},
					})
				}))

				co.WithChild("clone", co.New(mat.Button, func() {
					co.WithData(mat.ButtonData{
						Icon:    co.OpenImage("resources/icons/file-copy.png"),
						Text:    "Clone",
						Enabled: optional.Value(lifecycle.SelectedResource() != nil),
					})

					co.WithCallbackData(mat.ButtonCallbackData{
						ClickListener: func() {
							lifecycle.OnClone()
						},
					})
				}))

				co.WithChild("spacing", co.New(mat.Spacing, func() {
					co.WithData(mat.SpacingData{
						Width: 20,
					})
				}))

				co.WithChild("new", co.New(mat.Button, func() {
					co.WithData(mat.ButtonData{
						Icon: co.OpenImage("resources/icons/file-add.png"),
						Text: "New",
					})

					co.WithLayoutData(mat.LayoutData{
						Alignment: mat.AlignmentRight,
					})

					co.WithCallbackData(mat.ButtonCallbackData{
						ClickListener: func() {
							lifecycle.OnNew()
						},
					})
				}))
			}))

			co.WithChild("toolbar", co.New(mat.Toolbar, func() {
				co.WithData(mat.ToolbarData{
					Orientation: mat.ToolbarOrientationRightToLeft,
					Positioning: mat.ToolbarPositioningBottom,
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})

				co.WithChild("open", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text:    "Open",
						Enabled: optional.Value(lifecycle.SelectedResource() != nil),
					})
					co.WithLayoutData(mat.LayoutData{
						Alignment: mat.AlignmentRight,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.OnOpen(lifecycle.SelectedResource())
						},
					})
				}))

				co.WithChild("cancel", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text: "Cancel",
					})
					co.WithLayoutData(mat.LayoutData{
						Alignment: mat.AlignmentRight,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.OnCancel()
						},
					})
				}))
			}))
		}))
	})
})

type assetDialogLifecycle struct {
	co.Lifecycle
	handle           co.LifecycleHandle
	registry         *data.Registry
	onClose          func()
	onAssetSelected  func(id string)
	searchText       string
	selectedKind     data.ResourceKind
	selectedResource *data.Resource
}

func (l *assetDialogLifecycle) OnCreate(props co.Properties) {
	l.OnUpdate(props)
	l.selectedKind = data.ResourceKindTwoDTexture
	l.selectedResource = nil
	l.searchText = ""
}

func (l *assetDialogLifecycle) OnUpdate(props co.Properties) {
	var (
		data         = co.GetData[AssetDialogData](props)
		callbackData = co.GetOptionalCallbackData(props, AssetDialogCallbackData{})
	)

	l.registry = data.Registry
	l.onClose = callbackData.OnClose
	l.onAssetSelected = callbackData.OnAssetSelected
}

func (l *assetDialogLifecycle) OnCancel() {
	l.onClose()
}

func (l *assetDialogLifecycle) OnOpen(resource *data.Resource) {
	l.onAssetSelected(resource.ID())
	l.onClose()
}

func (l *assetDialogLifecycle) SelectedKind() data.ResourceKind {
	return l.selectedKind
}

func (l *assetDialogLifecycle) SetSelectedKind(kind data.ResourceKind) {
	l.selectedKind = kind
	l.selectedResource = nil
	l.searchText = ""
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) SetSelectedResource(resource *data.Resource) {
	l.selectedResource = resource
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) SelectedResource() *data.Resource {
	return l.selectedResource
}

func (l *assetDialogLifecycle) SearchText() string {
	return l.searchText
}

func (l *assetDialogLifecycle) SetSearchText(text string) {
	l.searchText = text
	l.selectedResource = nil
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) EachResource(fn func(*data.Resource)) {
	filter := data.FilterWithKind(l.selectedKind)
	if l.searchText != "" {
		filter = data.FilterAnd(
			filter,
			data.FilterWithSimilarName(l.searchText),
		)
	}
	l.registry.EachResource(filter, fn)
}

func (l *assetDialogLifecycle) OnNew() {
	resource, err := l.registry.NewResource(l.selectedKind)
	if err != nil {
		panic(err)
	}
	l.searchText = resource.Name()
	l.selectedResource = resource
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) OnClone() {
	resource, err := l.selectedResource.Clone()
	if err != nil {
		panic(err)
	}
	l.searchText = resource.Name()
	l.selectedResource = resource
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) OnDelete() {
	if err := l.selectedResource.Delete(); err != nil {
		panic(err)
	}
	l.selectedResource = nil
	l.handle.NotifyChanged()
}

type AssetItemData struct {
	Selected     bool
	PreviewImage image.Image
	ID           string
	Kind         data.ResourceKind
	Name         string
}

type AssetItemCallbackData struct {
	OnSelected func(id string)
}

var defaultAssetItemCallbackData = AssetItemCallbackData{
	OnSelected: func(id string) {},
}

var AssetItem = co.Define(func(props co.Properties) co.Instance {
	lifecycle := co.UseLifecycle(func(handle co.LifecycleHandle) *assetItemLifecycle {
		return &assetItemLifecycle{
			Lifecycle: co.NewBaseLifecycle(),
		}
	})

	return co.New(mat.ListItem, func() {
		co.WithData(mat.ListItemData{
			Selected: lifecycle.IsSelected(),
		})
		co.WithCallbackData(mat.ListItemCallbackData{
			OnSelected: lifecycle.OnSelected,
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("item", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
					ContentAlignment: mat.AlignmentLeft,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("preview", co.New(mat.Picture, func() {
				co.WithData(mat.PictureData{
					Image:           lifecycle.PreviewImage(),
					BackgroundColor: optional.Value(ui.Black()),
					ImageColor:      optional.Value(ui.White()),
					Mode:            mat.ImageModeFit,
				})
				co.WithLayoutData(mat.LayoutData{
					Width:  optional.Value(64),
					Height: optional.Value(64),
				})
			}))

			co.WithChild("info", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
						ContentAlignment: mat.AlignmentLeft,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("name", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.OpenFont("mat:///roboto-bold.ttf"),
						FontSize:  optional.Value(float32(16)),
						FontColor: optional.Value(ui.Black()),
						Text:      lifecycle.AssetName(),
					})
				}))

				co.WithChild("id", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.OpenFont("mat:///roboto-regular.ttf"),
						FontSize:  optional.Value(float32(16)),
						FontColor: optional.Value(ui.Black()),
						Text:      lifecycle.AssetID(),
					})
				}))
			}))
		}))
	})
})

type assetItemLifecycle struct {
	co.Lifecycle

	previewImage *ui.Image
	assetID      string
	assetKind    data.ResourceKind
	assetName    string
	selected     bool
	onSelected   func(id string)
}

func (l *assetItemLifecycle) OnCreate(props co.Properties) {
	l.OnUpdate(props)
}

func (l *assetItemLifecycle) OnUpdate(props co.Properties) {
	var (
		data         = co.GetData[AssetItemData](props)
		callbackData = co.GetOptionalCallbackData(props, defaultAssetItemCallbackData)
	)

	if l.previewImage != nil {
		l.previewImage.Destroy()
	}
	if data.PreviewImage != nil {
		l.previewImage = co.CreateImage(data.PreviewImage)
	}
	l.assetID = data.ID
	l.assetKind = data.Kind
	l.assetName = data.Name
	l.selected = data.Selected
	l.onSelected = callbackData.OnSelected
}

func (l *assetItemLifecycle) OnDestroy() {
	if l.previewImage != nil {
		l.previewImage.Destroy()
	}
	l.previewImage = nil
}

func (l *assetItemLifecycle) PreviewImage() *ui.Image {
	return l.previewImage
}

func (l *assetItemLifecycle) AssetID() string {
	return l.assetID
}

func (l *assetItemLifecycle) AssetKind() data.ResourceKind {
	return l.assetKind
}

func (l *assetItemLifecycle) AssetName() string {
	return l.assetName
}

func (l *assetItemLifecycle) IsSelected() bool {
	return l.selected
}

func (l *assetItemLifecycle) OnSelected() {
	l.onSelected(l.assetID)
}
