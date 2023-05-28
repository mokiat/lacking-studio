package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type StudioTabData struct {
	EditorModel      *model.Editor
	StudioController StudioController
	Selected         bool
}

var StudioTab = mvc.Wrap(co.Define(&studioTabComponent{}))

type studioTabComponent struct {
	co.BaseComponent

	editor     *model.Editor
	controller StudioController
	resource   *model.Resource
	isSelected bool
}

func (c *studioTabComponent) OnUpsert() {
	data := co.GetData[StudioTabData](c.Properties())
	c.editor = data.EditorModel
	c.controller = data.StudioController
	c.resource = c.editor.Resource()
	c.isSelected = data.Selected

	mvc.UseBinding(c.Scope(), c.resource, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeResourceName)
	})
}

func (c *studioTabComponent) Render() co.Instance {
	iconForModelKind := func(kind model.ResourceKind) *ui.Image {
		switch kind {
		case model.ResourceKindTwoDTexture:
			return co.OpenImage(c.Scope(), "icons/texture.png")
		case model.ResourceKindCubeTexture:
			return co.OpenImage(c.Scope(), "icons/texture.png")
		case model.ResourceKindModel:
			return co.OpenImage(c.Scope(), "icons/model.png")
		case model.ResourceKindScene:
			return co.OpenImage(c.Scope(), "icons/scene.png")
		default:
			return co.OpenImage(c.Scope(), "icons/broken-image.png")
		}
	}

	return co.New(std.TabbarTab, func() {
		co.WithData(std.TabbarTabData{
			Icon:     iconForModelKind(c.resource.Kind()),
			Text:     c.resource.Name(),
			Selected: c.isSelected,
		})
		co.WithCallbackData(std.TabbarTabCallbackData{
			OnClick: func() {
				c.controller.OnSelectEditor(c.editor)
			},
			OnClose: func() {
				c.controller.OnCloseEditor(c.editor)
			},
		})
	})
}
