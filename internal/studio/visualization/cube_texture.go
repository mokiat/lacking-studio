package visualization

import (
	"image"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewCubeTexture(api render.API, engine *graphics.Engine, texModel *model.CubeTexture) *CubeTexture {
	scene := engine.CreateScene()

	sky := scene.Sky()
	sky.SetBackgroundColor(sprec.NewVec3(0.2, 0.2, 0.2))

	camera := scene.CreateCamera()
	camera.SetMatrix(dprec.IdentityMat4())
	camera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	camera.SetFoV(sprec.Degrees(66))
	camera.SetAutoExposure(true)
	camera.SetExposure(1.0)
	camera.SetAutoFocus(false)

	result := &CubeTexture{
		texModel:    texModel,
		api:         api,
		engine:      engine,
		scene:       scene,
		camera:      camera,
		cameraPitch: dprec.Degrees(0),
		cameraYaw:   dprec.Degrees(0),
		cameraFoV:   sprec.Degrees(66),
	}
	result.createGraphicsRepresentation()
	result.subscribeToModel()
	return result
}

type CubeTexture struct {
	texModel        *model.CubeTexture
	texSubscription mvc.Subscription

	api         render.API
	engine      *graphics.Engine
	scene       *graphics.Scene
	camera      *graphics.Camera
	cameraPitch dprec.Angle
	cameraYaw   dprec.Angle
	cameraFoV   sprec.Angle
	texture     *graphics.CubeTexture

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
}

func (t *CubeTexture) TakeSnapshot(size ui.Size) image.Image {
	colorTexture := t.api.CreateColorTexture2D(render.ColorTexture2DInfo{
		Width:           size.Width,
		Height:          size.Height,
		Wrapping:        render.WrapModeClamp,
		Filtering:       render.FilterModeNearest,
		Mipmapping:      false,
		GammaCorrection: false,
		Format:          render.DataFormatRGBA8,
	})
	defer colorTexture.Release()

	framebuffer := t.api.CreateFramebuffer(render.FramebufferInfo{
		ColorAttachments: [4]render.Texture{
			colorTexture,
		},
	})
	defer framebuffer.Release()

	buffer := t.api.CreatePixelTransferBuffer(render.BufferInfo{
		Size: 4 * size.Width * size.Height,
	})
	defer buffer.Release()

	t.api.BeginRenderPass(render.RenderPassInfo{
		Framebuffer: framebuffer,
		Viewport: render.Area{
			X:      0,
			Y:      0,
			Width:  size.Width,
			Height: size.Height,
		},
		DepthLoadOp:    render.LoadOperationDontCare,
		DepthStoreOp:   render.StoreOperationDontCare,
		StencilLoadOp:  render.LoadOperationDontCare,
		StencilStoreOp: render.StoreOperationDontCare,
		Colors: [4]render.ColorAttachmentInfo{
			{
				LoadOp:     render.LoadOperationClear,
				ClearValue: [4]float32{0.0, 0.0, 0.0, 1.0},
			},
		},
	})

	t.scene.RenderFramebuffer(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  size.Width,
		Height: size.Height,
	})

	commands := t.api.CreateCommandQueue()
	defer commands.Release()
	commands.CopyContentToBuffer(render.CopyContentToBufferInfo{
		Buffer: buffer,
		X:      0,
		Y:      0,
		Width:  size.Width,
		Height: size.Height,
		Format: render.DataFormatRGBA8,
	})
	t.api.SubmitQueue(commands)

	previewImg := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))
	buffer.Fetch(render.BufferFetchInfo{
		Offset: 0,
		Target: previewImg.Pix,
	})
	for y := 0; y < size.Height/2; y++ {
		topOffset := y * (4 * size.Width)
		bottomOffset := (size.Height - y - 1) * (4 * size.Width)
		for x := 0; x < size.Width*4; x++ {
			previewImg.Pix[topOffset+x], previewImg.Pix[bottomOffset+x] =
				previewImg.Pix[bottomOffset+x], previewImg.Pix[topOffset+x]
		}
	}

	t.api.EndRenderPass()
	return previewImg
}

func (t *CubeTexture) OnViewportRender(framebuffer render.Framebuffer, size ui.Size) {
	t.camera.SetMatrix(dprec.Mat4MultiProd(
		dprec.RotationMat4(t.cameraYaw, 0.0, 1.0, 0.0),
		dprec.RotationMat4(t.cameraPitch, 1.0, 0.0, 0.0),
	))
	t.camera.SetFoV(t.cameraFoV)

	t.scene.RenderFramebuffer(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  size.Width,
		Height: size.Height,
	})
}

func (t *CubeTexture) OnViewportMouseEvent(event mat.ViewportMouseEvent) bool {
	switch event.Type {
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonMiddle {
			t.rotatingCamera = true
			t.oldMouseX = event.Position.X
			t.oldMouseY = event.Position.Y
			return true
		}
	case ui.MouseEventTypeMove:
		if t.rotatingCamera {
			t.cameraPitch += dprec.Degrees(float64(event.Position.Y-t.oldMouseY) / 5)
			t.cameraYaw += dprec.Degrees(float64(event.Position.X-t.oldMouseX) / 5)
			t.oldMouseX = event.Position.X
			t.oldMouseY = event.Position.Y
			return true
		}
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonMiddle {
			t.rotatingCamera = false
			return true
		}
	case ui.MouseEventTypeLeave:
		t.rotatingCamera = false
		return true
	case ui.MouseEventTypeScroll:
		fov := t.cameraFoV.Degrees()
		fov -= 2 * float32(event.ScrollY)
		fov = sprec.Clamp(fov, 0.1, 179.0)
		t.cameraFoV = sprec.Degrees(fov)
		return true
	}
	return false
}

func (t *CubeTexture) Destroy() {
	t.unsubscribeFromModel()
	t.deleteGraphicsRepresentation()
	t.scene.Delete()
}

func (t *CubeTexture) subscribeToModel() {
	t.texSubscription = t.texModel.Subscribe(func(ch mvc.Change) {
		t.deleteGraphicsRepresentation()
		t.createGraphicsRepresentation()
	})
}

func (t *CubeTexture) unsubscribeFromModel() {
	t.texSubscription.Delete()
}

func (t *CubeTexture) createGraphicsRepresentation() {
	definition := t.buildGraphicsDefinition(t.texModel)
	t.texture = t.engine.CreateCubeTexture(definition)
	sky := t.scene.Sky()
	sky.SetSkybox(t.texture)
}

func (t *CubeTexture) deleteGraphicsRepresentation() {
	sky := t.scene.Sky()
	sky.SetSkybox(nil)
	t.texture.Delete()
}

func (t *CubeTexture) buildGraphicsDefinition(src *model.CubeTexture) graphics.CubeTextureDefinition {
	return graphics.CubeTextureDefinition{
		Dimension:      src.Dimension(),
		Filtering:      assetToGraphicsFilter(src.Filtering()),
		InternalFormat: assetToGraphicsInternalFormat(src.Format()),
		DataFormat:     assetToGraphicsDataFormat(src.Format()),
		FrontSideData:  src.FrontData(),
		BackSideData:   src.BackData(),
		LeftSideData:   src.LeftData(),
		RightSideData:  src.RightData(),
		TopSideData:    src.TopData(),
		BottomSideData: src.BottomData(),
	}
}
