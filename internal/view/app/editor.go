package app

import (
	"github.com/mokiat/gog/opt"
	appmodel "github.com/mokiat/lacking-studio/internal/model/app"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking-studio/internal/view/common"
	editorview "github.com/mokiat/lacking-studio/internal/view/editor"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Editor = mvc.EventListener(co.Define(&editorComponent{}))

type EditorData struct {
	AppModel    *appmodel.Model
	EditorModel *editormodel.Model
	Visible     bool
}

var _ ui.ElementStateHandler = (*editorComponent)(nil)
var _ ui.ElementHistoryHandler = (*editorComponent)(nil)

type editorComponent struct {
	co.BaseComponent

	appModel    *appmodel.Model
	editorModel *editormodel.Model
	visible     bool
}

func (c *editorComponent) OnUpsert() {
	data := co.GetData[EditorData](c.Properties())
	c.appModel = data.AppModel
	c.editorModel = data.EditorModel
	c.visible = data.Visible
}

func (c *editorComponent) Render() co.Instance {
	// TODO: Use horizontal searator container so that navigator
	// and inspector can be resized by user.

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout:    layout.Frame(),
			Essence:   c,
			Visible:   opt.V(c.visible),
			Focusable: opt.V(true),
		})

		if c.appModel.IsNavigatorVisible() {
			co.WithChild("navigator", co.New(editorview.Navigator, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					Width:               opt.V(300),
				})
				co.WithData(editorview.NavigatorData{
					EditorModel: c.editorModel,
				})
			}))
		}

		co.WithChild("workbench", co.New(editorview.Workbench, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
			})
			co.WithData(editorview.WorkbenchData{
				EditorModel: c.editorModel,
			})
		}))

		if c.appModel.IsInspectorVisible() {
			co.WithChild("inspector", co.New(editorview.Inspector, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(300),
				})
				co.WithData(editorview.InspectorData{
					EditorModel: c.editorModel,
				})
			}))
		}
	})
}

func (c *editorComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case appmodel.NavigatorVisibleChangedEvent:
		c.Invalidate()
	case appmodel.InspectorVisibleChangedEvent:
		c.Invalidate()
	}
}

func (c *editorComponent) OnSave(element *ui.Element) bool {
	// Needs to be performed outside of the render pass, otherwise
	// it leads to broken command states / references.
	co.Schedule(c.Scope(), func() {
		if err := c.editorModel.Save(c.Scope()); err != nil {
			log.Error("Failed to save: %v", err.Error())
			common.OpenError(c.Scope(), "Error saving scene!\nCheck logs for more info.")
		}
	})
	return true
}

func (c *editorComponent) OnUndo(element *ui.Element) bool {
	c.editorModel.Undo()
	return true
}

func (c *editorComponent) OnRedo(element *ui.Element) bool {
	c.editorModel.Redo()
	return true
}
