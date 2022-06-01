package controller

import (
	"fmt"
	"log"

	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/optional"
)

func NewStudio(
	window *ui.Window,
	api render.API,
	registry asset.Registry,
	gfxEngine *graphics.Engine,
	physicsEngine *physics.Engine,
	ecsEngine *ecs.Engine,
) *Studio {
	dataStudio := data.NewRegistry(registry)
	if err := dataStudio.Init(); err != nil {
		panic(err) // TODO
	}

	result := &Studio{
		Controller: co.NewBaseController(),

		target: mvc.NewObservable(),

		api: api,

		window:    window,
		registry:  dataStudio,
		gfxEngine: gfxEngine,

		actionsVisible:    true,
		propertiesVisible: true,
	}
	result.editors = []model.Editor{}
	return result
}

type Studio struct {
	co.Controller

	target mvc.Observable

	api render.API

	window    *ui.Window
	registry  *data.Registry
	gfxEngine *graphics.Engine

	actionsVisible    bool
	propertiesVisible bool
	activeEditor      model.Editor
	editors           []model.Editor
}

func (s *Studio) Reduce(act mvc.Action) bool {
	switch act := act.(type) {
	case action.CloneResource:
		s.cloneResource(act.Resource.ID())
		return true
	case action.DeleteResource:
		s.deleteResource(act.Resource.ID())
		return true
	default:
		return false
	}
}

func (s *Studio) cloneResource(id string) {
	resource := s.registry.GetResourceByID(id)
	newResource, err := resource.Clone()
	if err != nil {
		s.HandleError(err)
		return
	}
	s.OpenAsset(newResource.ID())
}

func (s *Studio) deleteResource(id string) {
	// TODO: Open confirmation dialog
	for _, editor := range s.editors {
		if editor.ID() == id {
			s.CloseEditor(editor)
			break
		}
	}
	resource := s.registry.GetResourceByID(id)
	if err := resource.Delete(); err != nil {
		s.HandleError(err)
		return
	}
}

func (s *Studio) Target() mvc.Observable {
	return s.target
}

func (s *Studio) HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *Studio) Window() *ui.Window {
	return s.window
}

func (s *Studio) Registry() *data.Registry {
	return s.registry
}

func (s *Studio) GraphicsEngine() *graphics.Engine {
	return s.gfxEngine
}

func (s *Studio) IsActionsVisible() bool {
	return s.actionsVisible
}

func (s *Studio) SetActionsVisible(visible bool) {
	s.actionsVisible = visible
	s.NotifyChanged()
}

func (s *Studio) IsPropertiesVisible() bool {
	return s.propertiesVisible
}

func (s *Studio) SetPropertiesVisible(visible bool) {
	s.propertiesVisible = visible
	s.NotifyChanged()
}

func (s *Studio) UndoEnabled() bool {
	if s.activeEditor == nil {
		return false
	}
	return s.activeEditor.CanUndo()
}

func (s *Studio) Undo() {
	s.activeEditor.Undo()
	s.NotifyChanged()
}

func (s *Studio) RedoEnabled() bool {
	if s.activeEditor == nil {
		return false
	}
	return s.activeEditor.CanRedo()
}

func (s *Studio) Redo() {
	s.activeEditor.Redo()
	s.NotifyChanged()
}

func (s *Studio) SaveEnabled() bool {
	if s.activeEditor == nil {
		return false
	}
	return s.activeEditor.CanSave()
}

func (s *Studio) Save() {
	if err := s.activeEditor.Save(); err != nil {
		panic(err)
	}
	s.NotifyChanged()
}

func (s *Studio) OpenAsset(id string) {
	for _, editor := range s.editors {
		if editor.ID() == id {
			s.SelectEditor(editor)
			return
		}
	}

	resource := s.registry.GetResourceByID(id)
	switch resource.Kind() {
	case data.ResourceKindTwoDTexture:
		texModel, err := model.OpenTwoDTexture(s.registry, id)
		if err != nil {
			s.HandleError(err)
			return
		}
		s.OpenEditor(NewTwoDTextureEditor(s, texModel))
	case data.ResourceKindCubeTexture:
		texModel, err := model.OpenCubeTexture(s.registry, id)
		if err != nil {
			s.HandleError(err)
			return
		}
		s.OpenEditor(NewCubeTextureEditor(s, texModel))
	case data.ResourceKindModel:
		log.Println("TODO")
	case data.ResourceKindScene:
		log.Println("TODO")
	}
}

