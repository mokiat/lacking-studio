package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type StudioHeaderData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var StudioHeader = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data       = co.GetData[StudioHeaderData](props)
		studio     = data.StudioModel
		controller = data.StudioController
	)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
			}),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("toolbar", co.New(StudioToolbar, func() {
			co.WithData(StudioToolbarData{
				StudioModel:      studio,
				StudioController: controller,
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))

		co.WithChild("tabbar", co.New(StudioTabbar, func() {
			co.WithData(StudioTabbarData{
				StudioModel:      studio,
				StudioController: controller,
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))
	})
})
