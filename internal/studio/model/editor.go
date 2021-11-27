package model

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type Editor interface {
	ID() string
	Name() string
	Icon() ui.Image

	CanUndo() bool
	Undo()
	CanRedo() bool
	Redo()
	CanSave() bool
	Save() error

	Render(layoutData mat.LayoutData) co.Instance
}
