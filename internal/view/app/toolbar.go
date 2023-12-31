package app

import (
	"fmt"

	appmodel "github.com/mokiat/lacking-studio/internal/model/app"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	registrymodel "github.com/mokiat/lacking-studio/internal/model/registry"
	registryview "github.com/mokiat/lacking-studio/internal/view/registry"
	"github.com/mokiat/lacking/debug/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Toolbar = mvc.EventListener(co.Define(&toolbarComponent{}))

type ToolbarData struct {
	AppModel      *appmodel.Model
	RegistryModel *registrymodel.Model
}

type toolbarComponent struct {
	co.BaseComponent

	appModel      *appmodel.Model
	registryModel *registrymodel.Model
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
	case appmodel.ActiveEditorChangedEvent:
		c.Invalidate()
	}
}

func (c *toolbarComponent) onNewClicked() {
	var name string
	for i := 1; ; i++ {
		name = fmt.Sprintf("Untitled-%d", i)
		if !c.appModel.HasEditorWithName(name) {
			break
		}
	}

	eventBus := co.TypedValue[*mvc.EventBus](c.Scope())
	editor := editormodel.NewModel(eventBus, name)
	c.appModel.AddEditor(editor)
	c.appModel.SetActiveEditor(editor)
}

func (c *toolbarComponent) onBrowseClicked() {
	co.OpenOverlay(c.Scope(), co.New(registryview.Modal, func() {
		co.WithData(registryview.ModalData{
			RegistryModel: c.registryModel,
		})
		co.WithCallbackData(registryview.ModalCallbackData{
			OnOpen: c.onAssetOpen,
		})
	}))
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

func (c *toolbarComponent) onAssetOpen(asset *registrymodel.Asset) {
	log.Info("OPEN")
}
