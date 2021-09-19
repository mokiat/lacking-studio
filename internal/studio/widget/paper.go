package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var defaultPaperData = PaperData{
	Layout: ui.NewFillLayout(),
	Padding: ui.Spacing{
		Left:   1,
		Right:  1,
		Top:    1,
		Bottom: 1,
	},
}

type PaperData struct {
	Padding ui.Spacing
	Layout  ui.Layout
}

var Paper = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data PaperData
	props.InjectOptionalData(&data, defaultPaperData)

	var essence *paperEssence
	co.UseState(func() interface{} {
		return &paperEssence{}
	}).Inject(&essence)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Layout:  data.Layout,
			Padding: data.Padding,
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
}))

var _ ui.ElementRenderHandler = (*paperEssence)(nil)

type paperEssence struct{}

func (e *paperEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	size := element.Bounds().Size
	canvas.Shape().Begin(ui.Fill{
		Color: SurfaceColor,
	})
	canvas.Shape().Rectangle(
		ui.NewPosition(0, 0),
		size,
	)
	canvas.Shape().End()

	stroke := ui.Stroke{
		Size:  PaperBorderSize,
		Color: Gray,
	}
	canvas.Contour().Begin()
	canvas.Contour().MoveTo(ui.NewPosition(0, 0), stroke)
	canvas.Contour().LineTo(ui.NewPosition(0, size.Height), stroke)
	canvas.Contour().LineTo(ui.NewPosition(size.Width, size.Height), stroke)
	canvas.Contour().LineTo(ui.NewPosition(size.Width, 0), stroke)
	canvas.Contour().CloseLoop()
	canvas.Contour().End()
}
