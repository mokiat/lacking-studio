package controller

import (
	"fmt"
	"log"

	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

func NewStudio(
	projectDir string,
	window *ui.Window,
	registry asset.Registry,
	gfxEngine graphics.Engine,
	physicsEngine *physics.Engine,
	ecsEngine *ecs.Engine,
) *Studio {
	result := &Studio{
		Controller: co.NewBaseController(),

		projectDir: projectDir,
		window:     window,
		registry:   data.NewRegistry(registry),
		gfxEngine:  gfxEngine,

		actionsVisible:    true,
		propertiesVisible: true,
	}
	result.editors = []model.Editor{}
	return result
}

type Studio struct {
	co.Controller

	projectDir string
	window     *ui.Window
	registry   *data.Registry
	gfxEngine  graphics.Engine

	actionsVisible    bool
	propertiesVisible bool
	activeEditor      model.Editor
	editors           []model.Editor
}

func (s *Studio) ProjectDir() string {
	return s.projectDir
}

func (s *Studio) Window() *ui.Window {
	return s.window
}

func (s *Studio) Registry() *data.Registry {
	return s.registry
}

func (s *Studio) GraphicsEngine() graphics.Engine {
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
	resource := s.registry.GetResourceByID(id)
	for _, editor := range s.editors {
		if editor.ID() == resource.GUID {
			s.SelectEditor(editor)
			return
		}
	}

	switch resource.Kind {
	case "twod_texture":
		editor, err := NewTwoDTextureEditor(s, &resource.Resource)
		if err != nil {
			panic(err) // TODO
		}
		s.OpenEditor(editor)
	case "cube_texture":
		editor, err := NewCubeTextureEditor(s, &resource.Resource)
		if err != nil {
			panic(err) // TODO
		}
		s.OpenEditor(editor)
	case "model":
		log.Println("TODO")
	case "scene":
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

// type ApplicationController struct {
// 	co.Controller
// 	propertiesVisible bool

// 	viewportController *ViewportController
// }

// func (c *ApplicationController) Init() {
// 	c.viewportController.Init()
// }

// func (c *ApplicationController) Free() {
// 	c.viewportController.Free()
// }

// func (c *ApplicationController) IsPropertiesVisible() bool {
// 	return c.propertiesVisible
// }

// func (c *ApplicationController) TogglePropertiesVisible() {
// 	c.propertiesVisible = !c.propertiesVisible
// 	c.NotifyChanged()
// }

var StudioView = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	controller := props.Data().(*Studio)

	co.OpenFontCollection("resources/fonts/roboto.ttc")

	// co.Once(func() {
	// 	controller.Init()
	// })

	// co.Defer(func() {
	// 	controller.Free()
	// })

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(widget.BackgroundColor),
			Layout:          mat.NewFrameLayout(),
		})

		co.WithChild("top", co.New(StudioTopPanel, func() {
			co.WithData(props.Data())
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentTop,
			})
		}))

		// if controller.IsActionsVisible() {
		// 	co.WithChild("left", co.New(widget.Paper, func() {
		// 		co.WithData(widget.PaperData{
		// 			Padding: ui.Spacing{
		// 				Top:    20,
		// 				Bottom: 20,
		// 				Left:   1,
		// 				Right:  1,
		// 			},
		// 			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
		// 				ContentSpacing: 10,
		// 			}),
		// 		})
		// 		co.WithLayoutData(mat.LayoutData{
		// 			Alignment: mat.AlignmentLeft,
		// 		})

		// 		co.WithChild("model", co.New(widget.ToolbarButton, func() {
		// 			co.WithData(widget.ToolbarButtonData{
		// 				Icon:     co.OpenImage("resources/icons/model.png"),
		// 				Vertical: true,
		// 			})
		// 		}))

		// 		co.WithChild("light", co.New(widget.ToolbarButton, func() {
		// 			co.WithData(widget.ToolbarButtonData{
		// 				Icon:     co.OpenImage("resources/icons/light.png"),
		// 				Vertical: true,
		// 			})
		// 		}))

		// 		co.WithChild("camera", co.New(widget.ToolbarButton, func() {
		// 			co.WithData(widget.ToolbarButtonData{
		// 				Icon:     co.OpenImage("resources/icons/camera.png"),
		// 				Vertical: true,
		// 			})
		// 		}))
		// 	}))
		// }

		// if controller.IsPropertiesVisible() {
		// 	co.WithChild("right", co.New(mat.Container, func() {
		// 		co.WithData(mat.ContainerData{
		// 			BackgroundColor: optional.NewColor(ui.White()),
		// 			Layout:          mat.NewFillLayout(),
		// 		})
		// 		co.WithLayoutData(mat.LayoutData{
		// 			Alignment: mat.AlignmentRight,
		// 			Width:     optional.NewInt(500),
		// 		})

		// 		if editor := controller.ActiveEditor(); editor != nil {
		// 			key := fmt.Sprintf("content-%s", editor.ID())
		// 			co.WithChild(key, editor.RenderProperties())
		// 		}
		// 	}))
		// }

		if editor := controller.ActiveEditor(); editor != nil {
			key := fmt.Sprintf("center-%s", editor.ID())
			co.WithChild(key, editor.Render(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			}))
		}

		// co.WithChild("center", co.New(mat.Container, func() {
		// 	co.WithData(mat.ContainerData{
		// 		BackgroundColor: optional.NewColor(ui.Black()),
		// 	})
		// 	co.WithLayoutData(mat.LayoutData{
		// 		Alignment: mat.AlignmentCenter,
		// 	})
		// }))

		// co.WithChild("center", co.New(widget.Viewport, func() {
		// 	if editor := controller.ActiveEditor(); editor != nil {
		// 		co.WithData(widget.ViewportData{
		// 			Scene:  editor.Scene(),
		// 			Camera: editor.Camera(),
		// 		})
		// 		co.WithCallbackData(widget.ViewportCallbackData{
		// 			OnUpdate:     editor.Update,
		// 			OnMouseEvent: editor.OnViewportMouseEvent,
		// 		})
		// 	}
		// 	co.WithLayoutData(mat.LayoutData{
		// 		Alignment: mat.AlignmentCenter,
		// 	})
		// }))
	})
}))

