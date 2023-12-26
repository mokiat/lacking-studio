package app

import "github.com/mokiat/lacking/ui/mvc"

func NewModel(eventBus *mvc.EventBus) *Model {
	return &Model{
		eventBus: eventBus,

		navigatorVisible: true,
		inspectorVisible: true,
	}
}

type Model struct {
	eventBus *mvc.EventBus

	navigatorVisible bool
	inspectorVisible bool
}

func (m *Model) IsNavigatorVisible() bool {
	return m.navigatorVisible
}

func (m *Model) SetNavigatorVisible(visible bool) {
	if visible != m.navigatorVisible {
		m.navigatorVisible = visible
		m.eventBus.Notify(NavigatorVisibleChangedEvent{})
	}
}

func (m *Model) IsInspectorVisible() bool {
	return m.inspectorVisible
}

func (m *Model) SetInpsectorVisible(visible bool) {
	if visible != m.inspectorVisible {
		m.inspectorVisible = visible
		m.eventBus.Notify(InspectorVisibleChangedEvent{})
	}
}
