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
	}
}

type Model struct {
	eventBus *mvc.EventBus
	asset    *registry.Asset

	id   string
	name string
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
