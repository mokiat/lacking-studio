package registry

import (
	"github.com/mokiat/gog/opt"
	registrymodel "github.com/mokiat/lacking-studio/internal/model/registry"
	"github.com/mokiat/lacking/debug/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Modal = co.Define(&modalComponent{})

type ModalData struct {
	RegistryModel *registrymodel.Model
}

type ModalCallbackData struct {
	OnOpen func(asset *registrymodel.Asset)
}

type modalComponent struct {
	co.BaseComponent

	registryModel *registrymodel.Model
	selectedAsset *registrymodel.Asset

	onOpen func(asset *registrymodel.Asset)
}

func (c *modalComponent) OnUpsert() {
	data := co.GetData[ModalData](c.Properties())
	c.registryModel = data.RegistryModel

	callbackData := co.GetCallbackData[ModalCallbackData](c.Properties())
	c.onOpen = callbackData.OnOpen
}

func (c *modalComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(600),
			Height:           opt.V(600),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("header", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})

			co.WithChild("clone", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Icon:    co.OpenImage(c.Scope(), "icons/file-copy.png"),
					Text:    "Clone",
					Enabled: opt.V(c.selectedAsset != nil),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleClone,
				})
			}))

			co.WithChild("delete", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Icon:    co.OpenImage(c.Scope(), "icons/delete.png"),
					Text:    "Delete",
					Enabled: opt.V(c.selectedAsset != nil),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleDelete,
				})
			}))

			co.WithChild("separator-before-search", co.New(std.ToolbarSeparator, nil))

			co.WithChild("search", co.New(std.Editbox, func() {
				co.WithData(std.EditboxData{
					Text: "Search text....",
				})

				co.WithLayoutData(layout.Data{
					Width: opt.V(200),
				})

				co.WithCallbackData(std.EditboxCallbackData{
					OnChanged: func(text string) {
						// c.setSearchText(text)
					},
				})
			}))
		}))

		co.WithChild("content", co.New(std.ScrollPane, func() {
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

				c.eachAsset(func(asset *registrymodel.Asset) {
					co.WithChild(asset.ID(), co.New(Item, func() {
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithData(ItemData{
							Asset:    asset,
							Selected: c.selectedAsset == asset,
						})
						co.WithCallbackData(ItemCallbackData{
							OnSelected: c.handleItemSelected,
						})
					}))
				})
			}))
		}))

		co.WithChild("footer", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentBottom,
			})
			co.WithData(std.ToolbarData{
				Positioning: std.ToolbarPositioningBottom,
			})

			co.WithChild("open", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text:    "Open",
					Enabled: opt.V(c.selectedAsset != nil),
				})
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleOpen,
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

func (c *modalComponent) handleClone() {
	log.Info("CLONE!")
}

func (c *modalComponent) handleDelete() {
	log.Info("DELETE")
}

func (c *modalComponent) eachAsset(cb func(asset *registrymodel.Asset)) {
	for _, asset := range c.registryModel.Assets() {
		cb(asset)
	}
}

func (c *modalComponent) handleItemSelected(asset *registrymodel.Asset) {
	c.selectedAsset = asset
	c.Invalidate()
}

func (c *modalComponent) handleOpen() {
	co.CloseOverlay(c.Scope())
	c.onOpen(c.selectedAsset)
}

func (c *modalComponent) handleCancel() {
	co.CloseOverlay(c.Scope())
}
