package view

import (
	"fmt"

	"github.com/mokiat/gog/filter"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type StudioController interface {
	OnSave()
	OnUndo()
	OnRedo()
	OnToggleProperties()
	OnCreateResource(kind model.ResourceKind) *model.Resource
	OnOpenResource(id string)
	OnCloneResource(id string) *model.Resource
	OnDeleteResource(id string)
	OnSelectEditor(editor *model.Editor)
	OnCloseEditor(editor *model.Editor)
	RenderEditor(editor *model.Editor, scope co.Scope, layoutData any) co.Instance
}

type StudioData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var Studio = mvc.Wrap(co.Define(&studioComponent{}))

type studioComponent struct {
	co.BaseComponent

	studioModel      *model.Studio
	studioController StudioController
}

func (c *studioComponent) OnUpsert() {
	data := co.GetData[StudioData](c.Properties())
	c.studioModel = data.StudioModel
	c.studioController = data.StudioController
	mvc.UseBinding(c.Scope(), c.studioModel, filter.True[mvc.Change]()) // TODO
}

func (c *studioComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout:          layout.Frame(),
		})
		co.WithScope(c.Scope())

		if editor := c.studioModel.SelectedEditor(); editor != nil {
			key := fmt.Sprintf("center-%s", editor.Resource().ID())
			instance := c.studioController.RenderEditor(editor, c.Scope(), layout.Data{
				VerticalAlignment:   layout.VerticalAlignmentCenter,
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
			})
			co.WithChild(key, instance)
		}
	})
}
