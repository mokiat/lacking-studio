package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

var defaultAccordionData = AccordionData{}

var defaultAccordionCallbackData = AccordionCallbackData{
	OnToggle: func() {},
}

type AccordionData struct {
	Title    string
	Expanded bool
}

type AccordionCallbackData struct {
	OnToggle func()
}

var Accordion = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data AccordionData
	props.InjectOptionalData(&data, defaultAccordionData)

	var callbackData AccordionCallbackData
	props.InjectOptionalCallbackData(&callbackData, defaultAccordionCallbackData)

	headerEssence := co.UseState(func() *accordionHeaderEssence {
		return &accordionHeaderEssence{}
	}).Get()
	headerEssence.onToggle = callbackData.OnToggle

	var icon *ui.Image
	if data.Expanded {
		icon = co.OpenImage("resources/icons/expanded.png")
	} else {
		icon = co.OpenImage("resources/icons/collapsed.png")
	}

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
			}),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("header", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Essence: headerEssence,
				Padding: ui.Spacing{
					Left:   2,
					Right:  2,
					Top:    2,
					Bottom: 2,
				},
				Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
					ContentAlignment: mat.AlignmentCenter,
					ContentSpacing:   5,
				}),
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})

			co.WithChild("icon", co.New(mat.Picture, func() {
				co.WithData(mat.PictureData{
					Image:      icon,
					ImageColor: optional.Value(ui.Black()),
					Mode:       mat.ImageModeFit,
				})
				co.WithLayoutData(mat.LayoutData{
					Width:  optional.Value(24),
					Height: optional.Value(24),
				})
			}))

			co.WithChild("title", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "regular"),
					FontSize:  optional.Value(float32(20)),
					FontColor: optional.Value(ui.Black()),
					Text:      data.Title,
				})
			}))
		}))

		if data.Expanded {
			for _, child := range props.Children() {
				co.WithChild(child.Key(), child)
			}
		}
	})
}))

var _ ui.ElementMouseHandler = (*accordionHeaderEssence)(nil)
var _ ui.ElementRenderHandler = (*accordionHeaderEssence)(nil)

type accordionHeaderEssence struct {
	state    buttonState
	onToggle func()
}

func (e *accordionHeaderEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
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
				e.onToggle()
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

func (e *accordionHeaderEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var backgroundColor ui.Color
	switch e.state {
	case buttonStateOver:
		backgroundColor = LightGray
	case buttonStateDown:
		backgroundColor = DarkGray
	default:
		backgroundColor = ui.White()
	}

	size := element.Bounds().Size
	canvas.Reset()
	canvas.Rectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
	)
	canvas.Fill(ui.Fill{
		Color: backgroundColor,
	})

	canvas.Reset()
	canvas.SetStrokeSize(PaperBorderSize)
	canvas.SetStrokeColor(Gray)
	canvas.MoveTo(sprec.NewVec2(0, 0))
	canvas.LineTo(sprec.NewVec2(0, float32(size.Height)))
	canvas.LineTo(sprec.NewVec2(float32(size.Width), float32(size.Height)))
	canvas.LineTo(sprec.NewVec2(float32(size.Width), 0))
	canvas.CloseLoop()
	canvas.Stroke()
}
