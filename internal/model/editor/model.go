package editor

import (
	"github.com/google/uuid"
	"github.com/mokiat/lacking-studio/internal/model/registry"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewModel(eventBus *mvc.EventBus, name string) *Model {
	return &Model{
		eventBus: eventBus,
		asset:    nil,

		id:   uuid.NewString(),
		name: name,

		navigatorPage: NavigatorPageNodes,
		inspectorPage: InspectorPageAsset,
	}
}

type Model struct {
	eventBus *mvc.EventBus
	asset    *registry.Asset

	id   string
	name string

	navigatorPage NavigatorPage
	inspectorPage InspectorPage
}

func (m *Model) ID() string {
	return m.id
}

func (m *Model) Name() string {
	return m.name
}

func (m *Model) Image() *ui.Image {
	return nil
}

func (m *Model) Asset() *registry.Asset {
	return m.asset
}

func (m *Model) CanSave() bool {
	return false
}

func (m *Model) NavigatorPage() NavigatorPage {
	return m.navigatorPage
}

func (m *Model) SetNavigatorPage(page NavigatorPage) {
	if page != m.navigatorPage {
		m.navigatorPage = page
		m.eventBus.Notify(NavigatorPageChangedEvent{
			Editor: m,
		})
	}
}

func (m *Model) InspectorPage() InspectorPage {
	return m.inspectorPage
}

func (m *Model) SetInspectorPage(page InspectorPage) {
	if page != m.inspectorPage {
		m.inspectorPage = page
		m.eventBus.Notify(InspectorPageChangedEvent{
			Editor: m,
		})
	}
}
