package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type ToolbarButtonData struct {
	Icon     *ui.Image
	Text     string
	Disabled bool
	Selected bool
	Vertical bool
}

type ToolbarButtonCallbackData struct {
	ClickListener mat.ClickListener
}

var ToolbarButton = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data ToolbarButtonData
	props.InjectOptionalData(&data, ToolbarButtonData{})

	var layoutData mat.LayoutData
	props.InjectOptionalLayoutData(&layoutData, mat.LayoutData{})
	layoutData.Height = optional.Value(ToolbarItemHeight)

	var callbackData ToolbarButtonCallbackData
	props.InjectOptionalCallbackData(&callbackData, ToolbarButtonCallbackData{})

	essence := co.UseState(func() *toolbarButtonEssence {
		return &toolbarButtonEssence{}
	}).Get()
	essence.selected = data.Selected
	essence.clickListener = callbackData.ClickListener

	multiplierColor := ui.Black()
	if data.Disabled {
		multiplierColor = ui.Gray()
	}

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
				ContentAlignment: mat.AlignmentCenter,
				ContentSpacing:   5,
			}),
			Padding: ui.Spacing{
				Left:  4,
				Right: 4,
			},
			Enabled: optional.Value(!data.Disabled),
		})
		co.WithLayoutData(layoutData)

		if data.Icon != nil {
			co.WithChild("icon", co.New(mat.Picture, func() {
				co.WithData(mat.PictureData{
					Image:      data.Icon,
					ImageColor: optional.Value(multiplierColor),
					Mode:       mat.ImageModeFit,
				})
				co.WithLayoutData(mat.LayoutData{
					Width:  optional.Value(24),
					Height: optional.Value(24),
				})
			}))
		}

		if data.Text != "" {
			co.WithChild("text", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "regular"),
					FontSize:  optional.Value(float32(20)),
					FontColor: optional.Value(multiplierColor),
					Text:      data.Text,
				})
				co.WithLayoutData(mat.LayoutData{})
			}))
		}
	})
}))

var _ ui.ElementMouseHandler = (*toolbarButtonEssence)(nil)
var _ ui.ElementRenderHandler = (*toolbarButtonEssence)(nil)

type toolbarButtonEssence struct {
	state         buttonState
	selected      bool
	clickListener mat.ClickListener
}

func (e *toolbarButtonEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	context := element.Context()
	switch event.Type {
	case ui.MouseEventTypeEnter:
		e.state = buttonStateOver
		context.Window().Invalidate()
	case ui.MouseEventTypeLeave:
		e.state = buttonStateUp
		context.Window().Invalidate()
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonLeft {
			if e.state == buttonStateDown {
				e.onClick()
			}
			e.state = buttonStateOver
			context.Window().Invalidate()
		}
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonLeft {
			e.state = buttonStateDown
			context.Window().Invalidate()
		}
	}
	return true
}

func (e *toolbarButtonEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var backgroundColor ui.Color
	switch e.state {
	case buttonStateOver:
		backgroundColor = LightGray
	case buttonStateDown:
		backgroundColor = DarkGray
	default:
		backgroundColor = ui.Transparent()
	}

	size := element.Bounds().Size
	if !backgroundColor.Transparent() {
		canvas.Reset()
		canvas.Rectangle(
			sprec.ZeroVec2(),
			sprec.NewVec2(float32(size.Width), float32(size.Height)),
		)
		canvas.Fill(ui.Fill{
			Color: backgroundColor,
		})
	}
	if e.selected {
		canvas.Reset()
		canvas.Rectangle(
			sprec.NewVec2(0, float32(size.Height)-5),
			sprec.NewVec2(float32(size.Width), 5),
		)
		canvas.Fill(ui.Fill{
			Color: SecondaryColor,
		})
	}
}

func (e *toolbarButtonEssence) onClick() {
	if e.clickListener != nil {
		e.clickListener()
	}
}
