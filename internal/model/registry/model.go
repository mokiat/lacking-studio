package registry

import (
	"fmt"
	"slices"

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

func (m *Model) RenameAsset(asset *Asset, name string) error {
	resource := asset.delegate
	if err := resource.SetName(name); err != nil {
		return fmt.Errorf("error renaming resource: %w", err)
	}
	m.eventBus.Notify(AssetsChangedEvent{})
	return nil
}

func (m *Model) CloneAsset(asset *Asset) (*Asset, error) {
	newName := fmt.Sprintf("%s (clone)", asset.delegate.Name())
	newResource, err := m.delegate.CreateResource(newName)
	if err != nil {
		return nil, fmt.Errorf("error creating resource: %w", err)
	}
	if err := newResource.SetPreview(asset.delegate.Preview()); err != nil {
		return nil, fmt.Errorf("error setting preview: %w", err)
	}
	content, err := asset.delegate.OpenContent()
	if err != nil {
		return nil, fmt.Errorf("error opening content: %w", err)
	}
	if err := newResource.SaveContent(content); err != nil {
		return nil, fmt.Errorf("error saving content: %w", err)
	}
	clonedAsset := resourceToAsset(m.context, newResource)
	m.assets = append(m.assets, clonedAsset)
	m.eventBus.Notify(AssetsChangedEvent{})
	return clonedAsset, nil
}

func (m *Model) DeleteAsset(asset *Asset) error {
	resource := asset.delegate
	if err := resource.Delete(); err != nil {
		return fmt.Errorf("error deleting resource: %w", err)
	}
	m.assets = slices.DeleteFunc(m.assets, func(candidate *Asset) bool {
		return candidate == asset
	})
	m.eventBus.Notify(AssetsChangedEvent{})
	return nil
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
