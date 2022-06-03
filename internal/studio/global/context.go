package global

import (
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking/render"
)

type Context struct {
	API      render.API
	Registry *data.Registry
}
