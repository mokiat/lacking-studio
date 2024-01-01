package app

import (
	"github.com/mokiat/gog/opt"
	appmodel "github.com/mokiat/lacking-studio/internal/model/app"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	editorview "github.com/mokiat/lacking-studio/internal/view/editor"
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

type editorComponent struct {
	co.BaseComponent

	appModel *appmodel.Model
	visible  bool
}

func (c *editorComponent) OnUpsert() {
	data := co.GetData[EditorData](c.Properties())
	c.appModel = data.AppModel
	c.visible = data.Visible
}

func (c *editorComponent) Render() co.Instance {
	// TODO: Use horizontal searator container so that navigator
	// and inspector can be resized by user.

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout:  layout.Frame(),
			Visible: opt.V(c.visible),
		})

		if c.appModel.IsNavigatorVisible() {
			co.WithChild("navigator", co.New(editorview.Navigator, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					Width:               opt.V(300),
				})
				co.WithData(editorview.NavigatorData{
					// TODO
				})
			}))
		}

		co.WithChild("workbench", co.New(editorview.Workbench, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
			})
			co.WithData(editorview.WorkbenchData{
				// TODO
			})
		}))

		if c.appModel.IsInspectorVisible() {
			co.WithChild("inspector", co.New(editorview.Inspector, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(300),
				})
				co.WithData(editorview.InspectorData{
					// TODO
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
