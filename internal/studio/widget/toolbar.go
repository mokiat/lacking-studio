package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type ToolbarData struct {
	Flipped bool
}

var Toolbar = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data ToolbarData
	props.InjectOptionalData(&data, ToolbarData{})

	var essence *toolbarEssence
	co.UseState(func() interface{} {
		return &toolbarEssence{}
	}).Inject(&essence)

	var layoutData mat.LayoutData
	props.InjectOptionalLayoutData(&layoutData, mat.LayoutData{})
	layoutData.Height = optional.NewInt(ToolbarHeight)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    ToolbarBorderSize,
				Bottom: ToolbarBorderSize,
			},
			Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
				ContentAlignment: mat.AlignmentCenter,
				ContentSpacing:   ToolbarItemSpacing,
				Flipped:          data.Flipped,
			}),
		})
		co.WithLayoutData(layoutData)
		co.WithChildren(props.Children())
	})
}))

var _ ui.ElementRenderHandler = (*toolbarEssence)(nil)

type toolbarEssence struct{}

func (e *toolbarEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	size := element.Bounds().Size

	canvas.Shape().Begin(ui.Fill{
		Color: ToolbarColor,
	})
	canvas.Shape().Rectangle(
		ui.NewPosition(0, 0),
		size,
	)
	canvas.Shape().End()

	stroke := ui.Stroke{
		Color: ToolbarBorderColor,
		Size:  ToolbarBorderSize,
	}
	canvas.Contour().Begin()
	canvas.Contour().MoveTo(ui.NewPosition(0, size.Height), stroke)
	canvas.Contour().LineTo(ui.NewPosition(size.Width, size.Height), stroke)
	canvas.Contour().MoveTo(ui.NewPosition(size.Width, 0), stroke)
	canvas.Contour().LineTo(ui.NewPosition(0, 0), stroke)
	canvas.Contour().End()
}
