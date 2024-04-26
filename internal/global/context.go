package global

import (
	"github.com/mokiat/lacking-studio/internal/viewport"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

type Context struct {
	EventBus   *mvc.EventBus
	Registry   *asset.Registry
	GameEngine *game.Engine
	CommonData *viewport.CommonData
}
