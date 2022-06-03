package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

type StudioTabbarData struct {
	StudioModel *model.Studio
}

var StudioTabbar = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data   = co.GetData[StudioTabbarData](props)
		studio = data.StudioModel
	)

	mvc.UseBinding(studio, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeStudioEditorAdded) ||
			mvc.IsChange(ch, model.ChangeStudioEditorRemoved) ||
			mvc.IsChange(ch, model.ChangeStudioEditorSelection)
	})

	return co.New(mat.Tabbar, func() {
		co.WithLayoutData(props.LayoutData())

		studio.IterateEditors(func(editor *model.Editor) {
			key := editor.Resource().ID()
			co.WithChild(key, co.New(StudioTab, func() {
				co.WithData(StudioTabData{
					EditorModel: editor,
					Selected:    editor == studio.SelectedEditor(),
				})
			}))
		})
	})
})
