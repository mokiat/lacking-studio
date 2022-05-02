package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

var ToolbarSeparator = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	essence := co.UseState(func() *toolbarSeparatorEssence {
		return &toolbarSeparatorEssence{}
	}).Get()

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			IdealSize: optional.Value(ui.NewSize(
				ToolbarSeparatorWidth,
				ToolbarItemHeight,
			)),
		})
		co.WithLayoutData(props.LayoutData())
	})
}))

var _ ui.ElementRenderHandler = (*toolbarSeparatorEssence)(nil)

type toolbarSeparatorEssence struct{}

func (e *toolbarSeparatorEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	size := element.Bounds().Size

	lineLength := ToolbarSeparatorLineLength
	linePadding := (size.Height - lineLength) / 2

	canvas.Reset()
	canvas.SetStrokeSize(ToolbarSeparatorLineSize)
	canvas.SetStrokeColor(ToolbarSeparatorLineColor)
	canvas.MoveTo(sprec.NewVec2(float32(size.Width)/2, float32(linePadding)))
	canvas.LineTo(sprec.NewVec2(float32(size.Width)/2, float32(linePadding+lineLength)))
	canvas.Stroke()
}
