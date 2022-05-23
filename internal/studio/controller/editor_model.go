package controller

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewModelEditor(studio *Studio, resource *data.Resource) (*ModelEditor, error) {
	gfxScene := studio.GraphicsEngine().CreateScene()
	gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.1, 0.3, 0.5))

	dirLight := gfxScene.CreateDirectionalLight()
	dirLight.SetIntensity(sprec.NewVec3(1.0, 1.0, 1.0))
	dirLight.SetRotation(sprec.IdentityQuat())

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetPosition(sprec.NewVec3(0.0, 0.0, 3.0))
	gfxCamera.SetRotation(sprec.IdentityQuat())
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	gfxCamera.SetFoV(sprec.Degrees(66))
	gfxCamera.SetAutoExposure(false)
	gfxCamera.SetExposure(3.14)
	gfxCamera.SetAutoFocus(false)

	return &ModelEditor{
		Controller: co.NewBaseController(),
		BaseEditor: NewBaseEditor(),

		studio:   studio,
		resource: resource,

		propsAssetExpanded: false,

		gfxEngine:      studio.GraphicsEngine(),
		gfxScene:       gfxScene,
		gfxCamera:      gfxCamera,
		gfxCameraPitch: sprec.Degrees(0),
		gfxCameraYaw:   sprec.Degrees(0),
		gfxCameraFoV:   sprec.Degrees(66),
	}, nil
}

var _ model.Editor = (*ModelEditor)(nil)
var _ model.ModelEditor = (*ModelEditor)(nil)

type ModelEditor struct {
	co.Controller
	BaseEditor

	studio      *Studio
	resource    *data.Resource
	savedChange history.Change

	propsAssetExpanded bool

	gfxEngine       *graphics.Engine
	gfxScene        *graphics.Scene
	gfxCamera       *graphics.Camera
	gfxCameraPitch  sprec.Angle
	gfxCameraYaw    sprec.Angle
	gfxCameraFoV    sprec.Angle
	gfxMesh         *graphics.Mesh
	gfxMeshTemplate *graphics.MeshTemplate
	gfxMaterial     *graphics.Material
	gfxImage        *graphics.TwoDTexture

	assetModel asset.Model

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
}

func (e *ModelEditor) API() render.API {
	return e.studio.api
}

func (e *ModelEditor) IsPropertiesVisible() bool {
	return e.studio.IsPropertiesVisible()
}

func (e *ModelEditor) ID() string {
	return e.resource.ID()
}

func (e *ModelEditor) Name() string {
	return e.resource.Name()
}

func (e *ModelEditor) Icon() *ui.Image {
	return co.OpenImage("resources/icons/model.png")
}

func (e *ModelEditor) CanSave() bool {
	return e.savedChange != e.changes.LastChange()
}

func (e *ModelEditor) Save() error {
	return nil
}

func (e *ModelEditor) Update() {
	transform := sprec.Mat4MultiProd(
		sprec.RotationMat4(-e.gfxCameraYaw, 0.0, 1.0, 0.0),
		sprec.RotationMat4(-e.gfxCameraPitch, 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, 3.0),
	)
	e.gfxCamera.SetPosition(transform.Translation())
	e.gfxCamera.SetRotation(transform.RotationQuat())
	e.gfxCamera.SetFoV(e.gfxCameraFoV)
}

func (e *ModelEditor) OnViewportMouseEvent(event mat.ViewportMouseEvent) bool {
	switch event.Type {
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = true
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
		return true
	case ui.MouseEventTypeMove:
		if e.rotatingCamera {
			e.gfxCameraPitch += sprec.Degrees(float32(event.Position.Y-e.oldMouseY) / 5)
			e.gfxCameraYaw += sprec.Degrees(float32(event.Position.X-e.oldMouseX) / 5)
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
		return true
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = false
		}
		return true
	case ui.MouseEventTypeScroll:
		fov := e.gfxCameraFoV.Degrees()
		fov -= 2 * float32(event.ScrollY)
		fov = sprec.Clamp(fov, 0.1, 179.0)
		e.gfxCameraFoV = sprec.Degrees(fov)
		return true
	default:
		return false
	}
}

func (e *ModelEditor) Scene() *graphics.Scene {
	return nil
}

func (e *ModelEditor) Camera() *graphics.Camera {
	return nil
}

func (e *ModelEditor) IsAssetAccordionExpanded() bool {
	return e.propsAssetExpanded
}

func (e *ModelEditor) SetAssetAccordionExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.NotifyChanged()
}

func (e *ModelEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(view.Model, func() {
		co.WithData(e)
		co.WithLayoutData(layoutData)
	})
}

func (e *ModelEditor) Destroy() {
	e.gfxScene.Delete()
	if e.gfxMesh != nil {
		e.gfxMesh.Delete()
	}
	if e.gfxMeshTemplate != nil {
		e.gfxMeshTemplate.Delete()
	}
	if e.gfxMaterial != nil {
		e.gfxMaterial.Delete()
	}
	if e.gfxImage != nil {
		e.gfxImage.Delete()
	}
}
