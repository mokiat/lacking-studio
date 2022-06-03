package model

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewEditor(resource *Resource) *Editor {
	return &Editor{
		history:  NewHistory(),
		resource: resource,
	}
}

type Editor struct {
	history  *History
	resource *Resource
	handler  IEditor
}

func (e *Editor) History() *History {
	return e.history
}

func (e *Editor) Resource() *Resource {
	return e.resource
}

// // TODO: Remove
// func (e *Editor) SetHandler(handler IEditor) {
// 	e.handler = handler
// }

// // TODO: Remove
// func (e *Editor) Handler() IEditor {
// 	return e.handler
// }

type IEditor interface {
	// ID() string
	// Name() string
	// Icon(scope co.Scope) *ui.Image

	Save() error

	Render(scope co.Scope, layoutData mat.LayoutData) co.Instance

	Destroy()
}
