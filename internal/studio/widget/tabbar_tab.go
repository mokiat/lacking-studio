package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type TabbarTabData struct {
	Icon     ui.Image
	Text     string
	Selected bool
}

type TabbarTabCallbackData struct {
	OnClick func()
	OnClose func()
}

var TabbarTab = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data TabbarTabData
	props.InjectOptionalData(&data, TabbarTabData{})

	var layoutData mat.LayoutData
	props.InjectOptionalLayoutData(&layoutData, mat.LayoutData{})
	layoutData.Height = optional.NewInt(TabbarItemHeight)

	var callbackData TabbarTabCallbackData
	props.InjectOptionalCallbackData(&callbackData, TabbarTabCallbackData{})

	var essence *tabbarTabEssence
	co.UseState(func() interface{} {
		return &tabbarTabEssence{}
	}).Inject(&essence)
	essence.selected = data.Selected
	essence.onClick = callbackData.OnClick

	var closeEssence *buttonEssence
	co.UseState(func() interface{} {
		return &buttonEssence{}
	}).Inject(&closeEssence)
	closeEssence.onClick = callbackData.OnClose

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
				ContentAlignment: mat.AlignmentCenter,
				ContentSpacing:   5,
			}),
			Padding: ui.Spacing{
				Top:   5,
				Left:  10,
				Right: 10,
			},
		})
		co.WithLayoutData(layoutData)

		if data.Icon != nil {
			co.WithChild("icon", co.New(mat.Picture, func() {
				co.WithData(mat.PictureData{
					Image:      data.Icon,
					ImageColor: optional.NewColor(ui.Black()),
					Mode:       mat.ImageModeFit,
				})
				co.WithLayoutData(mat.LayoutData{
					Width:  optional.NewInt(24),
					Height: optional.NewInt(24),
				})
			}))
		}

		if data.Text != "" {
			co.WithChild("text", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "regular"),
					FontSize:  optional.NewInt(20),
					FontColor: optional.NewColor(ui.Black()),
					Text:      data.Text,
				})
				co.WithLayoutData(mat.LayoutData{})
			}))
		}

		if data.Selected {
			co.WithChild("close", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Essence: closeEssence,
					Layout:  mat.NewFillLayout(),
				})

				co.WithLayoutData(mat.LayoutData{
					Width:  optional.NewInt(24),
					Height: optional.NewInt(24),
				})

				co.WithChild("icon", co.New(mat.Picture, func() {
					co.WithData(mat.PictureData{
						Image:      co.OpenImage("resources/icons/close.png"),
						ImageColor: optional.NewColor(ui.Black()),
						Mode:       mat.ImageModeFit,
					})
					co.WithLayoutData(mat.LayoutData{
						Width:  optional.NewInt(24),
						Height: optional.NewInt(24),
					})
				}))
			}))
		}
	})
}))

var _ ui.ElementMouseHandler = (*tabbarTabEssence)(nil)
var _ ui.ElementRenderHandler = (*tabbarTabEssence)(nil)

type tabbarTabEssence struct {
	state    buttonState
	selected bool
	onClick  func()
}

func (e *tabbarTabEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
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
				e.handleClick()
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

func (e *tabbarTabEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	var backgroundColor ui.Color
	if e.selected {
		backgroundColor = ui.White()
	} else {
		switch e.state {
		case buttonStateOver:
			backgroundColor = LightGray
		case buttonStateDown:
			backgroundColor = DarkGray
		default:
			backgroundColor = ui.Transparent()
		}
	}

	size := element.Bounds().Size
	if !backgroundColor.Transparent() {
		canvas.Shape().Begin(ui.Fill{
			Color: backgroundColor,
		})
		canvas.Shape().MoveTo(ui.NewPosition(0, size.Height))
		canvas.Shape().LineTo(ui.NewPosition(size.Width, size.Height))
		canvas.Shape().LineTo(ui.NewPosition(size.Width, TabbarItemRadius))
		canvas.Shape().QuadTo(ui.NewPosition(size.Width, 0), ui.NewPosition(size.Width-TabbarItemRadius, 0))
		canvas.Shape().LineTo(ui.NewPosition(TabbarItemRadius, 0))
		canvas.Shape().QuadTo(ui.NewPosition(0, 0), ui.NewPosition(0, TabbarItemRadius))
		canvas.Shape().End()
	}
}

func (e *tabbarTabEssence) handleClick() {
	if e.onClick != nil {
		e.onClick()
	}
}

var _ ui.ElementMouseHandler = (*buttonEssence)(nil)
var _ ui.ElementRenderHandler = (*buttonEssence)(nil)

type buttonEssence struct {
	state   buttonState
	onClick func()
}

func (e *buttonEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
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
				e.handleClick()
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

func (e *buttonEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	var backgroundColor ui.Color
	switch e.state {
	case buttonStateOver:
		backgroundColor = Gray
	case buttonStateDown:
		backgroundColor = DarkGray
	default:
		backgroundColor = ui.Transparent()
	}

	size := element.Bounds().Size
	if !backgroundColor.Transparent() {
		canvas.Shape().Begin(ui.Fill{
			Color: backgroundColor,
		})
		canvas.Shape().Rectangle(ui.NewPosition(0, 0), size)
		canvas.Shape().End()
	}
}

func (e *buttonEssence) handleClick() {
	if e.onClick != nil {
		e.onClick()
	}
}
