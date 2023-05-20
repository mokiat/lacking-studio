package view

import (
	"github.com/mokiat/gog/filter"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type StudioToolbarData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var StudioToolbar = co.Define(&studioToolbarComponent{})

type studioToolbarComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	studioModel   *model.Studio
	history       *model.History
	controller    StudioController
	assetsOverlay co.Overlay
}

func (c *studioToolbarComponent) OnUpsert() {
	data := co.GetData[StudioToolbarData](c.Properties)
	c.studioModel = data.StudioModel
	c.controller = data.StudioController

	mvc.UseBinding(c.studioModel, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeStudioEditorSelection)
	})

	c.history = c.studioModel.SelectedHistory()
	mvc.UseBinding(c.history, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeHistory)
	})

	mvc.UseBinding(c.studioModel.Registry(), filter.True[mvc.Change]())
}

func (c *studioToolbarComponent) Render() co.Instance {
	return co.New(std.Toolbar, func() {
		co.WithLayoutData(c.Properties.LayoutData())

		co.WithChild("assets", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope, "icons/assets.png"),
				Text: "Assets",
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onAssetsClicked,
			})
		}))

		co.WithChild("separator1", co.New(std.ToolbarSeparator, nil))

		co.WithChild("save", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope, "icons/save.png"),
				Enabled: opt.V(c.history.CanSave()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.controller.OnSave()
				},
			})
		}))

		co.WithChild("separator2", co.New(std.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope, "icons/undo.png"),
				Enabled: opt.V(c.history.CanUndo()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.controller.OnUndo()
				},
			})
		}))

		co.WithChild("redo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope, "icons/redo.png"),
				Enabled: opt.V(c.history.CanRedo()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.controller.OnRedo()
				},
			})
		}))

		co.WithChild("separator3", co.New(std.ToolbarSeparator, nil))

		co.WithChild("properties", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope, "icons/properties.png"),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.controller.OnToggleProperties()
				},
			})
		}))
	})
}

func (c *studioToolbarComponent) onAssetsClicked() {
	c.assetsOverlay = co.OpenOverlay(c.Scope, co.New(AssetDialog, func() {
		co.WithData(AssetDialogData{
			Registry:   c.studioModel.Registry(),
			Controller: c.controller,
		})
		co.WithCallbackData(AssetDialogCallbackData{
			OnOpen: func(id string) {
				c.controller.OnOpenResource(id)
			},
			OnClose: func() {
				overlay := c.assetsOverlay
				overlay.Close()
			},
		})
	}))
}
