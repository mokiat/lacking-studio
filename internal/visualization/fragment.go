package visualization

import (
	"image"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/view/editor/viewport"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/blob"
)

type TestUniform struct {
	BaseColor   sprec.Vec4
	GlowColor   sprec.Vec4
	SplitHeight float32
}

func (u TestUniform) Std140Plot(plotter *blob.Plotter) {
	plotter.PlotSPVec4(u.BaseColor)
	plotter.PlotSPVec4(u.GlowColor)
	plotter.PlotFloat32(u.SplitHeight)
}

func (u TestUniform) Std140Size() int {
	return 4*render.SizeF32 + 4*render.SizeF32 + render.SizeF32
}

func NewFragment(renderAPI render.API, gfxEngine *graphics.Engine, commonData *viewport.CommonData) *Fragment {
	gfxScene := gfxEngine.CreateScene()

	// gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.01, 0.01, 0.02))

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetExposure(1.0)
	gfxCamera.SetAutoExposure(false)
	gfxCamera.SetFoV(sprec.Degrees(60))
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	cameraGizmo := viewport.NewCameraGizmo(gfxCamera)

	gridMeshDef := commonData.GridMeshDefinition()
	gridMesh := gfxScene.CreateMesh(graphics.MeshInfo{
		Definition: gridMeshDef,
	})
	gridMesh.SetMatrix(dprec.IdentityMat4())

	testShader := gfxEngine.CreateForwardShader(graphics.ShaderInfo{
		SourceCode: `
			uniforms {
				baseColor vec4,
				glowColor vec4,
			}

			func #fragment() {
				#color = glowColor
			}	
		`,
	})

	testMaterial := gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name:           "Test",
		GeometryPasses: []graphics.GeometryRenderPassInfo{},
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Culling: opt.V(render.CullModeNone),
				Shader:  testShader,
			},
		},
	})
	testMaterial.SetProperty("baseColor", sprec.NewVec4(1.0, 0.0, 0.0, 1.0))
	testMaterial.SetProperty("glowColor", sprec.NewVec4(0.0, 10.0, 0.0, 1.0))

	testMeshBuilder := graphics.NewShapeBuilder()
	testMeshBuilder.Solid(testMaterial).Sphere(sprec.ZeroVec3(), 1.0, 40)

	testMeshGeometry := gfxEngine.CreateMeshGeometry(testMeshBuilder.BuildGeometryInfo())
	testMeshDefinition := gfxEngine.CreateMeshDefinition(testMeshBuilder.BuildMeshDefinitionInfo(testMeshGeometry))

	testMesh := gfxScene.CreateMesh(graphics.MeshInfo{
		Definition: testMeshDefinition,
	})
	testMesh.SetMatrix(dprec.IdentityMat4())

	return &Fragment{
		renderAPI:  renderAPI,
		commonData: commonData,

		gfxEngine: gfxEngine,
		gfxScene:  gfxScene,
		gfxCamera: gfxCamera,

		cameraGizmo: cameraGizmo,
	}
}

type Fragment struct {
	renderAPI  render.API
	commonData *viewport.CommonData

	gfxEngine *graphics.Engine
	gfxScene  *graphics.Scene
	gfxCamera *graphics.Camera

	cameraGizmo *viewport.CameraGizmo
}

func (f *Fragment) Delete() {
	f.gfxScene.Delete()
}

func (f *Fragment) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	return f.cameraGizmo.OnKeyboardEvent(element, event)
}

func (f *Fragment) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	// TODO: Do camera motion. Have a "Gadget" concept and pass control initially to it.
	// Then trickle down until you get to here. If no gadget is interested, then do camera motion.

	return f.cameraGizmo.OnMouseEvent(element, event)
}

func (f *Fragment) OnRender(framebuffer render.Framebuffer, size ui.Size) {
	f.gfxScene.RenderFramebuffer(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  uint32(size.Width),
		Height: uint32(size.Height),
	})
}

func (f *Fragment) TakeSnapshot(size ui.Size) image.Image {
	colorTexture := f.renderAPI.CreateColorTexture2D(render.ColorTexture2DInfo{
		Width:           uint32(size.Width),
		Height:          uint32(size.Height),
		GammaCorrection: false,
		Format:          render.DataFormatRGBA8,
	})
	defer colorTexture.Release()

	framebuffer := f.renderAPI.CreateFramebuffer(render.FramebufferInfo{
		ColorAttachments: [4]render.Texture{
			colorTexture,
		},
	})
	defer framebuffer.Release()

	buffer := f.renderAPI.CreatePixelTransferBuffer(render.BufferInfo{
		Size: uint32(4 * size.Width * size.Height),
	})
	defer buffer.Release()

	f.gfxScene.RenderFramebuffer(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  uint32(size.Width),
		Height: uint32(size.Height),
	})

	commands := f.renderAPI.CreateCommandBuffer(1024)
	commands.BeginRenderPass(render.RenderPassInfo{
		Framebuffer: framebuffer,
		Viewport: render.Area{
			X:      0,
			Y:      0,
			Width:  uint32(size.Width),
			Height: uint32(size.Height),
		},
		DepthLoadOp:    render.LoadOperationLoad,
		DepthStoreOp:   render.StoreOperationDiscard,
		StencilLoadOp:  render.LoadOperationLoad,
		StencilStoreOp: render.StoreOperationDiscard,
		Colors: [4]render.ColorAttachmentInfo{
			{
				LoadOp: render.LoadOperationLoad,
			},
		},
	})
	commands.CopyFramebufferToBuffer(render.CopyFramebufferToBufferInfo{
		Buffer: buffer,
		Offset: 0,
		X:      0,
		Y:      0,
		Width:  uint32(size.Width),
		Height: uint32(size.Height),
		Format: render.DataFormatRGBA8,
	})
	commands.EndRenderPass()
	f.renderAPI.Queue().Submit(commands)

	previewImg := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))
	f.renderAPI.Queue().ReadBuffer(buffer, 0, previewImg.Pix)

	// Flip image.
	for y := 0; y < size.Height/2; y++ {
		topOffset := y * (4 * size.Width)
		bottomOffset := (size.Height - y - 1) * (4 * size.Width)
		for x := 0; x < size.Width*4; x++ {
			previewImg.Pix[topOffset+x], previewImg.Pix[bottomOffset+x] =
				previewImg.Pix[bottomOffset+x], previewImg.Pix[topOffset+x]
		}
	}
	return previewImg
}

func (f *Fragment) CreatePointLight() *PointLight {
	mesh := f.gfxScene.CreateMesh(graphics.MeshInfo{
		Definition: f.commonData.PointLightMeshDefinition(),
	})
	light := f.gfxScene.CreatePointLight(graphics.PointLightInfo{
		Position:  dprec.ZeroVec3(),             // TODO
		EmitRange: 100.0,                        // TODO
		EmitColor: dprec.NewVec3(1.0, 1.0, 1.0), // TODO
	})
	return &PointLight{
		mesh:  mesh,
		light: light,
	}
}
