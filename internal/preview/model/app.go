package model

import (
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewAppModel(eventBus *mvc.EventBus, registry *asset.Registry) *AppModel {
	return &AppModel{
		eventBus: eventBus,
		registry: registry,
	}
}

type AppModel struct {
	eventBus *mvc.EventBus
	registry *asset.Registry

	selectedResource *asset.Resource
}

func (m *AppModel) SelectedResource() *asset.Resource {
	return m.selectedResource
}

func (m *AppModel) SetSelectedResource(resource *asset.Resource) {
	m.selectedResource = resource
	m.eventBus.Notify(SelectedResourceChangedEvent{})
}

// TODO: Return a promise.
func (m *AppModel) Refresh() {
	if m.selectedResource == nil {
		if err := m.registry.Reload(); err != nil {
			m.eventBus.Notify(RefreshErrorEvent{
				Err: err,
			})
			return
		}
	}
	m.eventBus.Notify(RefreshEvent{})
}

func (m *AppModel) Resources() []*asset.Resource {
	return m.registry.Resources()
}

type SelectedResourceChangedEvent struct{}

type RefreshEvent struct{}

type RefreshErrorEvent struct {
	Err error
}
