package view

import (
	"fmt"
	"image"
	"log"

	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type AssetDialogData struct {
	Registry *data.Registry
}

type AssetDialogCallbackData struct {
	OnAssetSelected func(id string)
	OnClose         func()
}

var AssetDialog = co.Define(func(props co.Properties) co.Instance {
	var lifecycle *assetDialogLifecycle
	co.UseLifecycle(func(handle co.LifecycleHandle) co.Lifecycle {
		return &assetDialogLifecycle{
			Lifecycle: co.NewBaseLifecycle(),
			handle:    handle,
		}
	}, &lifecycle)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(ui.RGBA(0x00, 0x00, 0x00, 0xF0)),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("content", co.New(widget.Paper, func() {
			co.WithData(widget.PaperData{
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Width:            optional.NewInt(600),
				Height:           optional.NewInt(600),
				HorizontalCenter: optional.NewInt(0),
				VerticalCenter:   optional.NewInt(0),
			})

			co.WithChild("header", co.New(widget.Toolbar, func() {
				co.WithLayoutData(mat.LayoutData{
					Top:   optional.NewInt(0),
					Left:  optional.NewInt(0),
					Right: optional.NewInt(0),
				})

				co.WithChild("twod_texture", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/texture.png"),
						Text:     "2D Texture",
						Selected: lifecycle.SelectedKind() == "twod_texture",
					})
					co.WithCallbackData(widget.ToolbarButtonCallbackData{
						ClickListener: func() {
							lifecycle.SetSelectedKind("twod_texture")
						},
					})
				}))

				co.WithChild("cube_texture", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/texture.png"),
						Text:     "Cube Texture",
						Selected: lifecycle.SelectedKind() == "cube_texture",
					})
					co.WithCallbackData(widget.ToolbarButtonCallbackData{
						ClickListener: func() {
							lifecycle.SetSelectedKind("cube_texture")
						},
					})
				}))

				co.WithChild("model", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/model.png"),
						Text:     "Model",
						Selected: lifecycle.SelectedKind() == "model",
					})
					co.WithCallbackData(widget.ToolbarButtonCallbackData{
						ClickListener: func() {
							lifecycle.SetSelectedKind("model")
						},
					})
				}))

				co.WithChild("scene", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Text:     "Scene",
						Icon:     co.OpenImage("resources/icons/scene.png"),
						Selected: lifecycle.SelectedKind() == "scene",
					})
					co.WithCallbackData(widget.ToolbarButtonCallbackData{
						ClickListener: func() {
							lifecycle.SetSelectedKind("scene")
						},
					})
				}))
			}))

			co.WithChild("content", co.New(mat.Container, func() {
				co.WithData(mat.ContainerData{
					BackgroundColor: optional.NewColor(ui.RGB(240, 240, 240)),
					Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
						ContentAlignment: mat.AlignmentLeft,
					}),
				})
				co.WithLayoutData(mat.LayoutData{
					Left:   optional.NewInt(0),
					Right:  optional.NewInt(0),
					Top:    optional.NewInt(widget.ToolbarHeight),
					Bottom: optional.NewInt(widget.ToolbarHeight),
				})

				for _, resource := range lifecycle.Resources() {
					func(resource data.Resource) {
						co.WithChild(resource.GUID, co.New(AssetItem, func() {
							co.WithData(AssetItemData{
								PreviewImage: resource.PreviewImage,
								ID:           resource.GUID,
								Kind:         resource.Kind,
								Name:         resource.Name,
							})
							co.WithLayoutData(mat.LayoutData{
								GrowHorizontally: true,
							})
						}))
					}(resource)
				}
			}))

			co.WithChild("footer", co.New(widget.Toolbar, func() {
				co.WithData(widget.ToolbarData{
					Flipped: true,
				})
				co.WithLayoutData(mat.LayoutData{
					Left:   optional.NewInt(0),
					Right:  optional.NewInt(0),
					Bottom: optional.NewInt(0),
				})

				co.WithChild("open", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Text:     "Open",
						Disabled: true,
					})
				}))

				co.WithChild("cancel", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Text: "Cancel",
					})
					co.WithCallbackData(widget.ToolbarButtonCallbackData{
						ClickListener: func() {
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
	handle       co.LifecycleHandle
	registry     *data.Registry
	onClose      func()
	selectedKind string
}

func (l *assetDialogLifecycle) OnCreate(props co.Properties) {
	l.OnUpdate(props)
	l.selectedKind = "twod_texture"
}

func (l *assetDialogLifecycle) OnUpdate(props co.Properties) {
	var data AssetDialogData
	props.InjectData(&data)
	l.registry = data.Registry

	var callbackData AssetDialogCallbackData
	props.InjectOptionalCallbackData(&callbackData, AssetDialogCallbackData{})
	l.onClose = callbackData.OnClose
}

func (l *assetDialogLifecycle) OnCancel() {
	l.onClose()
}

func (l *assetDialogLifecycle) SelectedKind() string {
	return l.selectedKind
}

func (l *assetDialogLifecycle) SetSelectedKind(kind string) {
	l.selectedKind = kind
	l.handle.NotifyChanged()
}

func (l *assetDialogLifecycle) Resources() []data.Resource {
	return l.registry.ListResourcesOfKind(l.selectedKind)
}

type AssetItemData struct {
	Selected     bool
	PreviewImage image.Image
	ID           string
	Kind         string
	Name         string
}

var AssetItem = co.Define(func(props co.Properties) co.Instance {
	var lifecycle *assetItemLifecycle
	co.UseLifecycle(func(handle co.LifecycleHandle) co.Lifecycle {
		return &assetItemLifecycle{
			Lifecycle: co.NewBaseLifecycle(),
			handle:    handle,
		}
	}, &lifecycle)

	return co.New(widget.ListItem, func() {
		co.WithData(widget.ListItemData{
			Selected: lifecycle.IsSelected(),
		})
		co.WithCallbackData(widget.ListItemCallbackData{
			OnSelected: lifecycle.OnSelected,
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("item", co.New(widget.Paper, func() {
			co.WithData(widget.PaperData{
				Padding: ui.Spacing{
					Left:   5,
					Right:  5,
					Top:    5,
					Bottom: 5,
				},
				Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
					ContentAlignment: mat.AlignmentLeft,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("preview", co.New(mat.Picture, func() {
				co.WithData(mat.PictureData{
					Image:           lifecycle.PreviewImage(),
					BackgroundColor: optional.NewColor(ui.Black()),
					ImageColor:      optional.NewColor(ui.White()),
					Mode:            mat.ImageModeFit,
				})
				co.WithLayoutData(mat.LayoutData{
					Width:  optional.NewInt(64),
					Height: optional.NewInt(64),
				})
			}))

			co.WithChild("info", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
						ContentAlignment: mat.AlignmentLeft,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("id", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(16),
						FontColor: optional.NewColor(ui.Black()),
						Text:      fmt.Sprintf("ID: %s", lifecycle.AssetID()),
					})
				}))

				co.WithChild("kind", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(16),
						FontColor: optional.NewColor(ui.Black()),
						Text:      fmt.Sprintf("Kind: %s", lifecycle.AssetKind()),
					})
				}))

				co.WithChild("id", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(16),
						FontColor: optional.NewColor(ui.Black()),
						Text:      fmt.Sprintf("Name: %s", lifecycle.AssetName()),
					})
				}))
			}))
		}))
	})

})

type assetItemLifecycle struct {
	co.Lifecycle
	handle co.LifecycleHandle

	previewImage ui.Image
	assetID      string
	assetKind    string
	assetName    string
	selected     bool
}

func (l *assetItemLifecycle) OnCreate(props co.Properties) {
	var data AssetItemData
	props.InjectData(&data)

	l.previewImage = co.CreateImage(data.PreviewImage)
	l.assetID = data.ID
	l.assetKind = data.Kind
	l.assetName = data.Name
	l.selected = data.Selected
}

func (l *assetItemLifecycle) OnDestroy() {
	l.previewImage.Destroy()
	l.previewImage = nil
}

func (l *assetItemLifecycle) PreviewImage() ui.Image {
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
	log.Println("selected:", l.assetID)
}
