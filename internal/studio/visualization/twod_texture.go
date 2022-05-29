package visualization

import (
	"encoding/binary"
	"fmt"
	"image"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
)

func NewTwoDTexture(api render.API, engine *graphics.Engine, texModel *model.TwoDTexture) *TwoDTexture {
	scene := engine.CreateScene()

	sky := scene.Sky()
	sky.SetBackgroundColor(sprec.NewVec3(0.2, 0.2, 0.2))

	light := scene.CreateDirectionalLight()
	light.SetIntensity(sprec.NewVec3(1.0, 1.0, 1.0))
	light.SetRotation(sprec.IdentityQuat())

	camera := scene.CreateCamera()
	camera.SetPosition(sprec.NewVec3(0.0, 0.0, 3.0))
	camera.SetRotation(sprec.IdentityQuat())
	camera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	camera.SetFoV(sprec.Degrees(66))
	camera.SetAutoExposure(false)
	camera.SetExposure(3.14)
	camera.SetAutoFocus(false)

	result := &TwoDTexture{
		texModel:    texModel,
		api:         api,
		engine:      engine,
		scene:       scene,
		camera:      camera,
		cameraPitch: sprec.Degrees(0),
		cameraYaw:   sprec.Degrees(0),
		cameraFoV:   sprec.Degrees(66),
	}
	result.createGraphicsRepresentation()
	result.subscribeToModel()
	return result
}

type TwoDTexture struct {
	texModel        *model.TwoDTexture
	texSubscription observer.Subscription

	api          render.API
	engine       *graphics.Engine
	scene        *graphics.Scene
	camera       *graphics.Camera
	cameraPitch  sprec.Angle
	cameraYaw    sprec.Angle
	cameraFoV    sprec.Angle
	mesh         *graphics.Mesh
	meshTemplate *graphics.MeshTemplate
	material     *graphics.Material
	texture      *graphics.TwoDTexture

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
}

func (t *TwoDTexture) TakeSnapshot(size ui.Size) image.Image {
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
	}, t.camera)

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

func (t *TwoDTexture) OnViewportRender(framebuffer render.Framebuffer, size ui.Size) {
	transform := sprec.Mat4MultiProd(
		sprec.RotationMat4(-t.cameraYaw, 0.0, 1.0, 0.0),
		sprec.RotationMat4(-t.cameraPitch, 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, 3.0),
	)
	t.camera.SetPosition(transform.Translation())
	t.camera.SetRotation(transform.RotationQuat())
	t.camera.SetFoV(t.cameraFoV)

	t.scene.RenderFramebuffer(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  size.Width,
		Height: size.Height,
	}, t.camera)
}

func (t *TwoDTexture) OnViewportMouseEvent(event mat.ViewportMouseEvent) bool {
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
			t.cameraPitch += sprec.Degrees(float32(event.Position.Y-t.oldMouseY) / 5)
			t.cameraYaw += sprec.Degrees(float32(event.Position.X-t.oldMouseX) / 5)
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

func (t *TwoDTexture) Destroy() {
	t.deleteGraphicsRepresentation()
	t.unsubscribeFromModel()
	t.scene.Delete()
}

func (t *TwoDTexture) subscribeToModel() {
	t.texSubscription = t.texModel.Subscribe(func(ch observer.Change) {
		t.deleteGraphicsRepresentation()
		t.createGraphicsRepresentation()
	})
}

func (t *TwoDTexture) unsubscribeFromModel() {
	t.texSubscription.Delete()
}

func (t *TwoDTexture) createGraphicsRepresentation() {
	definition := t.buildGraphicsDefinition(t.texModel)

	t.texture = t.engine.CreateTwoDTexture(definition)

	t.material = t.engine.CreatePBRMaterial(graphics.PBRMaterialDefinition{
		BackfaceCulling: false,
		AlphaBlending:   false,
		AlphaTesting:    false,
		Metalness:       0.0,
		Roughness:       0.5,
		AlbedoColor:     sprec.NewVec4(1.0, 1.0, 1.0, 1.0),
		AlbedoTexture:   t.texture,
	})

	quadCount := 5
	vertexSize := 3*4 + 3*4 + 2*4
	vertexData := make([]byte, 4*vertexSize*quadCount)
	vertexPlotter := buffer.NewPlotter(vertexData, binary.LittleEndian)

	renderQuad := func(vertexPlotter *buffer.Plotter, offset sprec.Vec3, texOffset sprec.Vec2) {
		twoDTextureVertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(-0.5, 0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(0.0, 1.0), texOffset),
		}.Serialize(vertexPlotter)
		twoDTextureVertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(-0.5, -0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(0.0, 0.0), texOffset),
		}.Serialize(vertexPlotter)
		twoDTextureVertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(0.5, -0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(1.0, 0.0), texOffset),
		}.Serialize(vertexPlotter)
		twoDTextureVertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(0.5, 0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(1.0, 1.0), texOffset),
		}.Serialize(vertexPlotter)
	}

	renderQuad(vertexPlotter, sprec.NewVec3(0.0, 0.0, 0.0), sprec.NewVec2(0.0, 0.0))
	renderQuad(vertexPlotter, sprec.NewVec3(0.0, 1.01, 0.0), sprec.NewVec2(0.0, 1.0))
	renderQuad(vertexPlotter, sprec.NewVec3(0.0, -1.01, 0.0), sprec.NewVec2(0.0, -1.0))
	renderQuad(vertexPlotter, sprec.NewVec3(-1.01, 0.0, 0.0), sprec.NewVec2(-1.0, 0.0))
	renderQuad(vertexPlotter, sprec.NewVec3(1.01, 0.0, 0.0), sprec.NewVec2(1.0, 0.0))

	indexData := make([]byte, 6*2*quadCount)
	indexPlotter := buffer.NewPlotter(indexData, binary.LittleEndian)
	for i := uint16(0); i < uint16(quadCount); i++ {
		indexPlotter.PlotUint16(0 + i*4)
		indexPlotter.PlotUint16(1 + i*4)
		indexPlotter.PlotUint16(2 + i*4)

		indexPlotter.PlotUint16(0 + i*4)
		indexPlotter.PlotUint16(2 + i*4)
		indexPlotter.PlotUint16(3 + i*4)
	}

	t.meshTemplate = t.engine.CreateMeshTemplate(graphics.MeshTemplateDefinition{
		VertexData: vertexData,
		VertexFormat: graphics.VertexFormat{
			HasCoord:            true,
			CoordOffsetBytes:    0,
			CoordStrideBytes:    vertexSize,
			HasNormal:           true,
			NormalOffsetBytes:   3 * 4,
			NormalStrideBytes:   vertexSize,
			HasTexCoord:         true,
			TexCoordOffsetBytes: 3*4 + 3*4,
			TexCoordStrideBytes: vertexSize,
		},
		IndexData:   indexData,
		IndexFormat: graphics.IndexFormatU16,
		SubMeshes: []graphics.SubMeshTemplateDefinition{
			{
				Primitive:   graphics.PrimitiveTriangles,
				IndexOffset: 0,
				IndexCount:  6 * quadCount,
				Material:    t.material,
			},
		},
	})

	t.mesh = t.scene.CreateMesh(t.meshTemplate)
}

