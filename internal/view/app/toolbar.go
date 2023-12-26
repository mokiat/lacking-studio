package app

import (
	appmodel "github.com/mokiat/lacking-studio/internal/model/app"
	"github.com/mokiat/lacking/debug/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Toolbar = mvc.EventListener(co.Define(&toolbarComponent{}))

type ToolbarData struct {
	AppModel *appmodel.Model
}

type toolbarComponent struct {
	co.BaseComponent

	appModel *appmodel.Model
}

func (c *toolbarComponent) OnUpsert() {
	data := co.GetData[ToolbarData](c.Properties())
	c.appModel = data.AppModel
}

func (c *toolbarComponent) Render() co.Instance {
	return co.New(std.Toolbar, func() {
		co.WithLayoutData(c.Properties().LayoutData())

		co.WithChild("new", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope(), "icons/new.png"),
				Text: "New",
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onNewClicked,
			})
		}))

		co.WithChild("open", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope(), "icons/open.png"),
				Text: "Browse",
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onBrowseClicked,
			})
		}))

		co.WithChild("separator-after-assets", co.New(std.ToolbarSeparator, nil))

		co.WithChild("save", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope(), "icons/save.png"),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onSaveClicked,
			})
		}))

		co.WithChild("separator-after-save", co.New(std.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope(), "icons/undo.png"),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onUndoClicked,
			})
		}))

		co.WithChild("redo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope(), "icons/redo.png"),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onRedoClicked,
			})
		}))

		co.WithChild("separator-after-history", co.New(std.ToolbarSeparator, nil))

		if c.appModel.IsInspectorVisible() {
			co.WithChild("collapse-right", co.New(std.ToolbarButton, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithData(std.ToolbarButtonData{
					Icon: co.OpenImage(c.Scope(), "icons/collapse-right.png"),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.onCollapseInspector,
				})
			}))
		} else {
			co.WithChild("expand-right", co.New(std.ToolbarButton, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithData(std.ToolbarButtonData{
					Icon: co.OpenImage(c.Scope(), "icons/expand-right.png"),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.onExpandInspector,
				})
			}))
		}

		if c.appModel.IsNavigatorVisible() {
			co.WithChild("collapse-left", co.New(std.ToolbarButton, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithData(std.ToolbarButtonData{
					Icon: co.OpenImage(c.Scope(), "icons/collapse-left.png"),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.onCollapseNavigator,
				})
			}))
		} else {
			co.WithChild("expand-left", co.New(std.ToolbarButton, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithData(std.ToolbarButtonData{
					Icon: co.OpenImage(c.Scope(), "icons/expand-left.png"),
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.onExpandNavigator,
				})
			}))
		}

		co.WithChild("separator-before-expand-collapse", co.New(std.ToolbarSeparator, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
			})
		}))
	})
}

func (c *toolbarComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case appmodel.InspectorVisibleChangedEvent:
		c.Invalidate()
	case appmodel.NavigatorVisibleChangedEvent:
		c.Invalidate()
	}
}

func (c *toolbarComponent) onNewClicked() {
	log.Info("New")
}

func (c *toolbarComponent) onBrowseClicked() {
	log.Info("Browse")
	// c.assetsOverlay = co.OpenOverlay(c.Scope(), co.New(AssetDialog, func() {
	// 	co.WithData(AssetDialogData{
	// 		Registry:   c.studioModel.Registry(),
	// 		Controller: c.controller,
	// 	})
	// 	co.WithCallbackData(AssetDialogCallbackData{
	// 		OnOpen: func(id string) {
	// 			c.controller.OnOpenResource(id)
	// 		},
	// 		OnClose: func() {
	// 			overlay := c.assetsOverlay
	// 			overlay.Close()
	// 		},
	// 	})
	// }))
}

func (c *toolbarComponent) onSaveClicked() {
	log.Info("Save")
}

func (c *toolbarComponent) onUndoClicked() {
	log.Info("Undo")
}

func (c *toolbarComponent) onRedoClicked() {
	log.Info("Redo")
}

func (c *toolbarComponent) onExpandNavigator() {
	c.appModel.SetNavigatorVisible(true)
}

func (c *toolbarComponent) onCollapseNavigator() {
	c.appModel.SetNavigatorVisible(false)
}

func (c *toolbarComponent) onExpandInspector() {
	c.appModel.SetInpsectorVisible(true)
}

func (c *toolbarComponent) onCollapseInspector() {
	c.appModel.SetInpsectorVisible(false)
}
