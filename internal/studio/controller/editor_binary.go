package controller

import (
	"fmt"
	"os"

	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/change"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewBinaryEditor(globalCtx global.Context, studio *Studio, editorModel *model.Editor, binModel *model.Binary) *BinaryEditor {
	return &BinaryEditor{
		studio:    studio,
		histModel: editorModel.History(),
		binModel:  binModel,
		editorModel: model.NewBinaryEditor(
			editorModel,
		),
	}
}

var _ Editor = (*BinaryEditor)(nil)

type BinaryEditor struct {
	studio      *Studio
	histModel   *model.History
	binModel    *model.Binary
	editorModel *model.BinaryEditor
}

func (e *BinaryEditor) Save() error {
	if err := e.binModel.Save(); err != nil {
		return fmt.Errorf("error saving model %w", err)
	}
	e.histModel.Save()
	return nil
}

func (e *BinaryEditor) Render(scope co.Scope, layoutData mat.LayoutData) co.Instance {
	return co.New(view.BinaryEditor, func() {
		co.WithData(view.BinaryEditorData{
			ResourceModel:    e.binModel.Resource(),
			BinaryModel:      e.binModel,
			EditorModel:      e.editorModel,
			StudioController: e.studio,
			EditorController: e,
		})
		co.WithLayoutData(layoutData)
	})
}

func (e *BinaryEditor) Destroy() {}

func (e *BinaryEditor) OnRenameResource(name string) {
	e.histModel.Add(change.Name(e.binModel.Resource(),
		change.NameState{
			Value: e.binModel.Resource().Name(),
		},
		change.NameState{
			Value: name,
		},
	))
}

func (e *BinaryEditor) OnChangeContentFromPath(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	ch := change.BinaryContent(e.binModel,
		change.BinaryContentState{
			Data: e.binModel.Data(),
		},
		change.BinaryContentState{
			Data: content,
		},
	)
	if err := e.histModel.Add(ch); err != nil {
		panic(fmt.Errorf("failed to apply change: %w", err))
	}
}
