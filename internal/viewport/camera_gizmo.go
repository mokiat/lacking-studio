package viewport

import (
	"math"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
)

func NewCameraGizmo(camera *graphics.Camera) *CameraGizmo {
	gizmo := &CameraGizmo{
		camera:   camera,
		position: dprec.ZeroVec3(),
		yaw:      dprec.Degrees(15),
		pitch:    dprec.Degrees(30),
		zoom:     3,
	}
	gizmo.updateCamera()
	return gizmo
}

type CameraGizmo struct {
	camera *graphics.Camera

	position dprec.Vec3
	yaw      dprec.Angle
	pitch    dprec.Angle
	zoom     float64

	oldMouseX float64
	oldMouseY float64
	wheelDown bool
	shiftDown bool
}

func (g *CameraGizmo) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {
	case ui.KeyCodeLeftShift:
		g.shiftDown = (event.Action != ui.KeyboardActionUp)
		return true
	}

	return false
}

func (g *CameraGizmo) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	switch event.Action {
	case ui.MouseActionDown:
		if event.Button == ui.MouseButtonMiddle {
			g.wheelDown = true
			return true
		}
		return false

	case ui.MouseActionUp:
		if event.Button == ui.MouseButtonMiddle {
			g.wheelDown = false
			return true
		}
		return false

	case ui.MouseActionMove:
		newMouseX := float64(event.X)
		newMouseY := float64(event.Y)
		deltaMouseX := newMouseX - g.oldMouseX
		deltaMouseY := newMouseY - g.oldMouseY
		g.oldMouseX = newMouseX
		g.oldMouseY = newMouseY

		switch {
		case g.wheelDown && g.shiftDown:
			g.handlePan(deltaMouseX, -deltaMouseY)
			element.Invalidate()
			return true
		case g.wheelDown:
			g.handleRotation(deltaMouseX, deltaMouseY)
			element.Invalidate()
			return true
		default:
			return false
		}

	case ui.MouseActionScroll:
		g.handleZoom(-float64(event.ScrollY) * 0.01)
		element.Invalidate()
		return true
	}

	return false
}

func (g *CameraGizmo) handlePan(deltaX, deltaY float64) {
	matrix := g.cameraMatrix()
	vecX := matrix.OrientationX()
	vecY := matrix.OrientationY()

	translationAmount := math.Pow(2.0, g.zoom) / 300.0

	g.position = dprec.Vec3Diff(
		g.position,
		dprec.Vec3Sum(
			dprec.Vec3Prod(vecX, deltaX*translationAmount),
			dprec.Vec3Prod(vecY, deltaY*translationAmount),
		),
	)
	g.updateCamera()
}

func (g *CameraGizmo) handleRotation(deltaX, deltaY float64) {
	rotationAmount := dprec.Degrees(0.4)

	g.yaw -= rotationAmount * dprec.Angle(deltaX)
	g.pitch += rotationAmount * dprec.Angle(deltaY)
	g.updateCamera()

}

func (g *CameraGizmo) handleZoom(delta float64) {
	g.zoom += delta
	g.updateCamera()
}

func (g *CameraGizmo) cameraMatrix() dprec.Mat4 {
	return dprec.Mat4MultiProd(
		dprec.TranslationMat4(g.position.X, g.position.Y, g.position.Z),
		dprec.RotationMat4(g.yaw, 0.0, 1.0, 0.0),
		dprec.RotationMat4(g.pitch, -1.0, 0.0, 0.0),
		dprec.TranslationMat4(0.0, 0.0, math.Pow(2.0, g.zoom)),
	)
}

func (g *CameraGizmo) updateCamera() {
	g.camera.SetMatrix(g.cameraMatrix())
}
