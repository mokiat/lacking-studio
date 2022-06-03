package model

import (
	"github.com/mokiat/lacking/ui/mvc"
	"golang.org/x/exp/slices"
)

var (
	ChangeStudio                  = mvc.NewChange("studio")
	ChangeStudioPropertiesVisible = mvc.SubChange(ChangeStudio, "properties_visible")
	ChangeStudioEditorAdded       = mvc.SubChange(ChangeStudio, "editor_added")
	ChangeStudioEditorRemoved     = mvc.SubChange(ChangeStudio, "editor_removed")
	ChangeStudioEditorSelection   = mvc.SubChange(ChangeStudio, "editor_selection")
)

func NewStudio() *Studio {
	return &Studio{
		Observable:   mvc.NewObservable(),
		dummyHistory: NewHistory(),
	}
}

type Studio struct {
	mvc.Observable

	isPropertiesVisible bool
	editors             []*Editor
	selectedEditor      *Editor
	dummyHistory        *History
}

func (s *Studio) IsPropertiesVisible() bool {
	return s.isPropertiesVisible
}

func (s *Studio) SetPropertiesVisible(visible bool) {
	s.isPropertiesVisible = visible
	s.SignalChange(ChangeStudioPropertiesVisible)
}

func (s *Studio) Editors() []*Editor {
	return s.editors
}

func (s *Studio) IterateEditors(cb func(*Editor)) {
	for _, editor := range s.editors {
		cb(editor)
	}
}

func (s *Studio) AddEditor(editor *Editor) {
	s.editors = append(s.editors, editor)
	s.SignalChange(ChangeStudioEditorAdded)
}

func (s *Studio) RemoveEditor(editor *Editor) {
	if index := slices.Index(s.editors, editor); index >= 0 {
		s.editors = slices.Delete(s.editors, index, index+1)
		if s.selectedEditor == editor {
			switch {
			case len(s.editors) == 0:
				s.selectedEditor = nil
			case index < len(s.editors):
				s.selectedEditor = s.editors[index]
			default:
				s.selectedEditor = s.editors[index-1]
			}
			s.SignalChange(ChangeStudioEditorSelection)
		}
		s.SignalChange(ChangeStudioEditorRemoved)
	}
}

func (s *Studio) FindEditorByID(id string) *Editor {
	index := slices.IndexFunc(s.editors, func(editor *Editor) bool {
		return editor.Resource().ID() == id
	})
	if index < 0 {
		return nil
	}
	return s.editors[index]
}

func (s *Studio) SelectedEditor() *Editor {
	return s.selectedEditor
}

func (s *Studio) SetSelectedEditor(editor *Editor) {
	s.selectedEditor = editor
	s.SignalChange(ChangeStudioEditorSelection)
}

func (s *Studio) SelectedHistory() *History {
	if s.selectedEditor == nil {
		return s.dummyHistory
	}
	return s.selectedEditor.History()
}
