package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type EditboxData struct {
	Text string
}

type EditboxCallbackData struct {
	OnChanged func(text string)
}

var Editbox = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data EditboxData
	props.InjectOptionalData(&data, EditboxData{})

	var layoutData mat.LayoutData
	props.InjectOptionalLayoutData(&layoutData, mat.LayoutData{})
	layoutData.Height = optional.NewInt(EditboxHeight)

	var callbackData EditboxCallbackData
	props.InjectOptionalCallbackData(&callbackData, EditboxCallbackData{})

	var essence *editboxEssence
	co.UseState(func() interface{} {
		return &editboxEssence{}
	}).Inject(&essence)
	if data.Text != essence.text {
		essence.text = data.Text
		essence.volatileText = data.Text
	}
	essence.font = co.GetFont("roboto", "regular")
	essence.fontSize = 18
	essence.onChanged = callbackData.OnChanged

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   essence,
			Focusable: optional.NewBool(true),
			IdealSize: optional.NewSize(ui.NewSize(200, essence.font.TextSize(essence.text, essence.fontSize).Height)),
		})
		co.WithLayoutData(layoutData)
	})
}))

var _ ui.ElementKeyboardHandler = (*editboxEssence)(nil)
var _ ui.ElementRenderHandler = (*editboxEssence)(nil)

type editboxEssence struct {
	text         string
	volatileText string
	font         ui.Font
	fontSize     int
	onChanged    func(string)
}

func (e *editboxEssence) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Type {
	case ui.KeyboardEventTypeKeyDown, ui.KeyboardEventTypeRepeat:
		switch event.Code {
		case ui.KeyCodeBackspace:
			if len(e.volatileText) > 0 {
				e.volatileText = e.volatileText[:len(e.volatileText)-1]
				co.Window().Invalidate()
			}
		case ui.KeyCodeEscape:
			e.volatileText = e.text
			co.Window().DiscardFocus()
		case ui.KeyCodeEnter:
			e.text = e.volatileText
			e.handleTextChange(e.text)
			co.Window().DiscardFocus()
		}
	case ui.KeyboardEventTypeType:
		e.volatileText += string(event.Rune)
		co.Window().Invalidate()
	}
	return true
}

func (e *editboxEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	var strokeColor ui.Color
	var text string
	if co.Window().IsElementFocused(element) {
		strokeColor = SecondaryColor
		text = e.volatileText + "|"
	} else {
		strokeColor = ui.Black()
		text = e.volatileText
	}

	size := element.Bounds().Size
	canvas.Shape().Begin(ui.Fill{
		Color: BackgroundColor,
	})
	canvas.Shape().RoundRectangle(
		ui.NewPosition(0, 0),
		size,
		ui.RectRoundness{
			TopLeftRadius:     5,
			TopRightRadius:    5,
			BottomLeftRadius:  5,
			BottomRightRadius: 5,
		},
	)
	canvas.Shape().End()

	textBounds := e.font.TextSize(text, e.fontSize)
	canvas.Text().Begin(ui.Typography{
		Font:  e.font,
		Size:  e.fontSize,
		Color: ui.Black(),
	})
	canvas.Text().Line(text, ui.NewPosition(2, (size.Height-textBounds.Height)/2))
	canvas.Text().End()

	canvas.Contour().Begin()
	canvas.Contour().RoundRectangle(
		ui.NewPosition(0, 0),
		size,
		ui.RectRoundness{
			TopLeftRadius:     5,
			TopRightRadius:    5,
			BottomLeftRadius:  5,
			BottomRightRadius: 5,
		},
		ui.Stroke{
			Size:  1,
			Color: strokeColor,
		},
	)
	canvas.Contour().End()
}

func (e *editboxEssence) handleTextChange(text string) {
	if e.onChanged != nil {
		e.onChanged(text)
	}
}
