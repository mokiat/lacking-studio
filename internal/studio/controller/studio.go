package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/change"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewStudio(globalCtx global.Context, studioModel *model.Studio) *Studio {
	return &Studio{
		globalCtx:         globalCtx,
		studioModel:       studioModel,
		resourceModels:    make(map[asset.Resource]*model.Resource),
		editorControllers: make(map[*model.Editor]Editor),
	}
}

type Studio struct {
	globalCtx   global.Context
	studioModel *model.Studio

	resourceModels    map[asset.Resource]*model.Resource
	editorControllers map[*model.Editor]Editor
}

func (s *Studio) OnSave() {
	editor := s.studioModel.SelectedEditor()
	if editor == nil {
		log.Warn("[studio] Trying to save unselected editor")
		return
	}
	controller := s.editorControllers[editor]
	if err := controller.Save(); err != nil {
		panic(err)
	}
}

func (s *Studio) OnUndo() {
	if err := s.studioModel.SelectedHistory().Undo(); err != nil {
		panic("TODO")
	}
}

func (s *Studio) OnRedo() {
	if err := s.studioModel.SelectedHistory().Redo(); err != nil {
		panic("TODO")
	}
}

func (s *Studio) OnToggleProperties() {
	editor := s.studioModel.SelectedEditor()
	if editor == nil {
		return
	}
	editor.SetPropertiesVisible(!editor.IsPropertiesVisible())
}

func (s *Studio) OnCreateResource(kind model.ResourceKind) {
	log.Warn("TODO: Create Resource")
}

func (s *Studio) OnOpenResource(id string) {
	if editor := s.studioModel.FindEditorByID(id); editor != nil {
		s.studioModel.SetSelectedEditor(editor)
		return
	}

	resourceModel := s.registry().ResourceByID(id)
	editorModel := model.NewEditor(resourceModel)

	switch resourceModel.Kind() {
	case model.ResourceKindTwoDTexture:
		texModel, err := model.OpenTwoDTexture(resourceModel)
		if err != nil {
			panic("TODO")
		}
		controller := NewTwoDTextureEditor(s.globalCtx, s, editorModel, texModel)
		s.editorControllers[editorModel] = controller

	case model.ResourceKindCubeTexture:
		texModel, err := model.OpenCubeTexture(resourceModel)
		if err != nil {
			panic("TODO")
		}
		controller := NewCubeTextureEditor(s.globalCtx, s, editorModel, texModel)
		s.editorControllers[editorModel] = controller

	case model.ResourceKindModel:
		log.Info("TODO: Open Model")
		return

	case model.ResourceKindScene:
		log.Info("TODO: Open Scene")
		return
	}

	s.studioModel.AddEditor(editorModel)
	s.studioModel.SetSelectedEditor(editorModel)
}

func (s *Studio) OnRenameResource(id, name string) {
	resource := s.registry().ResourceByID(id)
	if resource == nil {
		log.Warn("[studio] Trying to rename missing resource")
		return
	}
	s.studioModel.SelectedHistory().Add(change.Name(resource,
		change.NameState{
			Value: resource.Name(),
		},
		change.NameState{
			Value: name,
		},
	))
}

func (s *Studio) OnCloneResource(id string) *model.Resource {
	resource := s.registry().ResourceByID(id)
	if resource == nil {
		log.Warn("[studio] Trying to clone missing resource")
		return nil
	}
	newResource, err := resource.Clone()
	if err != nil {
		panic(err)
	}
	return newResource
}

func (s *Studio) OnDeleteResource(id string) {
	editor := s.studioModel.FindEditorByID(id)
	if editor != nil {
		s.OnCloseEditor(editor)
	}
	resource := s.registry().ResourceByID(id)
	if resource == nil {
		log.Warn("[studio] Trying to delete missing resource")
	}
	if err := resource.Delete(); err != nil {
		panic(err)
	}
}

func (s *Studio) OnSelectEditor(editor *model.Editor) {
	s.studioModel.SetSelectedEditor(editor)
}

func (s *Studio) OnCloseEditor(editor *model.Editor) {
	// TODO: Confirmation dialog if unsaved (can be in the UI?)
	s.studioModel.RemoveEditor(editor)
	if controller, ok := s.editorControllers[editor]; ok {
		controller.Destroy()
		delete(s.editorControllers, editor)
	}
}

func (s *Studio) Render(scope co.Scope) co.Instance {
	return co.New(view.Studio, func() {
		co.WithData(view.StudioData{
			StudioModel:      s.studioModel,
			StudioController: s,
		})
	})
}

func (s *Studio) RenderEditor(editorModel *model.Editor, scope co.Scope, layoutData mat.LayoutData) co.Instance {
	controller := s.editorControllers[editorModel]
	return controller.Render(scope, layoutData)
}

func (s *Studio) registry() *model.Registry {
	return s.studioModel.Registry()
}