func (s *Studio) OpenEditor(editor model.Editor) {
	s.editors = append(s.editors, editor)
	s.activeEditor = editor
	s.NotifyChanged()
}

func (s *Studio) ActiveEditor() model.Editor {
	return s.activeEditor
}

func (s *Studio) SelectEditor(editor model.Editor) {
	s.activeEditor = editor
	s.NotifyChanged()
}

func (s *Studio) CloseEditor(editor model.Editor) {
	editorIndex := s.editorIndex(editor)
	if editorIndex < 0 {
		return
	}

	s.editors = append(s.editors[:editorIndex], s.editors[editorIndex+1:]...)
	if editor == s.activeEditor {
		switch {
		case len(s.editors) == 0:
			s.activeEditor = nil
		case editorIndex < len(s.editors):
			s.activeEditor = s.editors[editorIndex]
		default:
			s.activeEditor = s.editors[editorIndex-1]
		}
	}

	editor.Destroy()

	s.NotifyChanged()
}

func (s *Studio) Editors() []model.Editor {
	return s.editors
}

func (s *Studio) EachEditor(cb func(editor model.Editor)) {
	for _, editor := range s.editors {
		cb(editor)
	}
}

func (s *Studio) Render() co.Instance {
	return co.New(StudioView, func() {
		co.WithData(s)
	})
}

func (s *Studio) editorIndex(editor model.Editor) int {
	for i, candidate := range s.editors {
		if candidate == editor {
			return i
		}
	}
	return -1
}

// TODO: Move to view package
var StudioView = co.Controlled(co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	controller := props.Data().(*Studio)
	scope = mvc.UseReducer(scope, controller)

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

		if editor := controller.ActiveEditor(); editor != nil {
			key := fmt.Sprintf("center-%s", editor.ID())
			co.WithChild(key, editor.Render(scope, mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			}))
		}
	})
}))

var StudioTopPanel = co.Controlled(co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
			}),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("toolbar", co.New(Toolbar, func() {
			co.WithData(props.Data())
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))

		co.WithChild("tabbar", co.New(Tabbar, func() {
			co.WithData(props.Data())
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
		}))
	})
}))

var Toolbar = co.Controlled(co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	controller := props.Data().(*Studio)

	assetsOverlay := co.UseState(func() co.Overlay {
		return nil
	})

	onAssetsClicked := func() {
		assetsOverlay.Set(co.OpenOverlay(co.New(view.AssetDialog, func() {
			co.WithData(view.AssetDialogData{
				Registry: controller.Registry(),
			})
			co.WithCallbackData(view.AssetDialogCallbackData{
				OnAssetSelected: func(id string) {
					controller.OpenAsset(id)
				},
				OnClose: func() {
					overlay := assetsOverlay.Get()
					overlay.Close()
				},
			})
		})))
	}

	onPropertiesVisibleClicked := func() {
		controller.SetPropertiesVisible(!controller.IsPropertiesVisible())
	}

	return co.New(mat.Toolbar, func() {
		co.WithLayoutData(props.LayoutData())

		co.WithChild("assets", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon: co.OpenImage(scope, "icons/assets.png"),
				Text: "Assets",
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: onAssetsClicked,
			})
		}))

		co.WithChild("separator1", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("save", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/save.png"),
				Enabled: optional.Value(controller.SaveEnabled()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					controller.Save()
				},
			})
		}))

		co.WithChild("separator2", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/undo.png"),
				Enabled: optional.Value(controller.UndoEnabled()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					controller.Undo()
				},
			})
		}))

		co.WithChild("redo", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/redo.png"),
				Enabled: optional.Value(controller.RedoEnabled()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					controller.Redo()
				},
			})
		}))

		co.WithChild("separator3", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("properties", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon: co.OpenImage(scope, "icons/properties.png"),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: onPropertiesVisibleClicked,
			})
		}))
	})
}))

var Tabbar = co.Controlled(co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	controller := props.Data().(*Studio)

	return co.New(mat.Tabbar, func() {
		co.WithLayoutData(props.LayoutData())

		controller.EachEditor(func(editor model.Editor) {
			co.WithChild(editor.ID(), co.New(mat.TabbarTab, func() {
				co.WithData(mat.TabbarTabData{
					Icon:     editor.Icon(scope),
					Text:     editor.Name(),
					Selected: editor == controller.ActiveEditor(),
				})
				co.WithCallbackData(mat.TabbarTabCallbackData{
					OnClick: func() {
						controller.SelectEditor(editor)
					},
					OnClose: func() {
						controller.CloseEditor(editor)
					},
				})
			}))
		})
	})
}))
