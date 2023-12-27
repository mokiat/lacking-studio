package registry

import "github.com/mokiat/lacking/ui/mvc"

func NewModel(eventBus *mvc.EventBus) *Model {
	return &Model{
		eventBus: eventBus,
	}
}

type Model struct {
	eventBus *mvc.EventBus
}

func (m *Model) Assets() []*Asset {
	return []*Asset{
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		// TODO
	}
}
