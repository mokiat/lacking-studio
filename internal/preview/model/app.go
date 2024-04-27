package model

import (
	"os"
	"os/exec"

	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/async"
)

func NewAppModel(window *ui.Window, eventBus *mvc.EventBus, registry *asset.Registry) *AppModel {
	return &AppModel{
		window:   window,
		eventBus: eventBus,
		registry: registry,

		cameraSectionExpanded: true,
		autoExposure:          false,

		sceneSectionExpanded: true,
		showGrid:             true,
		showAmbientLight:     true,
		showDirectionalLight: true,
		showSky:              true,

		refreshEnabled: true,
	}
}

type AppModel struct {
	window   *ui.Window
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

	refreshEnabled bool
}

func (m *AppModel) SelectedResource() *asset.Resource {
	return m.selectedResource
}

func (m *AppModel) SetSelectedResource(resource *asset.Resource) {
	m.selectedResource = resource
	m.eventBus.Notify(SelectedResourceChangedEvent{})
}

func (m *AppModel) RefreshEnabled() bool {
	return m.refreshEnabled
}

func (m *AppModel) Refresh() {
	if m.refreshEnabled {
		m.refreshEnabled = false
		var promise async.Promise[struct{}]
		if m.selectedResource == nil {
			promise = m.refreshRegistry()
		} else {
			promise = m.refreshResource(m.selectedResource)
		}
		promise.OnSuccess(func(struct{}) {
			m.window.Schedule(func() {
				m.refreshEnabled = true
				m.eventBus.Notify(RefreshEvent{})
			})
		})
		promise.OnError(func(err error) {
			m.window.Schedule(func() {
				m.refreshEnabled = true
				m.eventBus.Notify(RefreshErrorEvent{
					Err: err,
				})
			})
		})
	}
}

func (m *AppModel) refreshRegistry() async.Promise[struct{}] {
	promise := async.NewPromise[struct{}]()
	go func() {
		if err := m.packAssets(""); err != nil {
			promise.Fail(err)
		}

		reloadErr := make(chan error)
		m.window.Schedule(func() {
			reloadErr <- m.registry.Reload()
		})
		if err := <-reloadErr; err != nil {
			promise.Fail(err)
		}

		promise.Deliver(struct{}{})
	}()
	return promise
}

func (m *AppModel) refreshResource(resource *asset.Resource) async.Promise[struct{}] {
	promise := async.NewPromise[struct{}]()
	go func() {
		if err := m.packAssets(resource.Name()); err != nil {
			promise.Fail(err)
		}
		promise.Deliver(struct{}{})
	}()
	return promise
}

func (m *AppModel) packAssets(model string) error {
	args := []string{"pack"}
	if model != "" {
		args = append(args, "--", model)
	}
	cmd := exec.Command("task", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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
