package visualization

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game/graphics"
)

type PointLight struct {
	mesh  *graphics.Mesh
	light *graphics.PointLight
}

func (l *PointLight) Delete() {
	defer l.mesh.Delete()
	defer l.light.Delete()
}

func (l *PointLight) SetMatrix(matrix dprec.Mat4) {
	l.mesh.SetMatrix(matrix)
	l.light.SetPosition(matrix.Translation())
}
