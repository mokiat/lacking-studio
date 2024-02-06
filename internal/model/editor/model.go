package editor

import (
	registrymodel "github.com/mokiat/lacking-studio/internal/model/registry"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewModel(eventBus *mvc.EventBus, asset *registrymodel.Asset) *Model {
	return &Model{
		eventBus: eventBus,
		asset:    asset,

		navigatorPage: NavigatorPageNodes,
		inspectorPage: InspectorPageAsset,
	}
}

type Model struct {
	eventBus *mvc.EventBus
	asset    *registrymodel.Asset

	navigatorPage NavigatorPage
	inspectorPage InspectorPage

	textures []*Texture

	selection any
}

func (m *Model) ID() string {
	return m.asset.ID()
}

func (m *Model) Name() string {
	return m.asset.Name()
}

func (m *Model) Image() *ui.Image {
	return m.asset.Image()
}

func (m *Model) Asset() *registrymodel.Asset {
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

func (m *Model) Textures() []*Texture {
	return m.textures
}

func (m *Model) Selection() any {
	return m.selection
}

func (m *Model) SetSelection(selection any) {
	if selection != m.selection {
		m.selection = selection
		m.eventBus.Notify(SelectionChangedEvent{
			Editor: m,
		})
	}
}
