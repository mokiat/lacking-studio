package view

import (
	"fmt"

	"github.com/mokiat/gog/filter"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
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
	RenderEditor(editor *model.Editor, scope co.Scope, layoutData mat.LayoutData) co.Instance
}

type StudioData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var Studio = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data             = co.GetData[StudioData](props)
		studioModel      = data.StudioModel
		studioController = data.StudioController
	)

	mvc.UseBinding(studioModel, filter.True[mvc.Change]()) // TODO

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: opt.V(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithScope(scope)

		co.WithChild("top", co.New(StudioHeader, func() {
			co.WithData(StudioHeaderData{
				StudioModel:      studioModel,
				StudioController: studioController,
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentTop,
			})
		}))

		if editor := studioModel.SelectedEditor(); editor != nil {
			key := fmt.Sprintf("center-%s", editor.Resource().ID())
			instance := studioController.RenderEditor(editor, scope, mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})
			co.WithChild(key, instance)
		}
	})
})
