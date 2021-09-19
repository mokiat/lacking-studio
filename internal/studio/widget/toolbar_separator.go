package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

var ToolbarSeparator = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var essence *toolbarSeparatorEssence
	co.UseState(func() interface{} {
		return &toolbarSeparatorEssence{}
	}).Inject(&essence)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			IdealSize: optional.NewSize(ui.NewSize(
				ToolbarSeparatorWidth,
				ToolbarItemHeight,
			)),
		})
		co.WithLayoutData(props.LayoutData())
	})
}))

var _ ui.ElementRenderHandler = (*toolbarSeparatorEssence)(nil)

type toolbarSeparatorEssence struct{}

func (e *toolbarSeparatorEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	size := element.Bounds().Size

	lineLength := ToolbarSeparatorLineLength
	linePadding := (size.Height - lineLength) / 2

	stroke := ui.Stroke{
		Color: ToolbarSeparatorLineColor,
		Size:  ToolbarSeparatorLineSize,
	}
	canvas.Contour().Begin()
	canvas.Contour().MoveTo(ui.NewPosition(size.Width/2, linePadding), stroke)
	canvas.Contour().LineTo(ui.NewPosition(size.Width/2, linePadding+lineLength), stroke)
	canvas.Contour().End()
}
