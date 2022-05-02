package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
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
	layoutData.Height = optional.Value(EditboxHeight)

	var callbackData EditboxCallbackData
	props.InjectOptionalCallbackData(&callbackData, EditboxCallbackData{})

	essence := co.UseState(func() *editboxEssence {
		return &editboxEssence{}
	}).Get()
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
			Focusable: optional.Value(true),
			IdealSize: optional.Value(ui.NewSize(
				200,
				int(essence.font.TextSize(essence.text, essence.fontSize).Y)),
			),
		})
		co.WithLayoutData(layoutData)
	})
}))

var _ ui.ElementKeyboardHandler = (*editboxEssence)(nil)
var _ ui.ElementRenderHandler = (*editboxEssence)(nil)

type editboxEssence struct {
	text         string
	volatileText string
	font         *ui.Font
	fontSize     float32
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

func (e *editboxEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
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
	canvas.Reset()
	canvas.RoundRectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
		sprec.NewVec4(5, 5, 5, 5),
	)
	canvas.Fill(ui.Fill{
		Color: BackgroundColor,
	})

	textBounds := e.font.TextSize(text, e.fontSize)
	canvas.Reset()
	canvas.FillText(text, sprec.NewVec2(2, (float32(size.Height)-textBounds.Y)/2), ui.Typography{
		Font:  e.font,
		Size:  e.fontSize,
		Color: ui.Black(),
	})

	canvas.Reset()
	canvas.SetStrokeSize(1.0)
	canvas.SetStrokeColor(strokeColor)
	canvas.RoundRectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
		sprec.NewVec4(5, 5, 5, 5),
	)
	canvas.Stroke()
}

func (e *editboxEssence) handleTextChange(text string) {
	if e.onChanged != nil {
		e.onChanged(text)
	}
}
