package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type ToolbarData struct {
	Flipped bool
}

var Toolbar = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data ToolbarData
	props.InjectOptionalData(&data, ToolbarData{})

	essence := co.UseState(func() *toolbarEssence {
		return &toolbarEssence{}
	}).Get()

	var layoutData mat.LayoutData
	props.InjectOptionalLayoutData(&layoutData, mat.LayoutData{})
	layoutData.Height = optional.Value(ToolbarHeight)

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

func (e *toolbarEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	size := element.Bounds().Size

	canvas.Reset()
	canvas.Rectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
	)
	canvas.Fill(ui.Fill{
		Color: ToolbarColor,
	})

	canvas.Reset()
	canvas.SetStrokeSize(float32(ToolbarBorderSize))
	canvas.SetStrokeColor(ToolbarBorderColor)
	canvas.MoveTo(sprec.NewVec2(0, float32(size.Height)))
	canvas.LineTo(sprec.NewVec2(float32(size.Width), float32(size.Height)))
	canvas.MoveTo(sprec.NewVec2(float32(size.Width), 0))
	canvas.LineTo(sprec.NewVec2(0, 0))
	canvas.Stroke()
}
