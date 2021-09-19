package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

var Tabbar = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var layoutData mat.LayoutData
	props.InjectOptionalLayoutData(&layoutData, mat.LayoutData{})
	layoutData.Height = optional.NewInt(TabbarHeight)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(LightGray),
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    0,
				Bottom: 0,
			},
			Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
				ContentAlignment: mat.AlignmentCenter,
				ContentSpacing:   TabbarItemSpacing,
			}),
		})
		co.WithLayoutData(layoutData)
		co.WithChildren(props.Children())
	})
}))
