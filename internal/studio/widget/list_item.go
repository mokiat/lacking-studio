package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type ListItemData struct {
	Selected bool
}

type ListItemCallbackData struct {
	OnSelected func()
}

var defaultListItemCallbackData = ListItemCallbackData{
	OnSelected: func() {},
}

var ListItem = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var lifecycle *listItemLifecycle
	co.UseLifecycle(func(handle co.LifecycleHandle) co.Lifecycle {
		return &listItemLifecycle{
			Lifecycle: co.NewBaseLifecycle(),
		}
	}, &lifecycle)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Padding: ui.Spacing{
				Left:   2,
				Right:  2,
				Top:    2,
				Bottom: 2,
			},
			Essence: lifecycle,
			Layout:  mat.NewFillLayout(),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
}))

var _ ui.ElementMouseHandler = (*listItemLifecycle)(nil)
var _ ui.ElementRenderHandler = (*listItemLifecycle)(nil)

type listItemLifecycle struct {
	co.Lifecycle

	selected   bool
	onSelected func()
	state      buttonState
}

func (e *listItemLifecycle) OnCreate(props co.Properties) {
	e.OnUpdate(props)
	e.state = buttonStateUp
}

func (e *listItemLifecycle) OnUpdate(props co.Properties) {
	var data ListItemData
	props.InjectOptionalData(&data, ListItemData{})
	e.selected = data.Selected

	var callbackData ListItemCallbackData
	props.InjectOptionalCallbackData(&callbackData, defaultListItemCallbackData)
	e.onSelected = callbackData.OnSelected
}

func (e *listItemLifecycle) IsSelected() bool {
	return e.selected
}

func (e *listItemLifecycle) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
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
				e.onSelected()
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

func (e *listItemLifecycle) OnRender(element *ui.Element, canvas ui.Canvas) {
	var backgroundColor ui.Color
	switch e.state {
	case buttonStateOver:
		backgroundColor = ui.ColorWithAlpha(SecondaryColor, 128)
	case buttonStateDown:
		backgroundColor = ui.ColorWithAlpha(SecondaryColor, 196)
	default:
		backgroundColor = ui.Transparent()
	}
	if e.selected {
		backgroundColor = SecondaryColor
	}

	if !backgroundColor.Transparent() {
		canvas.Shape().Begin(ui.Fill{
			Color: backgroundColor,
		})
		canvas.Shape().Rectangle(
			ui.NewPosition(0, 0),
			element.Bounds().Size,
		)
		canvas.Shape().End()
	}
}
