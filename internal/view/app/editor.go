package app

import (
	"github.com/mokiat/gog/opt"
	editorview "github.com/mokiat/lacking-studio/internal/view/editor"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Editor = co.Define(&editorComponent{})

type EditorData struct {
	Visible bool
}

type editorComponent struct {
	co.BaseComponent

	visible bool
}

func (c *editorComponent) OnUpsert() {
	data := co.GetData[EditorData](c.Properties())
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

		co.WithChild("navigator", co.New(editorview.Navigator, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentLeft,
				Width:               opt.V(300),
			})
			co.WithData(editorview.NavigatorData{
				// TODO
			})
		}))

		co.WithChild("workbench", co.New(editorview.Workbench, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
			})
			co.WithData(editorview.WorkbenchData{
				// TODO
			})
		}))

		co.WithChild("inspector", co.New(editorview.Inspector, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
				Width:               opt.V(300),
			})
			co.WithData(editorview.InspectorData{
				// TODO
			})
		}))

	})
}
