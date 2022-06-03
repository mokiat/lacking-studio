package view

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/filter"
	"github.com/mokiat/lacking/util/optional"
)

type StudioController interface {
	mvc.Reducer
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

	mvc.UseBinding(studioModel, filter.Always[mvc.Change]()) // TODO

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithScope(scope)

		co.WithChild("top", co.New(StudioHeader, func() {
			co.WithData(StudioHeaderData{
				StudioModel: studioModel,
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
