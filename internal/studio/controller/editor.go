package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

const undoCount = 10

func NewBaseEditor() BaseEditor {
	return BaseEditor{
		changes: history.NewQueue(undoCount),
	}
}

type BaseEditor struct {
	changes     *history.Queue
	savedChange history.Change
}

func (e *BaseEditor) CanUndo() bool {
	return e.changes.CanPop()
}

func (e *BaseEditor) Undo() {
	if err := e.changes.Pop(); err != nil {
		panic(err)
	}
}

func (e *BaseEditor) CanRedo() bool {
	return e.changes.CanUnpop()
}

func (e *BaseEditor) Redo() {
	if err := e.changes.Unpop(); err != nil {
		panic(err)
	}
}

func (e *BaseEditor) CanSave() bool {
	return e.savedChange != e.changes.LastChange()
}

func (e *BaseEditor) Save() error {
	e.savedChange = e.changes.LastChange()
	return nil
}

func (e *BaseEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(mat.Element, func() {
		co.WithLayoutData(layoutData)
	})
}
