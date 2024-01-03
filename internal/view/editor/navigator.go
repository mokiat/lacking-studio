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

var Navigator = mvc.EventListener(co.Define(&navigatorComponent{}))

type NavigatorData struct {
	EditorModel *editormodel.Model
}

type navigatorComponent struct {
	co.BaseComponent

	editorModel *editormodel.Model
}

func (c *navigatorComponent) OnUpsert() {
	data := co.GetData[NavigatorData](c.Properties())
	c.editorModel = data.EditorModel
}

func (c *navigatorComponent) Render() co.Instance {
	pageItems := []std.DropdownItem{
		{
			Key:   editormodel.NavigatorPageNodes,
			Label: "Nodes",
		},
		{
			Key:   editormodel.NavigatorPageTextures,
			Label: "Textures",
		},
		{
			Key:   editormodel.NavigatorPageMaterials,
			Label: "Materials",
		},
		{
			Key:   editormodel.NavigatorPageMeshes,
			Label: "Meshes",
		},
		{
			Key:   editormodel.NavigatorPageAnimations,
			Label: "Animations",
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
				SelectedKey: c.editorModel.NavigatorPage(),
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
				BackgroundColor: opt.V(ui.Red()),
			})
		}))
	})
}

func (c *navigatorComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case editormodel.NavigatorPageChangedEvent:
		if event.Editor == c.editorModel {
			c.Invalidate()
		}
	}
}

func (c *navigatorComponent) handlePageSelected(key any) {
	c.editorModel.SetNavigatorPage(key.(editormodel.NavigatorPage))
}
