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

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(mat.ModalOverlayColor),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("content", co.New(mat.Paper, func() {
			co.WithData(mat.PaperData{
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Width:            optional.Value(600),
				Height:           optional.Value(600),
				HorizontalCenter: optional.Value(0),
				VerticalCenter:   optional.Value(0),
			})

			co.WithChild("header", co.New(mat.Toolbar, func() {
				co.WithLayoutData(mat.LayoutData{
					Top:   optional.Value(0),
					Left:  optional.Value(0),
					Right: optional.Value(0),
				})

				co.WithChild("twod_texture", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/texture.png"),
						Text:     "2D Texture",
						Selected: lifecycle.SelectedKind() == "twod_texture",
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind("twod_texture")
						},
					})
				}))

				co.WithChild("cube_texture", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/texture.png"),
						Text:     "Cube Texture",
						Selected: lifecycle.SelectedKind() == "cube_texture",
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind("cube_texture")
						},
					})
				}))

				co.WithChild("model", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/model.png"),
						Text:     "Model",
						Selected: lifecycle.SelectedKind() == "model",
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind("model")
						},
					})
				}))

				co.WithChild("scene", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text:     "Scene",
						Icon:     co.OpenImage("resources/icons/scene.png"),
						Selected: lifecycle.SelectedKind() == "scene",
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.SetSelectedKind("scene")
						},
					})
				}))
			}))

			co.WithChild(fmt.Sprintf("content-%s", lifecycle.SelectedKind()), co.New(mat.ScrollPane, func() {
				co.WithData(mat.ScrollPaneData{
					DisableHorizontal: true,
				})

				co.WithLayoutData(mat.LayoutData{
					Left:   optional.Value(0),
					Right:  optional.Value(0),
					Top:    optional.Value(mat.ToolbarHeight), // FIXME: Use frame layout
					Bottom: optional.Value(mat.ToolbarHeight), // FIXME: Use frame layout
				})

				co.WithChild("content", co.New(mat.List, func() {
					co.WithLayoutData(mat.LayoutData{
						GrowHorizontally: true,
					})

					for _, resource := range lifecycle.Resources() {
						func(resource *data.Resource) {
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
									Selected:     resource.ID() == lifecycle.SelectedResourceID(),
								})
								co.WithLayoutData(mat.LayoutData{
									GrowHorizontally: true,
								})
								co.WithCallbackData(AssetItemCallbackData{
									OnSelected: func(id string) {
										lifecycle.SetSelectedResourceID(id)
									},
								})
							}))
						}(resource)
					}
				}))
			}))

			co.WithChild("footer", co.New(mat.Toolbar, func() {
				co.WithData(mat.ToolbarData{
					Orientation: mat.ToolbarOrientationRightToLeft,
					Positioning: mat.ToolbarPositioningBottom,
				})
				co.WithLayoutData(mat.LayoutData{
					Left:   optional.Value(0),
					Right:  optional.Value(0),
					Bottom: optional.Value(0),
				})

				co.WithChild("open", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text:     "Open",
						Disabled: lifecycle.SelectedResourceID() == "",
					})
					co.WithCallbackData(mat.ToolbarButtonCallbackData{
						OnClick: func() {
							lifecycle.OnOpen(lifecycle.SelectedResourceID())
						},
					})
				}))

				co.WithChild("cancel", co.New(mat.ToolbarButton, func() {
					co.WithData(mat.ToolbarButtonData{
						Text: "Cancel",
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
	handle          co.LifecycleHandle
	registry        *data.Registry
	onClose         func()
	onAssetSelected func(id string)
	selectedKind    string
	selectedAssetID string
}

func (l *assetDialogLifecycle) OnCreate(props co.Properties) {
	l.OnUpdate(props)
	l.selectedKind = "twod_texture"
	l.selectedAssetID = ""
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

func (l *assetDialogLifecycle) OnOpen(id string) {
	l.onAssetSelected(id)
	l.onClose()
}

func (l *assetDialogLifecycle) SelectedKind() string {
	return l.selectedKind
}

func (l *assetDialogLifecycle) SetSelectedKind(kind string) {
	l.selectedKind = kind
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) SetSelectedResourceID(id string) {
	l.selectedAssetID = id
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) SelectedResourceID() string {
	return l.selectedAssetID
}

func (l *assetDialogLifecycle) Resources() []*data.Resource {
	return l.registry.ListResourcesOfKind(l.selectedKind)
}

type AssetItemData struct {
	Selected     bool
	PreviewImage image.Image
	ID           string
	Kind         string
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
	assetKind    string
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

func (l *assetItemLifecycle) AssetKind() string {
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
