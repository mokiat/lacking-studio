package app

import (
	appmodel "github.com/mokiat/lacking-studio/internal/model/app"
	"github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking-studio/internal/view/common"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Tabbar = mvc.EventListener(co.Define(&tabbarComponent{}))

type TabbarData struct {
	AppModel *appmodel.Model
}

type tabbarComponent struct {
	co.BaseComponent

	appModel *appmodel.Model
}

func (c *tabbarComponent) OnUpsert() {
	data := co.GetData[TabbarData](c.Properties())
	c.appModel = data.AppModel
}

func (c *tabbarComponent) Render() co.Instance {
	return co.New(std.Tabbar, func() {
		co.WithLayoutData(c.Properties().LayoutData())

		c.appModel.EachEditor(func(editor *editor.Model) {
			co.WithChild(editor.ID(), co.New(std.TabbarTab, func() {
				co.WithData(std.TabbarTabData{
					Icon:     c.editorIcon(editor),
					Text:     c.editorTitle(editor),
					Selected: c.appModel.ActiveEditor() == editor,
				})
				co.WithCallbackData(std.TabbarTabCallbackData{
					OnClick: func() {
						c.selectEditor(editor)
					},
					OnClose: func() {
						c.closeEditor(editor, false)
					},
				})
			}))
		})
	})
}

func (c *tabbarComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case appmodel.EditorsChangedEvent:
		c.Invalidate()
	case appmodel.ActiveEditorChangedEvent:
		c.Invalidate()
	}
}

func (c *tabbarComponent) editorIcon(editor *editor.Model) *ui.Image {
	return editor.Image()
}

func (c *tabbarComponent) editorTitle(editor *editor.Model) string {
	text := editor.Name()
	if editor.CanSave() {
		text += " *"
	}
	return text
}

func (c *tabbarComponent) selectEditor(editor *editor.Model) {
	c.appModel.SetActiveEditor(editor)
}

func (c *tabbarComponent) closeEditor(editor *editor.Model, force bool) {
	if !force && editor.CanSave() {
		common.OpenConfirmation(c.Scope(),
			"There are unsaved changes!\n\nAre you sure you want to continue?",
			func() { c.closeEditor(editor, true) },
		)
	} else {
		c.appModel.RemoveEditor(editor)
	}
}
