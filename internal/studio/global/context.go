package global

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
)

type Context struct {
	API            render.API
	GraphicsEngine *graphics.Engine
}
