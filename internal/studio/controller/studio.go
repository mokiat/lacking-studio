package controller

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

func NewStudio(
	window *ui.Window,
	gfxEngine graphics.Engine,
	physicsEngine *physics.Engine,
	ecsEngine *ecs.Engine,
) *Studio {
	result := &Studio{
		Controller: co.NewBaseController(),

		window:    window,
		gfxEngine: gfxEngine,

		actionsVisible:    true,
		propertiesVisible: true,
	}
	result.editors = []Editor{
		NewCubeTextureEditor(result, gfxEngine),
		NewModelEditor(),
	}
	return result
}

type Studio struct {
	co.Controller

	window    *ui.Window
	gfxEngine graphics.Engine

	actionsVisible    bool
	propertiesVisible bool
	activeEditor      Editor
	editors           []Editor
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

func (s *Studio) OpenEditor(editor Editor) {
	panic("TODO")
}

func (s *Studio) CloseEditor(editor Editor) {
	panic("TODO")
}

func (s *Studio) ActiveEditor() Editor {
	return s.activeEditor
}

func (s *Studio) OnEditorClicked(editor Editor) {
	s.activeEditor = editor
	s.NotifyChanged()
}

func (s *Studio) OnEditorClosed(editor Editor) {
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

func (s *Studio) Editors() []Editor {
	return s.editors
}

func (s *Studio) EachEditor(cb func(editor Editor)) {
	for _, editor := range s.editors {
		cb(editor)
	}
}

func (s *Studio) Render() co.Instance {
	return co.New(StudioView, func() {
		co.WithContext(global.Context{
			GFXEngine: s.gfxEngine,
		})
		co.WithData(s)
	})
}

func (s *Studio) editorIndex(editor Editor) int {
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

		if controller.IsActionsVisible() {
			co.WithChild("left", co.New(widget.Paper, func() {
				co.WithData(widget.PaperData{
					Padding: ui.Spacing{
						Top:    20,
						Bottom: 20,
						Left:   1,
						Right:  1,
					},
					Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
						ContentSpacing: 10,
					}),
				})
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentLeft,
				})

				co.WithChild("model", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/model.png"),
						Vertical: true,
					})
				}))

				co.WithChild("light", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/light.png"),
						Vertical: true,
					})
				}))

				co.WithChild("camera", co.New(widget.ToolbarButton, func() {
					co.WithData(widget.ToolbarButtonData{
						Icon:     co.OpenImage("resources/icons/camera.png"),
						Vertical: true,
					})
				}))
			}))
		}

		if controller.IsPropertiesVisible() {
			co.WithChild("right", co.New(mat.Container, func() {
				co.WithData(mat.ContainerData{
					BackgroundColor: optional.NewColor(ui.White()),
					Layout:          mat.NewFillLayout(),
				})
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     optional.NewInt(500),
				})

				if editor := controller.ActiveEditor(); editor != nil {
					key := fmt.Sprintf("content-%s", editor.ID())
					co.WithChild(key, editor.RenderProperties())
				}
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

		co.WithChild("center", co.New(widget.Viewport, func() {
			if editor := controller.ActiveEditor(); editor != nil {
				co.WithData(widget.ViewportData{
					Scene:  editor.Scene(),
					Camera: editor.Camera(),
				})
				co.WithCallbackData(widget.ViewportCallbackData{
					OnUpdate:     editor.Update,
					OnMouseEvent: editor.OnViewportMouseEvent,
				})
			}
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})
		}))
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
		assetsOverlay.Set(co.OpenOverlay(co.New(widget.AssetDialog, func() {
			co.WithCallbackData(widget.AssetDialogCallbackData{
				OnAssetSelected: func() {
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
				Disabled: true,
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

		controller.EachEditor(func(editor Editor) {
			co.WithChild(editor.ID(), co.New(widget.TabbarTab, func() {
				co.WithData(widget.TabbarTabData{
					Icon:     editor.Icon(),
					Text:     editor.Name(),
					Selected: editor == controller.ActiveEditor(),
				})
				co.WithCallbackData(widget.TabbarTabCallbackData{
					OnClick: func() {
						controller.OnEditorClicked(editor)
					},
					OnClose: func() {
						controller.OnEditorClosed(editor)
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
