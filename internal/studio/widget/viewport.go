package widget

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type ViewportData struct {
	Scene  graphics.Scene
	Camera graphics.Camera
}

type ViewportCallbackData struct {
	OnUpdate     func()
	OnMouseEvent func(event ViewportMouseEvent)
}

var defaultViewportCallbackData = ViewportCallbackData{
	OnUpdate:     func() {},
	OnMouseEvent: func(event ViewportMouseEvent) {},
}

type ViewportMouseEvent struct {
	ui.MouseEvent
	X float32
	Y float32
}

var Viewport = co.Define(func(props co.Properties) co.Instance {
	var data ViewportData
	props.InjectOptionalData(&data, ViewportData{})

	var callbackData ViewportCallbackData
	props.InjectOptionalCallbackData(&callbackData, defaultViewportCallbackData)

	var essence *viewportEssence
	co.UseState(func() interface{} {
		return &viewportEssence{}
	}).Inject(&essence)
	essence.onUpdate = callbackData.OnUpdate
	essence.onMouseEvent = callbackData.OnMouseEvent
	essence.scene = data.Scene
	essence.camera = data.Camera

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
		})
		co.WithLayoutData(props.LayoutData())
	})
})

var _ ui.ElementMouseHandler = (*viewportEssence)(nil)
var _ ui.ElementRenderHandler = (*viewportEssence)(nil)

type viewportEssence struct {
	onUpdate     func()
	onMouseEvent func(event ViewportMouseEvent)
	scene        graphics.Scene
	camera       graphics.Camera
}

func (e *viewportEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	x := -1.0 + 2.0*float32(event.Position.X)/float32(element.Bounds().Width)
	y := 1.0 - 2.0*float32(event.Position.Y)/float32(element.Bounds().Height)
	e.onMouseEvent(ViewportMouseEvent{
		MouseEvent: event,
		X:          x,
		Y:          y,
	})
	return true
}

func (e *viewportEssence) OnRender(element *ui.Element, canvas ui.Canvas) {
	e.onUpdate()
	if e.scene != nil && e.camera != nil {
		canvas.DrawSurface(e)
	} else {
		canvas.Shape().Begin(ui.Fill{
			Rule:  ui.FillRuleSimple,
			Color: ui.Black(),
		})
		canvas.Shape().Rectangle(
			ui.NewPosition(0, 0),
			element.ContentBounds().Size,
		)
		canvas.Shape().End()
	}
}

func (e *viewportEssence) Render(x, y, width, height int) {
	if e.scene != nil && e.camera != nil {
		e.scene.Render(graphics.Viewport{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}, e.camera)
	}
}
