package view

import (
	"github.com/mokiat/gblob"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/global"
	"github.com/mokiat/lacking-studio/internal/preview/model"
	"github.com/mokiat/lacking-studio/internal/viewport"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Viewport = mvc.EventListener(co.Define(&viewportComponent{}))

type ViewportData struct {
	AppModel *model.AppModel
	Resource *asset.Resource
}

type viewportComponent struct {
	co.BaseComponent

	appModel *model.AppModel
	resource *asset.Resource

	renderAPI render.API

	ambientTexture render.Texture

	gameEngine *game.Engine
	gameScene  *game.Scene

	commonData  *viewport.CommonData
	cameraGizmo *viewport.CameraGizmo

	currentResourceSet *game.ResourceSet
}

func (c *viewportComponent) OnCreate() {
	data := co.GetData[ViewportData](c.Properties())
	c.appModel = data.AppModel
	c.resource = data.Resource

	window := co.Window(c.Scope())
	c.renderAPI = window.RenderAPI()

	ambientData := make(gblob.LittleEndianBlock, 4*4)
	ambientData.SetFloat32(0, 1.0)
	ambientData.SetFloat32(4, 1.0)
	ambientData.SetFloat32(8, 1.5)
	ambientData.SetFloat32(12, 1.0)
	c.ambientTexture = c.renderAPI.CreateColorTextureCube(render.ColorTextureCubeInfo{
		Dimension:       1,
		GenerateMipmaps: false,
		GammaCorrection: false,
		Format:          render.DataFormatRGBA32F,
		FrontSideData:   ambientData,
		BackSideData:    ambientData,
		LeftSideData:    ambientData,
		RightSideData:   ambientData,
		TopSideData:     ambientData,
		BottomSideData:  ambientData,
	})

	ctx := co.TypedValue[*global.Context](c.Scope())
	c.commonData = ctx.CommonData
	c.gameEngine = ctx.GameEngine

	c.gameScene = c.gameEngine.CreateScene()
	gfxScene := c.gameScene.Graphics()

	gfxScene.CreateSky(graphics.SkyInfo{
		Definition: c.commonData.SkyDefinition(),
	})

	dirLightNode := c.gameScene.CreateDirectionalLight(game.DirectionalLightInfo{
		EmitColor:  opt.V(dprec.NewVec3(1.0, 1.0, 1.0)),
		CastShadow: opt.V(true),
	})
	dirLightNode.SetPosition(dprec.NewVec3(0.0, 20.0, 20.0))
	dirLightNode.SetRotation(dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()))

	c.gameScene.CreateAmbientLight(game.AmbientLightInfo{
		ReflectionTexture: c.ambientTexture,
		RefractionTexture: c.ambientTexture,
		OuterRadius:       opt.V(2000.0),
		InnerRadius:       opt.V(2000.0),
		CastShadow:        opt.V(false),
	})

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetExposure(1.0)
	gfxCamera.SetAutoExposure(false)
	gfxCamera.SetFoV(sprec.Degrees(60))
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)

	gridMeshDef := c.commonData.GridMeshDefinition()
	gridMesh := gfxScene.CreateMesh(graphics.MeshInfo{
		Definition: gridMeshDef,
	})
	gridMesh.SetMatrix(dprec.IdentityMat4())

	c.cameraGizmo = viewport.NewCameraGizmo(gfxCamera)

	c.currentResourceSet = c.gameEngine.CreateResourceSet()
	promise := c.currentResourceSet.OpenModelByID(c.resource.ID())
	promise.OnSuccess(func(modelDefinition *game.ModelDefinition) {
		co.Schedule(c.Scope(), func() {
			c.handleModelLoaded(modelDefinition)
		})
	})
	promise.OnError(func(err error) {
		co.Schedule(c.Scope(), func() {
			c.handleModelLoadError(err)
		})
	})
}

func (c *viewportComponent) OnDelete() {
	c.ambientTexture.Release()
	c.currentResourceSet.Delete()
	c.gameScene.Delete()
}

func (c *viewportComponent) Render() co.Instance {
	return co.New(std.Viewport, func() {
		co.WithData(std.ViewportData{
			API: c.renderAPI,
		})
		co.WithCallbackData(std.ViewportCallbackData{
			OnKeyboardEvent: c.handleViewportKeyboardEvent,
			OnMouseEvent:    c.handleViewportMouseEvent,
			OnRender:        c.handleViewportRender,
		})
	})
}

func (c *viewportComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.RefreshEvent:
		// TODO: Reload the resource
		c.Invalidate()
	}
}

func (c *viewportComponent) handleViewportKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	return c.cameraGizmo.OnKeyboardEvent(element, event)
}

func (c *viewportComponent) handleViewportMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return c.cameraGizmo.OnMouseEvent(element, event)
}

func (c *viewportComponent) handleViewportRender(framebuffer render.Framebuffer, size ui.Size) {
	c.gameEngine.Update()
	c.gameEngine.Render(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  uint32(size.Width),
		Height: uint32(size.Height),
	})
}

func (c *viewportComponent) handleModelLoaded(modelDefinition *game.ModelDefinition) {
	// TODO: Remove existing model, if there is one.
	// TODO: Relese old resources.

	// TODO: Track model
	model := c.gameScene.CreateModel(game.ModelInfo{
		Name:       "Model",
		Definition: modelDefinition,
		Position:   dprec.ZeroVec3(),
		Rotation:   dprec.IdentityQuat(),
		Scale:      dprec.NewVec3(1.0, 1.0, 1.0),
		IsDynamic:  false, // FIXME: Setting this to true kills large scenes
	})
	if len(model.Animations()) > 0 {
		animation := model.Animations()[0]
		c.gameScene.PlayAnimation(animation)
	}
	// TODO: Find camera and light nodes and attach indicator gizmos to them
	// from the common data.
}

func (c *viewportComponent) handleModelLoadError(err error) {
}
