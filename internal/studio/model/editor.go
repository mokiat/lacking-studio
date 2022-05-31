package model

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type Editor interface {
	ID() string
	Name() string
	Icon(scope co.Scope) *ui.Image

	CanUndo() bool
	Undo()
	CanRedo() bool
	Redo()
	CanSave() bool
	Save() error

	Render(scope co.Scope, layoutData mat.LayoutData) co.Instance

	Destroy()
}
