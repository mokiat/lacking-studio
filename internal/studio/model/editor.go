package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeEditor                  = mvc.NewChange("editor")
	ChangeEditorPropertiesVisible = mvc.SubChange(ChangeEditor, "properties_visible")
)

func NewEditor(resource *Resource) *Editor {
	return &Editor{
		Observable:          mvc.NewObservable(),
		history:             NewHistory(),
		resource:            resource,
		isPropertiesVisible: true,
	}
}

type Editor struct {
	mvc.Observable
	history             *History
	resource            *Resource
	isPropertiesVisible bool
}

func (e *Editor) IsPropertiesVisible() bool {
	return e.isPropertiesVisible
}

func (e *Editor) SetPropertiesVisible(visible bool) {
	e.isPropertiesVisible = visible
	e.SignalChange(ChangeEditorPropertiesVisible)
}

func (e *Editor) History() *History {
	return e.history
}

func (e *Editor) Resource() *Resource {
	return e.resource
}
