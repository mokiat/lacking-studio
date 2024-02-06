package registry

import (
	"fmt"

	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/debug/log"
	asset "github.com/mokiat/lacking/game/newasset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewModel(eventBus *mvc.EventBus, context *ui.Context, delegate *asset.Registry) *Model {
	assets := gog.Map(delegate.Resources(), func(resource *asset.Resource) *Asset {
		return resourceToAsset(context, resource)
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

func (m *Model) CreateAsset(name string) (*Asset, error) {
	resource, err := m.delegate.CreateResource(name)
	if err != nil {
		return nil, fmt.Errorf("error creating resource: %w", err)
	}
	asset := resourceToAsset(m.context, resource)
	m.assets = append(m.assets, asset)
	m.eventBus.Notify(AssetsChangedEvent{})
	return asset, nil
}

func resourceToAsset(context *ui.Context, resource *asset.Resource) *Asset {
	var previewImage *ui.Image
	if previewImg := resource.Preview(); previewImg != nil {
		var err error
		previewImage, err = context.CreateImage(previewImg)
		if err != nil {
			log.Error("Failed to create preview image: %v", err)
			previewImage = nil
		}
	}
	return &Asset{
		delegate:     resource,
		previewImage: previewImage,
	}
}

type AssetsChangedEvent struct{}
