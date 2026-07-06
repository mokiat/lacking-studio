package global

import (
	"github.com/mokiat/lacking-studio/internal/viewport"
	"github.com/mokiat/lacking/core/resource"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui/mvc"
)

type Context struct {
	EventBus   *mvc.EventBus
	Store      resource.Store
	GameEngine *game.Engine
	CommonData *viewport.CommonData
}
