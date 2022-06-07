package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

type StudioTabData struct {
	EditorModel      *model.Editor
	StudioController StudioController
	Selected         bool
}

var StudioTab = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data       = co.GetData[StudioTabData](props)
		editor     = data.EditorModel
		controller = data.StudioController
		resource   = editor.Resource()
	)

	mvc.UseBinding(resource, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeResourceName)
	})

	iconForModelKind := func(kind model.ResourceKind) *ui.Image {
		switch kind {
		case model.ResourceKindTwoDTexture:
			return co.OpenImage(scope, "icons/texture.png")
		case model.ResourceKindCubeTexture:
			return co.OpenImage(scope, "icons/texture.png")
		case model.ResourceKindModel:
			return co.OpenImage(scope, "icons/model.png")
		case model.ResourceKindScene:
			return co.OpenImage(scope, "icons/scene.png")
		default:
			return co.OpenImage(scope, "icons/broken-image.png")
		}
	}

	return co.New(mat.TabbarTab, func() {
		co.WithData(mat.TabbarTabData{
			Icon:     iconForModelKind(resource.Kind()),
			Text:     resource.Name(),
			Selected: data.Selected,
		})
		co.WithCallbackData(mat.TabbarTabCallbackData{
			OnClick: func() {
				controller.OnSelectEditor(editor)
			},
			OnClose: func() {
				controller.OnCloseEditor(editor)
			},
		})
	})
})