var StudioTopPanel = co.Controlled(co.Define(func(props co.Properties) co.Instance {
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

var Toolbar = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	controller := props.Data().(*Studio)

	assetsOverlay := co.UseState(func() interface{} {
		return co.Overlay{}
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
					overlay := assetsOverlay.Get().(co.Overlay)
					overlay.Close()
				},
			})
		})))
	}

	onPropertiesVisibleClicked := func() {
		controller.SetPropertiesVisible(!controller.IsPropertiesVisible())
	}

	return co.New(widget.Toolbar, func() {
		co.WithLayoutData(props.LayoutData())

		co.WithChild("assets", co.New(widget.ToolbarButton, func() {
			co.WithData(widget.ToolbarButtonData{
				Icon: co.OpenImage("resources/icons/assets.png"),
				Text: "Assets",
			})
			co.WithCallbackData(widget.ToolbarButtonCallbackData{
				ClickListener: onAssetsClicked,
			})
		}))

		co.WithChild("separator1", co.New(widget.ToolbarSeparator, nil))

		co.WithChild("save", co.New(widget.ToolbarButton, func() {
			co.WithData(widget.ToolbarButtonData{
				Icon:     co.OpenImage("resources/icons/save.png"),
				Disabled: !controller.SaveEnabled(),
			})
			co.WithCallbackData(widget.ToolbarButtonCallbackData{
				ClickListener: func() {
					controller.Save()
				},
			})
		}))

		co.WithChild("separator2", co.New(widget.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(widget.ToolbarButton, func() {
			co.WithData(widget.ToolbarButtonData{
				Icon:     co.OpenImage("resources/icons/undo.png"),
				Disabled: !controller.UndoEnabled(),
			})
			co.WithCallbackData(widget.ToolbarButtonCallbackData{
				ClickListener: func() {
					controller.Undo()
				},
			})
		}))

		co.WithChild("redo", co.New(widget.ToolbarButton, func() {
			co.WithData(widget.ToolbarButtonData{
				Icon:     co.OpenImage("resources/icons/redo.png"),
				Disabled: !controller.RedoEnabled(),
			})
			co.WithCallbackData(widget.ToolbarButtonCallbackData{
				ClickListener: func() {
					controller.Redo()
				},
			})
		}))

		co.WithChild("separator3", co.New(widget.ToolbarSeparator, nil))

		co.WithChild("properties", co.New(widget.ToolbarButton, func() {
			co.WithData(widget.ToolbarButtonData{
				Icon: co.OpenImage("resources/icons/properties.png"),
			})
			co.WithCallbackData(widget.ToolbarButtonCallbackData{
				ClickListener: onPropertiesVisibleClicked,
			})
		}))
	})
}))

var Tabbar = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	controller := props.Data().(*Studio)

	return co.New(widget.Tabbar, func() {
		co.WithLayoutData(props.LayoutData())

		controller.EachEditor(func(editor model.Editor) {
			co.WithChild(editor.ID(), co.New(widget.TabbarTab, func() {
				co.WithData(widget.TabbarTabData{
					Icon:     editor.Icon(),
					Text:     editor.Name(),
					Selected: editor == controller.ActiveEditor(),
				})
				co.WithCallbackData(widget.TabbarTabCallbackData{
					OnClick: func() {
						controller.SelectEditor(editor)
					},
					OnClose: func() {
						controller.CloseEditor(editor)
					},
				})
			}))
		})

		// co.WithChild("forest", co.New(widget.TabbarTab, func() {
		// 	co.WithData(widget.TabbarTabData{
		// 		Icon: co.OpenImage("resources/icons/scene.png"),
		// 		Text: "Мега Сцена",
		// 	})
		// }))

		// co.WithChild("tree", co.New(widget.TabbarTab, func() {
		// 	co.WithData(widget.TabbarTabData{
		// 		Icon: co.OpenImage("resources/icons/model.png"),
		// 		Text: "Дърво",
		// 	})
		// }))

		// co.WithChild("car", co.New(widget.TabbarTab, func() {
		// 	co.WithData(widget.TabbarTabData{
		// 		Icon: co.OpenImage("resources/icons/model.png"),
		// 		Text: "Кола",
		// 	})
		// }))

		// co.WithChild("stone", co.New(widget.TabbarTab, func() {
		// 	co.WithData(widget.TabbarTabData{
		// 		Icon: co.OpenImage("resources/icons/model.png"),
		// 		Text: "Камък",
		// 	})
		// }))
	})
}))
