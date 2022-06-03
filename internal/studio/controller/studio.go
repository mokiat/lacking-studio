package controller

import (
	"fmt"

	studiodata "github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/filter"
	"github.com/mokiat/lacking/util/optional"
)

func NewStudio(
	window *ui.Window,
	api render.API,
	registry *studiodata.Registry,
	gfxEngine *graphics.Engine,
	physicsEngine *physics.Engine,
	ecsEngine *ecs.Engine,
) *Studio {
	return &Studio{
		api:       api,
		window:    window,
		registry:  registry,
		gfxEngine: gfxEngine,

		studioModel:       model.NewStudio(),
		editorControllers: make(map[*model.Editor]model.IEditor),

		actionsVisible:    true,
		propertiesVisible: true,
	}
}

type Studio struct {
	api       render.API
	window    *ui.Window
	registry  *studiodata.Registry
	gfxEngine *graphics.Engine

	studioModel       *model.Studio
	editorControllers map[*model.Editor]model.IEditor

	actionsVisible    bool
	propertiesVisible bool
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
		// TODO: Send to the currently selected editor
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

func (s *Studio) HandleError(err error) {
	if err != nil {
		panic(err)
	}
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

func (s *Studio) Render() co.Instance {
	return co.New(StudioView, func() {
		co.WithData(s)
	})
}

// TODO: Move to view package
var StudioView = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	controller := props.Data().(*Studio)
	scope = mvc.UseReducer(scope, controller)

	studioModel := controller.Model()
	mvc.UseBinding(studioModel, filter.Always[mvc.Change]()) // TODO

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithScope(scope)

		co.WithChild("top", co.New(StudioTopPanel, func() {
			co.WithData(props.Data())
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentTop,
			})
		}))

		if editor := studioModel.SelectedEditor(); editor != nil {
			editorController := controller.editorControllers[editor]
			key := fmt.Sprintf("center-%s", editor.Resource().ID())
			co.WithChild(key, editorController.Render(scope, mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			}))
		}
	})
})

var StudioTopPanel = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	controller := props.Data().(*Studio)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
			}),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("toolbar", co.New(view.StudioToolbar, func() {
			co.WithData(view.StudioToolbarData{
				StudioModel: controller.Model(),
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))

		co.WithChild("tabbar", co.New(view.StudioTabbar, func() {
			co.WithData(view.StudioTabbarData{
				StudioModel: controller.Model(),
			})
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))
	})
})
