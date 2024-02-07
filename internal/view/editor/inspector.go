package editor

import (
	"github.com/mokiat/gog/opt"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Inspector = mvc.EventListener(co.Define(&inspectorComponent{}))

type InspectorData struct {
	EditorModel *editormodel.Model
}

type inspectorComponent struct {
	co.BaseComponent

	editorModel *editormodel.Model
}

func (c *inspectorComponent) OnUpsert() {
	data := co.GetData[InspectorData](c.Properties())
	c.editorModel = data.EditorModel
}

func (c *inspectorComponent) Render() co.Instance {
	pageItems := []std.DropdownItem{
		{
			Key:   editormodel.InspectorPageAsset,
			Label: "Asset",
		},
		{
			Key:   editormodel.InspectorPageSelection,
			Label: "Selection",
		},
		{
			Key:   editormodel.InspectorPageViewport,
			Label: "Viewport",
		},
	}

	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderColor:     opt.V(std.OutlineColor),
			BorderSize:      ui.UniformSpacing(1),
			BackgroundColor: opt.V(std.SurfaceColor),
			Padding:         ui.UniformSpacing(2),
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.UniformSpacing(2),
			}),
		})

		co.WithChild("selector", co.New(std.Dropdown, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(std.DropdownData{
				Items:       pageItems,
				SelectedKey: c.editorModel.InspectorPage(),
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: c.handlePageSelected,
			})
		}))

		// TODO: Switch based on selected page
		co.WithChild("content", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ContainerData{
				BackgroundColor: opt.V(std.SurfaceColor),
			})
		}))
	})
}

func (c *inspectorComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case editormodel.InspectorPageChangedEvent:
		if event.Editor == c.editorModel {
			c.Invalidate()
		}
	}
}

func (c *inspectorComponent) handlePageSelected(key any) {
	c.editorModel.SetInspectorPage(key.(editormodel.InspectorPage))
}
