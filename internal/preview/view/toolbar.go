package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/preview/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Toolbar = mvc.EventListener(co.Define(&toolbarComponent{}))

type ToolbarData struct {
	AppModel *model.AppModel
}

type toolbarComponent struct {
	co.BaseComponent

	appModel *model.AppModel
}

func (c *toolbarComponent) OnUpsert() {
	data := co.GetData[ToolbarData](c.Properties())
	c.appModel = data.AppModel
}

func (c *toolbarComponent) Render() co.Instance {
	return co.New(std.Toolbar, func() {
		co.WithLayoutData(c.Properties().LayoutData())

		co.WithChild("refresh", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "icons/refresh.png"),
				Text:    "Refresh",
				Enabled: opt.V(c.appModel.RefreshEnabled()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.handleRefresh,
			})
		}))

		// The following are listed in reverse.

		co.WithChild("quit", co.New(std.ToolbarButton, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
			})
			co.WithData(std.ToolbarButtonData{
				Icon: co.OpenImage(c.Scope(), "icons/quit.png"),
				Text: "Quit",
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.handleQuit,
			})
		}))

		co.WithChild("separator-between-quit-back", co.New(std.ToolbarSeparator, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
			})
		}))

		co.WithChild("back", co.New(std.ToolbarButton, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
			})
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "icons/back.png"),
				Text:    "Back",
				Enabled: opt.V(c.appModel.SelectedResource() != nil),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.handleBack,
			})
		}))
	})
}

func (c *toolbarComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.RefreshEvent:
		// TODO: Hide loading indicator
		c.Invalidate()
	case model.RefreshErrorEvent:
		// TODO: Hide loading indicator and show error
		c.Invalidate()
	case model.SelectedResourceChangedEvent:
		c.Invalidate()
	}
}

func (c *toolbarComponent) handleRefresh() {
	c.appModel.Refresh()
	// TODO: Show loading indicator
	c.Invalidate()
}

func (c *toolbarComponent) handleBack() {
	c.appModel.SetSelectedResource(nil)
	c.Invalidate()
}

func (c *toolbarComponent) handleQuit() {
	co.Window(c.Scope()).Close()
}
