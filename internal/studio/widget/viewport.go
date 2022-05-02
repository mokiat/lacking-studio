package widget

import (
	"log"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type ViewportData struct {
	API    render.API
	Scene  *graphics.Scene
	Camera *graphics.Camera
}

type ViewportCallbackData struct {
	OnUpdate     func()
	OnMouseEvent func(event ViewportMouseEvent) bool
}

var defaultViewportCallbackData = ViewportCallbackData{
	OnUpdate: func() {},
	OnMouseEvent: func(event ViewportMouseEvent) bool {
		return false
	},
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

	essence := co.UseLifecycle(func(handle co.LifecycleHandle) *viewportEssence {
		return &viewportEssence{
			api: data.API,
		}
	})

	essence.onUpdate = callbackData.OnUpdate
	essence.onMouseEvent = callbackData.OnMouseEvent
	essence.scene = data.Scene
	essence.camera = data.Camera

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   essence,
			Focusable: optional.Value(true),
		})
		co.WithLayoutData(props.LayoutData())
	})
})

var _ ui.ElementKeyboardHandler = (*viewportEssence)(nil)
var _ ui.ElementMouseHandler = (*viewportEssence)(nil)
var _ ui.ElementRenderHandler = (*viewportEssence)(nil)

type viewportEssence struct {
	co.BaseLifecycle

	api         render.API
	texture     render.Texture
	framebuffer render.Framebuffer
	width       int
	height      int

	onUpdate     func()
	onMouseEvent func(event ViewportMouseEvent) bool
	scene        *graphics.Scene
	camera       *graphics.Camera
}

func (e *viewportEssence) OnCreate(props co.Properties) {
	e.ensureFramebuffer(800, 600)
}

func (e *viewportEssence) OnDestroy() {
	e.releaseFramebuffer()
}

func (e *viewportEssence) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	log.Printf("Viewport keyboard event: %#v", event)
	return true
}

func (e *viewportEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	x := -1.0 + 2.0*float32(event.Position.X)/float32(element.Bounds().Width)
	y := 1.0 - 2.0*float32(event.Position.Y)/float32(element.Bounds().Height)
	return e.onMouseEvent(ViewportMouseEvent{
		MouseEvent: event,
		X:          x,
		Y:          y,
	})
}

func (e *viewportEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	e.onUpdate()
	if e.scene != nil && e.camera != nil {
		canvas.DrawSurface(e, ui.NewPosition(0, 0), element.ContentBounds().Size)
	} else {
		size := element.ContentBounds().Size
		canvas.Reset()
		canvas.Rectangle(
			sprec.ZeroVec2(),
			sprec.NewVec2(float32(size.Width), float32(size.Height)),
		)
		canvas.Fill(ui.Fill{
			Rule:  ui.FillRuleSimple,
			Color: ui.Black(),
		})
	}
}

func (e *viewportEssence) Render(width, height int) render.Texture {
	if e.scene != nil && e.camera != nil {
		e.ensureFramebuffer(width, height)
		e.scene.RenderFramebuffer(e.framebuffer, graphics.Viewport{
			X:      0,
			Y:      0,
			Width:  e.width,
			Height: e.height,
		}, e.camera)
	}
	return e.texture
}

func (e *viewportEssence) ensureFramebuffer(width, height int) {
	e.releaseFramebuffer()
	e.width = width
	e.height = height
	e.texture = e.api.CreateColorTexture2D(render.ColorTexture2DInfo{
		Width:           width,
		Height:          height,
		Wrapping:        render.WrapModeClamp,
		Filtering:       render.FilterModeNearest,
		Mipmapping:      false,
		GammaCorrection: true,
		Format:          render.DataFormatRGBA8,
	})
	e.framebuffer = e.api.CreateFramebuffer(render.FramebufferInfo{
		ColorAttachments: [4]render.Texture{
			e.texture,
		},
	})
}

func (e *viewportEssence) releaseFramebuffer() {
	if e.framebuffer != nil {
		e.framebuffer.Release()
		e.framebuffer = nil
	}
	if e.texture != nil {
		e.texture.Release()
		e.texture = nil
	}
}
