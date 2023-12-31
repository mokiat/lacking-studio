package app

import (
	"slices"
	"strings"

	"github.com/mokiat/gog"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking/ui/mvc"
)

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

	editors      []*editormodel.Model
	activeEditor *editormodel.Model
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

func (m *Model) Editors() []*editormodel.Model {
	return m.editors
}

func (m *Model) EachEditor(cb func(editor *editormodel.Model)) {
	for _, editor := range m.editors {
		cb(editor)
	}
}

func (m *Model) HasEditorWithName(name string) bool {
	_, has := gog.FindFunc(m.editors, func(editor *editormodel.Model) bool {
		return strings.EqualFold(editor.Name(), name)
	})
	return has
}

func (m *Model) AddEditor(editor *editormodel.Model) {
	if index := slices.Index(m.editors, editor); index >= 0 {
		return
	}
	m.editors = append(m.editors, editor)
	m.eventBus.Notify(EditorsChangedEvent{})
}

func (m *Model) RemoveEditor(editor *editormodel.Model) {
	index := slices.Index(m.editors, editor)
	if index < 0 {
		return
	}

	if m.activeEditor == editor {
		if index > 0 {
			m.SetActiveEditor(m.editors[index-1])
		} else if len(m.editors) > 1 {
			m.SetActiveEditor(m.editors[1])
		} else {
			m.SetActiveEditor(nil)
		}
	}

	m.editors = slices.DeleteFunc(m.editors, func(candidate *editormodel.Model) bool {
		return candidate == editor
	})
	m.eventBus.Notify(EditorsChangedEvent{})
}

func (m *Model) ActiveEditor() *editormodel.Model {
	return m.activeEditor
}

func (m *Model) SetActiveEditor(editor *editormodel.Model) {
	if editor != m.activeEditor {
		m.activeEditor = editor
		m.eventBus.Notify(ActiveEditorChangedEvent{})
	}
}
