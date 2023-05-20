package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type StudioTabbarData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var StudioTabbar = co.Define(&studioTabbarComponent{})

type studioTabbarComponent struct {
	Properties co.Properties `co:"properties"`

	studio     *model.Studio
	controller StudioController
}

func (c *studioTabbarComponent) OnUpsert() {
	data := co.GetData[StudioTabbarData](c.Properties)
	c.studio = data.StudioModel
	c.controller = data.StudioController

	mvc.UseBinding(c.studio, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeStudioEditorAdded) ||
			mvc.IsChange(ch, model.ChangeStudioEditorRemoved) ||
			mvc.IsChange(ch, model.ChangeStudioEditorSelection)
	})
}

func (c *studioTabbarComponent) Render() co.Instance {
	return co.New(std.Tabbar, func() {
		co.WithLayoutData(c.Properties.LayoutData())

		c.studio.IterateEditors(func(editor *model.Editor) {
			key := editor.Resource().ID()
			co.WithChild(key, co.New(StudioTab, func() {
				co.WithData(StudioTabData{
					EditorModel:      editor,
					StudioController: c.controller,
					Selected:         editor == c.studio.SelectedEditor(),
				})
			}))
		})
	})
}
