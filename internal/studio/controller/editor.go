package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

const undoCount = 10

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

	Destroy()
}

func NewBaseEditor() BaseEditor {
	return BaseEditor{
		Controller: co.NewBaseController(),
		changes:    history.NewQueue(undoCount),
	}
}

type BaseEditor struct {
	co.Controller
	changes *history.Queue
}

func (e BaseEditor) CanUndo() bool {
	return e.changes.CanPop()
}

func (e BaseEditor) Undo() {
	if err := e.changes.Pop(); err != nil {
		panic(err)
	}
}

func (e BaseEditor) CanRedo() bool {
	return e.changes.CanUnpop()
}

func (e BaseEditor) Redo() {
	if err := e.changes.Unpop(); err != nil {
		panic(err)
	}
}

func (e BaseEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(mat.Element, func() {
		co.WithLayoutData(layoutData)
	})
}
