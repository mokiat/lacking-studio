package registry

import (
	"fmt"
	"strings"

	"github.com/mokiat/gog/opt"
	registrymodel "github.com/mokiat/lacking-studio/internal/model/registry"
	"github.com/mokiat/lacking-studio/internal/view/common"
	"github.com/mokiat/lacking/debug/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var BrowseAssetsModal = mvc.EventListener(co.Define(&browseAssetsModalComponent{}))

type BrowseAssetsModalData struct {
	RegistryModel *registrymodel.Model
}

type BrowseAssetsModalCallbackData struct {
	OnOpen func(asset *registrymodel.Asset)
}

type browseAssetsModalComponent struct {
	co.BaseComponent

	registryModel *registrymodel.Model

	selectedAsset *registrymodel.Asset
	searchText    string

	onOpen func(asset *registrymodel.Asset)
}

func (c *browseAssetsModalComponent) OnUpsert() {
	data := co.GetData[BrowseAssetsModalData](c.Properties())
	c.registryModel = data.RegistryModel

	callbackData := co.GetCallbackData[BrowseAssetsModalCallbackData](c.Properties())
	c.onOpen = callbackData.OnOpen

	if c.selectedAsset != nil && !c.assetMatchesSearch(c.selectedAsset) {
		c.selectedAsset = nil
	}
}

func (c *browseAssetsModalComponent) Render() co.Instance {
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

			co.WithChild("rename", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Icon:    co.OpenImage(c.Scope(), "icons/edit.png"),
					Text:    "Rename",
					Enabled: opt.V(c.selectedAsset != nil),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleRename,
				})
			}))

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

			co.WithChild("search", co.New(std.EditBox, func() {
				co.WithData(std.EditBoxData{
					Text: c.searchText,
				})

				co.WithLayoutData(layout.Data{
					Width: opt.V(200),
				})

				co.WithCallbackData(std.EditBoxCallbackData{
					OnChange: c.handleSearchTextChange,
					OnReject: c.handleSearchReject,
				})
			}))
		}))

		co.WithChild("content", co.New(std.ScrollPane, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ScrollPaneData{
				DisableHorizontal: true,
				Focused:           false,
			})

			co.WithChild("list", co.New(std.List, func() {
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

func (c *browseAssetsModalComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case registrymodel.AssetsChangedEvent:
		c.Invalidate()
	}
}

func (c *browseAssetsModalComponent) handleRename() {
	co.OpenOverlay(c.Scope(), co.New(RenameAssetModal, func() {
		co.WithData(RenameAssetModalData{
			OldName: c.selectedAsset.Name(),
		})
		co.WithCallbackData(RenameAssetModalCallbackData{
			OnApply: c.handleRenameApply,
		})
	}))
}

func (c *browseAssetsModalComponent) handleRenameApply(newName string) {
	if err := c.registryModel.RenameAsset(c.selectedAsset, newName); err != nil {
		log.Error("Failed to rename scene: %s", err.Error())
		common.OpenError(c.Scope(), "Error renaming scene!")
		return
	}
}

func (c *browseAssetsModalComponent) handleClone() {
	newAsset, err := c.registryModel.CloneAsset(c.selectedAsset)
	if err != nil {
		log.Error("Failed to clone scene: %s", err.Error())
		common.OpenError(c.Scope(), "Error cloning scene!")
		return
	}
	c.searchText = newAsset.Name()
	c.selectedAsset = newAsset
	c.Invalidate()
}

func (c *browseAssetsModalComponent) handleDelete() {
	message := fmt.Sprintf("Will delete scene:\n%q\n\nAre you sure?", c.selectedAsset.Name())
	common.OpenConfirmation(c.Scope(), message, c.handleDeleteConfirm)
}

func (c *browseAssetsModalComponent) handleDeleteConfirm() {
	defer func() {
		c.selectedAsset = nil
	}()
	if err := c.registryModel.DeleteAsset(c.selectedAsset); err != nil {
		log.Error("Failed to delete scene: %s", err.Error())
		common.OpenError(c.Scope(), "Error deleting scene!")
		return
	}
}

func (c *browseAssetsModalComponent) handleSearchTextChange(text string) {
	c.searchText = text
	c.Invalidate()
}

func (c *browseAssetsModalComponent) handleSearchReject() {
	c.searchText = ""
	c.Invalidate()
}

func (c *browseAssetsModalComponent) eachAsset(cb func(asset *registrymodel.Asset)) {
	for _, asset := range c.registryModel.Assets() {
		if c.assetMatchesSearch(asset) {
			cb(asset)
		}
	}
}

func (c *browseAssetsModalComponent) assetMatchesSearch(asset *registrymodel.Asset) bool {
	if c.searchText == "" {
		return true
	}
	return strings.Contains(asset.ID(), c.searchText) || strings.Contains(asset.Name(), c.searchText)
}

func (c *browseAssetsModalComponent) handleItemSelected(asset *registrymodel.Asset) {
	c.selectedAsset = asset
	c.Invalidate()
}

func (c *browseAssetsModalComponent) handleOpen() {
	co.CloseOverlay(c.Scope())
	c.onOpen(c.selectedAsset)
}

func (c *browseAssetsModalComponent) handleCancel() {
	co.CloseOverlay(c.Scope())
}
