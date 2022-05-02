package widget

import (
	"github.com/mokiat/gomath/sprec"
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

	essence := co.UseState(func() *paperEssence {
		return &paperEssence{}
	}).Get()

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

func (e *paperEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	size := element.Bounds().Size
	canvas.Reset()
	canvas.Rectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
	)
	canvas.Fill(ui.Fill{
		Color: SurfaceColor,
	})

	canvas.Reset()
	canvas.SetStrokeSize(PaperBorderSize)
	canvas.SetStrokeColor(Gray)
	canvas.MoveTo(sprec.NewVec2(0, 0))
	canvas.LineTo(sprec.NewVec2(0, float32(size.Height)))
	canvas.LineTo(sprec.NewVec2(float32(size.Width), float32(size.Height)))
	canvas.LineTo(sprec.NewVec2(float32(size.Width), 0))
	canvas.CloseLoop()
}
