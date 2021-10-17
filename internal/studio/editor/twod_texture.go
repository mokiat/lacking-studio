package editor

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewTwoDTexture() *TwoDTexture {
	return &TwoDTexture{}
}

type TwoDTexture struct {
}

func (t *TwoDTexture) RenderSidePanel() co.Instance {
	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentCenter,
				ContentSpacing:   5,
			}),
		})

		co.WithChild("drawer_asset", co.New(AssetDrawer, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})

		}))
	})
}

var AssetDrawer = co.Define(func(props co.Properties) co.Instance {
	return co.New(widget.Accordion, func() {
		co.WithData(widget.AccordionData{})

		co.WithLayoutData(mat.LayoutData{})

		co.WithChild("title", co.New(mat.Container, func() {

		}))

		if len(props.Children()) > 0 {
			co.WithChild("content", props.Children()[0])
		}
	})
})
