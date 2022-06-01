package model

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeHistory     = mvc.NewChange("history")
	ChangeHistoryAdd  = mvc.SubChange(ChangeHistory, "add")
	ChangeHistoryUndo = mvc.SubChange(ChangeHistory, "undo")
	ChangeHistoryRedo = mvc.SubChange(ChangeHistory, "redo")
	ChangeHistorySave = mvc.SubChange(ChangeHistory, "save")
)

const maxUndoCount = 10

func NewHistory() *History {
	return &History{
		Observable:  mvc.NewObservable(),
		changes:     history.NewQueue(maxUndoCount),
		savedChange: nil,
	}
}

type History struct {
	mvc.Observable
	changes     *history.Queue
	savedChange history.Change
}

func (h *History) Add(ch history.Change) error {
	if err := h.changes.Push(ch); err != nil {
		return fmt.Errorf("error pushing change: %w", err)
	}
	h.SignalChange(ChangeHistoryAdd)
	return nil
}

func (h *History) CanUndo() bool {
	return h.changes.CanPop()
}

func (h *History) Undo() error {
	if err := h.changes.Pop(); err != nil {
		return fmt.Errorf("error popping change: %w", err)
	}
	h.SignalChange(ChangeHistoryUndo)
	return nil
}

func (h *History) CanRedo() bool {
	return h.changes.CanUnpop()
}

func (h *History) Redo() error {
	if err := h.changes.Unpop(); err != nil {
		return fmt.Errorf("error unpopping change: %w", err)
	}
	h.SignalChange(ChangeHistoryRedo)
	return nil
}

func (h *History) CanSave() bool {
	return h.savedChange != h.changes.LastChange()
}

func (h *History) Save() {
	h.savedChange = h.changes.LastChange()
	h.SignalChange(ChangeHistorySave)
}
