package controller

import (
	studiodata "github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewStudio(globalCtx global.Context, studioModel *model.Studio) *Studio {
	return &Studio{
		api:       globalCtx.API,
		registry:  globalCtx.Registry,
		gfxEngine: globalCtx.GraphicsEngine,

		studioModel:       studioModel,
		editorControllers: make(map[*model.Editor]model.IEditor),
	}
}

type Studio struct {
	api       render.API
	registry  *studiodata.Registry
	gfxEngine *graphics.Engine

	studioModel       *model.Studio
	editorControllers map[*model.Editor]model.IEditor
}

func (s *Studio) Model() *model.Studio {
	return s.studioModel
}

func (s *Studio) Reduce(act mvc.Action) bool {
	switch act := act.(type) {
	case action.OpenResource:
		s.openResource(act.ID)
		return true
	case action.CloneResource:
		s.cloneResource(act.Resource.ID())
		return true
	case action.DeleteResource:
		s.deleteResource(act.Resource.ID())
		return true
	case action.ChangeSelectedEditor:
		s.changeSelectedEditor(act.Editor)
		return true
	case action.CloseEditor:
		s.closeEditor(act.Editor)
		return true
	case action.Undo:
		s.undo()
		return true
	case action.Redo:
		s.redo()
		return true
	case action.Save:
		s.save()
		return true
	default:
		return false
	}
}

func (s *Studio) openResource(id string) {
	if editor := s.studioModel.FindEditorByID(id); editor != nil {
		s.studioModel.SetSelectedEditor(editor)
		return
	}

	resource := s.registry.GetResourceByID(id)
	resourceModel := model.NewResource(resource)
	editorModel := model.NewEditor(resourceModel)

	switch resource.Kind() {
	case studiodata.ResourceKindTwoDTexture:
		texModel, err := model.OpenTwoDTexture(resourceModel)
		if err != nil {
			panic("TODO")
		}
		controller := NewTwoDTextureEditor(s, editorModel, texModel)
		s.editorControllers[editorModel] = controller

	case studiodata.ResourceKindCubeTexture:
		texModel, err := model.OpenCubeTexture(resourceModel)
		if err != nil {
			panic("TODO")
		}
		controller := NewCubeTextureEditor(s, editorModel, texModel)
		s.editorControllers[editorModel] = controller

	case studiodata.ResourceKindModel:
		log.Info("TODO: Open Model")
		return

	case studiodata.ResourceKindScene:
		log.Info("TODO: Open Scene")
		return
	}

	s.studioModel.AddEditor(editorModel)
	s.studioModel.SetSelectedEditor(editorModel)
}

func (s *Studio) cloneResource(id string) {
	// resource := s.registry.GetResourceByID(id)
	// newResource, err := resource.Clone()
	// if err != nil {
	// 	s.HandleError(err)
	// 	return
	// }
	// s.OpenAsset(newResource.ID())
}

func (s *Studio) deleteResource(id string) {
	// // TODO: Open confirmation dialog
	// for _, editor := range s.editors {
	// 	if editor.ID() == id {
	// 		s.CloseEditor(editor)
	// 		break
	// 	}
	// }
	// resource := s.registry.GetResourceByID(id)
	// if err := resource.Delete(); err != nil {
	// 	s.HandleError(err)
	// 	return
	// }
}

func (s *Studio) changeSelectedEditor(editor *model.Editor) {
	s.studioModel.SetSelectedEditor(editor)
}

func (s *Studio) closeEditor(editor *model.Editor) {
	s.studioModel.RemoveEditor(editor)
	// editor.Destroy() // TODO
}

func (s *Studio) undo() {
	if err := s.studioModel.SelectedHistory().Undo(); err != nil {
		panic("TODO")
	}
}

func (s *Studio) redo() {
	if err := s.studioModel.SelectedHistory().Redo(); err != nil {
		panic("TODO")
	}
}

func (s *Studio) save() {
	// TODO
}

func (s *Studio) GraphicsEngine() *graphics.Engine {
	return s.gfxEngine
}

func (s *Studio) Save() {
	// TODO: Hm.... we need to call Save on the underlying model...
	// This might not be possible without some type of reverse dispatch
	// or via the handle after all....
	// s.SelectedHistory().Save() // TODO: Handle error
	// Or maybe we can have ownership over the controller (as we do)
	// and can register it as a Reducer in the UI so that the Save
	// button would see it...
	// if err := s.activeEditor.Save(); err != nil {
	// 	panic(err)
	// }
}

func (s *Studio) Render(scope co.Scope) co.Instance {
	return co.New(view.Studio, func() {
		co.WithData(view.StudioData{
			StudioModel:      s.studioModel,
			StudioController: s,
		})
		co.WithScope(mvc.UseReducer(scope, s))
	})
}

func (s *Studio) RenderEditor(editorModel *model.Editor, scope co.Scope, layoutData mat.LayoutData) co.Instance {
	controller := s.editorControllers[editorModel]
	return controller.Render(scope, layoutData)
}