func (t *TwoDTexture) deleteGraphicsRepresentation() {
	t.mesh.Delete()
	t.meshTemplate.Delete()
	t.material.Delete()
	t.texture.Delete()
}

func (t *TwoDTexture) buildGraphicsDefinition(src *model.TwoDTexture) graphics.TwoDTextureDefinition {
	return graphics.TwoDTextureDefinition{
		Width:           src.Width(),
		Height:          src.Height(),
		Wrapping:        assetToGraphicsWrap(src.Wrapping()),
		Filtering:       assetToGraphicsFilter(src.Filtering()),
		GenerateMipmaps: src.Mipmapping(),
		GammaCorrection: src.GammaCorrection(),
		InternalFormat:  assetToGraphicsInternalFormat(src.Format()),
		DataFormat:      assetToGraphicsDataFormat(src.Format()),
		Data:            src.Data(),
	}
}

func assetToGraphicsWrap(wrap asset.WrapMode) graphics.Wrap {
	switch wrap {
	case asset.WrapModeClampToEdge:
		return graphics.WrapClampToEdge
	case asset.WrapModeRepeat:
		return graphics.WrapRepeat
	case asset.WrapModeMirroredRepeat:
		return graphics.WrapMirroredRepat
	default:
		panic(fmt.Errorf("unsupported wrap: %v", wrap))
	}
}

func assetToGraphicsFilter(filter asset.FilterMode) graphics.Filter {
	switch filter {
	case asset.FilterModeNearest:
		return graphics.FilterNearest
	case asset.FilterModeLinear:
		return graphics.FilterLinear
	case asset.FilterModeAnisotropic:
		return graphics.FilterAnisotropic
	default:
		panic(fmt.Errorf("unsupported filter: %v", filter))
	}
}

func assetToGraphicsInternalFormat(format asset.TexelFormat) graphics.InternalFormat {
	switch format {
	case asset.TexelFormatRGBA8:
		return graphics.InternalFormatRGBA8
	case asset.TexelFormatRGBA16F:
		return graphics.InternalFormatRGBA16F
	case asset.TexelFormatRGBA32F:
		return graphics.InternalFormatRGBA32F
	default:
		panic(fmt.Errorf("unsupported format: %v", format))
	}
}

func assetToGraphicsDataFormat(format asset.TexelFormat) graphics.DataFormat {
	switch format {
	case asset.TexelFormatRGBA8:
		return graphics.DataFormatRGBA8
	case asset.TexelFormatRGBA16F:
		return graphics.DataFormatRGBA16F
	case asset.TexelFormatRGBA32F:
		return graphics.DataFormatRGBA32F
	default:
		panic(fmt.Errorf("unsupported format: %v", format))
	}
}

type twoDTextureVertex struct {
	Coord    sprec.Vec3
	TexCoord sprec.Vec2
}

func (v twoDTextureVertex) Serialize(plotter *buffer.Plotter) {
	plotter.PlotFloat32(v.Coord.X)
	plotter.PlotFloat32(v.Coord.Y)
	plotter.PlotFloat32(v.Coord.Z)
	plotter.PlotFloat32(0.0)
	plotter.PlotFloat32(0.0)
	plotter.PlotFloat32(1.0)
	plotter.PlotFloat32(v.TexCoord.X)
	plotter.PlotFloat32(v.TexCoord.Y)
}
