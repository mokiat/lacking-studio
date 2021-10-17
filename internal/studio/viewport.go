package studio

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var Viewport = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	controller := props.Data().(*ViewportController)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: controller,
		})
		co.WithLayoutData(props.LayoutData())
	})
}))

var _ ui.ElementMouseHandler = (*ViewportController)(nil)
var _ ui.ElementRenderHandler = (*ViewportController)(nil)

type ViewportController struct {
	co.Controller

	gfxEngine     graphics.Engine
	physicsEngine *physics.Engine
	ecsEngine     *ecs.Engine

	scene *Scene
}

func (c *ViewportController) Init() {
	c.gfxEngine.Create()
	c.scene = c.createScene()
}

func (c *ViewportController) Free() {
	c.gfxEngine.Destroy()
}

func (c *ViewportController) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return true
}

func (c *ViewportController) OnRender(element *ui.Element, canvas ui.Canvas) {
	canvas.DrawSurface(c.scene)
	element.Context().Window().Invalidate()
}

func (c *ViewportController) createScene() *Scene {
	scene := &Scene{
		gfxScene:     c.gfxEngine.CreateScene(),
		physicsScene: c.physicsEngine.CreateScene(0.015),
		ecsScene:     c.ecsEngine.CreateScene(),
	}
	scene.Init()
	return scene
}

type Scene struct {
	gfxScene     graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	camera graphics.Camera
}

func (s *Scene) Init() {
	s.camera = s.gfxScene.CreateCamera()
	s.camera.SetPosition(sprec.NewVec3(0.0, 0.0, 0.0))
	s.camera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	s.camera.SetFoV(sprec.Degrees(66))
	s.camera.SetAutoExposure(true)
	s.camera.SetExposure(1.0)
	s.camera.SetAutoFocus(false)

	s.gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(1.0, 0.0, 0.0))
}

func (s *Scene) Render(x, y, width, height int) {
	s.gfxScene.Render(graphics.Viewport{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}, s.camera)
}

func (s *Scene) Delete() {
	s.gfxScene.Delete()
	s.physicsScene.Delete()
	s.ecsScene.Delete()
}
