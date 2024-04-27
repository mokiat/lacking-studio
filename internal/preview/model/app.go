package model

import (
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewAppModel(eventBus *mvc.EventBus, registry *asset.Registry) *AppModel {
	return &AppModel{
		eventBus: eventBus,
		registry: registry,

		cameraSectionExpanded: true,
		autoExposure:          false,

		sceneSectionExpanded: true,
		showGrid:             true,
		showAmbientLight:     true,
		showDirectionalLight: true,
		showSky:              true,
	}
}

type AppModel struct {
	eventBus *mvc.EventBus
	registry *asset.Registry

	selectedResource *asset.Resource

	cameraSectionExpanded bool
	autoExposure          bool

	sceneSectionExpanded bool
	showGrid             bool
	showAmbientLight     bool
	showDirectionalLight bool
	showSky              bool
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

func (m *AppModel) CameraSectionExpanded() bool {
	return m.cameraSectionExpanded
}

func (m *AppModel) SetCameraSectionExpanded(value bool) {
	if value != m.cameraSectionExpanded {
		m.cameraSectionExpanded = value
		m.eventBus.Notify(CameraSectionExpandedChangedEvent{})
	}
}

func (m *AppModel) AutoExposure() bool {
	return m.autoExposure
}

func (m *AppModel) SetAutoExposure(value bool) {
	if value != m.autoExposure {
		m.autoExposure = value
		m.eventBus.Notify(AutoExposureChangedEvent{})
	}
}

func (m *AppModel) SceneSectionExpanded() bool {
	return m.sceneSectionExpanded
}

func (m *AppModel) SetSceneSectionExpanded(value bool) {
	if value != m.sceneSectionExpanded {
		m.sceneSectionExpanded = value
		m.eventBus.Notify(SceneSectionExpandedChangedEvent{})
	}
}

func (m *AppModel) ShowGrid() bool {
	return m.showGrid
}

func (m *AppModel) SetShowGrid(value bool) {
	if value != m.showGrid {
		m.showGrid = value
		m.eventBus.Notify(ShowGridChangedEvent{})
	}
}

func (m *AppModel) ShowAmbientLight() bool {
	return m.showAmbientLight
}

func (m *AppModel) SetShowAmbientLight(value bool) {
	if value != m.showAmbientLight {
		m.showAmbientLight = value
		m.eventBus.Notify(ShowAmbientLightChangedEvent{})
	}
}

func (m *AppModel) ShowDirectionalLight() bool {
	return m.showDirectionalLight
}

func (m *AppModel) SetShowDirectionalLight(value bool) {
	if value != m.showDirectionalLight {
		m.showDirectionalLight = value
		m.eventBus.Notify(ShowDirectionalLightChangedEvent{})
	}
}

func (m *AppModel) ShowSky() bool {
	return m.showSky
}

func (m *AppModel) SetShowSky(value bool) {
	if value != m.showSky {
		m.showSky = value
		m.eventBus.Notify(ShowSkyChangedEvent{})
	}
}

type SelectedResourceChangedEvent struct{}

type RefreshEvent struct{}

type RefreshErrorEvent struct {
	Err error
}

type CameraSectionExpandedChangedEvent struct{}

type AutoExposureChangedEvent struct{}

type SceneSectionExpandedChangedEvent struct{}

type ShowGridChangedEvent struct{}

type ShowAmbientLightChangedEvent struct{}

type ShowDirectionalLightChangedEvent struct{}

type ShowSkyChangedEvent struct{}
