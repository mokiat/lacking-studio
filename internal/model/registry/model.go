package registry

import (
	"github.com/mokiat/gog"
	asset "github.com/mokiat/lacking/game/newasset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewModel(eventBus *mvc.EventBus, context *ui.Context, delegate *asset.Registry) *Model {
	assets := gog.Map(delegate.Resources(), func(resource *asset.Resource) *Asset {
		return &Asset{
			delegate: resource,
		}
	})

	return &Model{
		eventBus: eventBus,
		context:  context,
		delegate: delegate,

		assets: assets,
	}
}

type Model struct {
	eventBus *mvc.EventBus
	context  *ui.Context
	delegate *asset.Registry

	assets []*Asset
}

func (m *Model) Assets() []*Asset {
	return m.assets
}
