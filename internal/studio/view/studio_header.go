package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type StudioHeaderData struct {
	StudioModel *model.Studio
}

var StudioHeader = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data   = co.GetData[StudioHeaderData](props)
		studio = data.StudioModel
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
				StudioModel: studio,
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))

		co.WithChild("tabbar", co.New(StudioTabbar, func() {
			co.WithData(StudioTabbarData{
				StudioModel: studio,
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))
	})
})
