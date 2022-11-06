package view

import (
	"fmt"
	"image"

	"github.com/mokiat/gog/filter"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
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

var AssetDialog = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	lifecycle := co.UseLifecycle(func(handle co.LifecycleHandle) *assetDialogLifecycle {
		return &assetDialogLifecycle{
			Lifecycle: co.NewBaseLifecycle(),
			handle:    handle,
		}
	})

	return co.New(mat.Modal, func() {
		co.WithLayoutData(mat.LayoutData{
			Width:            opt.V(600),
			Height:           opt.V(600),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
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
						Icon:     co.OpenImage(scope, "icons/texture.png"),
						Text:     "2D Texture",
						Selected: lifecycle.SelectedKind() == model.ResourceKindTwoDTexture,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(model.ResourceKindTwoDTexture)
						},
					})
				}))

				co.WithChild("cube_texture", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage(scope, "icons/texture.png"),
						Text:     "Cube Texture",
						Selected: lifecycle.SelectedKind() == model.ResourceKindCubeTexture,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(model.ResourceKindCubeTexture)
						},
					})
				}))

				co.WithChild("model", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage(scope, "icons/model.png"),
						Text:     "Model",
						Selected: lifecycle.SelectedKind() == model.ResourceKindModel,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(model.ResourceKindModel)
						},
					})
				}))

				co.WithChild("scene", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text:     "Scene",
						Icon:     co.OpenImage(scope, "icons/scene.png"),
						Selected: lifecycle.SelectedKind() == model.ResourceKindScene,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(model.ResourceKindScene)
						},
					})
				}))

				co.WithChild("binary", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text:     "Binary",
						Icon:     co.OpenImage(scope, "icons/broken-image.png"),
						Selected: lifecycle.SelectedKind() == model.ResourceKindBinary,
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind(model.ResourceKindBinary)
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
						Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(18)),
						FontColor: opt.V(mat.OnSurfaceColor),
						Text:      "Search:",
					})
				}))

				co.WithChild("editbox", co.New(mat.Editbox, func() {
					co.WithData(mat.EditboxData{
						Text: lifecycle.SearchText(),
					})

					co.WithLayoutData(mat.LayoutData{
						Width: opt.V(200),
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

				lifecycle.EachResource(func(resource *model.Resource) {
					previewImage := resource.PreviewImage()
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
						Icon:    co.OpenImage(scope, "icons/delete.png"),
						Text:    "Delete",
						Enabled: opt.V(lifecycle.SelectedResource() != nil),
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
						Icon:    co.OpenImage(scope, "icons/file-copy.png"),
						Text:    "Clone",
						Enabled: opt.V(lifecycle.SelectedResource() != nil),
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
						Icon: co.OpenImage(scope, "icons/file-add.png"),
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
						Enabled: opt.V(lifecycle.SelectedResource() != nil),
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
	controller       StudioController
	registry         *model.Registry
	onOpen           func(id string)
	onClose          func()
	searchText       string
	selectedKind     model.ResourceKind
	selectedResource *model.Resource
}

func (l *assetDialogLifecycle) OnCreate(props co.Properties, scope co.Scope) {
	l.OnUpdate(props, scope)
	l.selectedKind = model.ResourceKindTwoDTexture
	l.selectedResource = nil
	l.searchText = ""
}

func (l *assetDialogLifecycle) OnUpdate(props co.Properties, scope co.Scope) {
	var (
		data         = co.GetData[AssetDialogData](props)
		callbackData = co.GetOptionalCallbackData(props, AssetDialogCallbackData{})
	)

	l.registry = data.Registry
	l.controller = data.Controller
	l.onOpen = callbackData.OnOpen
	l.onClose = callbackData.OnClose
}

func (l *assetDialogLifecycle) OnCancel() {
	l.onClose()
}

func (l *assetDialogLifecycle) OnOpen(resource *model.Resource) {
	l.onOpen(resource.ID())
	l.onClose()
}

func (l *assetDialogLifecycle) SelectedKind() model.ResourceKind {
	return l.selectedKind
}

func (l *assetDialogLifecycle) SetSelectedKind(kind model.ResourceKind) {
	l.selectedKind = kind
	l.selectedResource = nil
	l.searchText = ""
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) SetSelectedResource(resource *model.Resource) {
	l.selectedResource = resource
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) SelectedResource() *model.Resource {
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

func (l *assetDialogLifecycle) EachResource(fn func(*model.Resource)) {
	fltrs := []filter.Func[*model.Resource]{
		model.ResourcesWithKind(l.selectedKind),
	}
	if l.searchText != "" {
		fltrs = append(fltrs, model.ResourcesWithSimilarName(l.searchText))
	}
	l.registry.IterateResources(fn, fltrs...)
}

func (l *assetDialogLifecycle) OnNew() {
	resource := l.controller.OnCreateResource(l.selectedKind)
	if resource != nil {
		l.searchText = resource.Name()
		l.selectedResource = resource
		l.handle.NotifyChanged()
	}
}

func (l *assetDialogLifecycle) OnClone() {
	resource := l.controller.OnCloneResource(l.selectedResource.ID())
	if resource != nil {
		l.searchText = resource.Name()
		l.selectedResource = resource
		l.handle.NotifyChanged()
	}
}

func (l *assetDialogLifecycle) OnDelete() {
	l.controller.OnDeleteResource(l.selectedResource.ID())
	l.selectedResource = nil
	l.handle.NotifyChanged()
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

var AssetItem = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
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
					BackgroundColor: opt.V(ui.Black()),
					ImageColor:      opt.V(ui.White()),
					Mode:            mat.ImageModeFit,
				})
				co.WithLayoutData(mat.LayoutData{
					Width:  opt.V(64),
					Height: opt.V(64),
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
						Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      lifecycle.AssetName(),
					})
				}))

				co.WithChild("id", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      lifecycle.AssetID(),
					})
				}))
			}))
		}))
	})
})

type assetItemLifecycle struct {
	co.Lifecycle

	previewImage        *ui.Image
	defaultPreviewImage *ui.Image
	assetID             string
	assetKind           model.ResourceKind
	assetName           string
	selected            bool
	onSelected          func(id string)
}

func (l *assetItemLifecycle) OnCreate(props co.Properties, scope co.Scope) {
	l.OnUpdate(props, scope)
}

func (l *assetItemLifecycle) OnUpdate(props co.Properties, scope co.Scope) {
	var (
		data         = co.GetData[AssetItemData](props)
		callbackData = co.GetOptionalCallbackData(props, defaultAssetItemCallbackData)
	)

	if l.previewImage != nil {
		l.previewImage.Destroy()
	}
	if data.PreviewImage != nil {
		l.previewImage = co.CreateImage(scope, data.PreviewImage)
	}
	l.defaultPreviewImage = co.OpenImage(scope, "icons/broken-image.png")
	l.assetID = data.ID
	l.assetKind = data.Kind
	l.assetName = data.Name
	l.selected = data.Selected
	l.onSelected = callbackData.OnSelected
}

func (l *assetItemLifecycle) OnDestroy(scope co.Scope) {
	if l.previewImage != nil {
		l.previewImage.Destroy()
	}
	l.previewImage = nil
}

func (l *assetItemLifecycle) PreviewImage() *ui.Image {
	if l.previewImage == nil {
		return l.defaultPreviewImage
	}
	return l.previewImage
}

func (l *assetItemLifecycle) AssetID() string {
	return l.assetID
}

func (l *assetItemLifecycle) AssetKind() model.ResourceKind {
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
