package model

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
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

type IEditor interface {
	mvc.Reducer
	Save() error
	Render(scope co.Scope, layoutData mat.LayoutData) co.Instance
	Destroy()
}
