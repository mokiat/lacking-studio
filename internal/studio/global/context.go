package global

import (
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
)

type Context struct {
	Window         *ui.Window
	API            render.API
	Registry       *data.Registry
	GraphicsEngine *graphics.Engine
	PhysicsEngine  *physics.Engine
	ECSEngine      *ecs.Engine
}
