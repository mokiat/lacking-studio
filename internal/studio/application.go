package studio

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

func BootstrapApplication(window *ui.Window, gfxEngine graphics.Engine) {
	co.Initialize(window, co.New(co.StoreProvider, func() {
		co.WithData(co.StoreProviderData{
			Entries: []co.StoreProviderEntry{
				co.NewStoreProviderEntry(global.Reducer()),
			},
		})

		co.WithChild("app", co.New(Application, func() {
			co.WithContext(global.Context{
				GFXEngine: gfxEngine,
			})
			co.WithData(&ApplicationController{
				Controller:        co.NewBaseController(),
				propertiesVisible: true,
			})
		}))
	}))
}

// type ApplicationData = mat.SwitchData

type ApplicationController struct {
	co.Controller
	propertiesVisible bool
}

func (c *ApplicationController) IsPropertiesVisible() bool {
	return c.propertiesVisible
}

func (c *ApplicationController) TogglePropertiesVisible() {
	c.propertiesVisible = !c.propertiesVisible
	c.NotifyChanged()
}

var Application = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	co.OpenFontCollection("resources/fonts/roboto.ttc")

	controller := props.Data().(*ApplicationController)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(widget.BackgroundColor),
			Layout:          mat.NewFrameLayout(),
		})

		co.WithChild("top", co.New(TopPanel, func() {
			co.WithData(props.Data())
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentTop,
			})
		}))

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

		if controller.IsPropertiesVisible() {
			co.WithChild("right", co.New(mat.Container, func() {
				co.WithData(mat.ContainerData{
					BackgroundColor: optional.NewColor(ui.White()),
				})
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     optional.NewInt(300),
				})
			}))
		}

		co.WithChild("center", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				BackgroundColor: optional.NewColor(ui.Black()),
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})
		}))
	})
}))

var TopPanel = co.Controlled(co.Define(func(props co.Properties) co.Instance {
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
	controller := props.Data().(*ApplicationController)

	assetsOverlay := co.UseState(func() interface{} {
		return co.Overlay{}
	})

	onAssetsClicked := func() {
		assetsOverlay.Set(co.OpenOverlay(co.New(widget.AssetDialog, func() {
			co.WithCallbackData(widget.AssetDialogCallbackData{
				OnAssetSelected: func() {
					fmt.Println("ASSET SELECTED")
				},
				OnClose: func() {
					overlay := assetsOverlay.Get().(co.Overlay)
					overlay.Close()
				},
			})
		})))
	}

	onPropertiesVisibleClicked := func() {
		controller.TogglePropertiesVisible()
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
				Disabled: true,
			})
		}))

		co.WithChild("redo", co.New(widget.ToolbarButton, func() {
			co.WithData(widget.ToolbarButtonData{
				Icon:     co.OpenImage("resources/icons/redo.png"),
				Disabled: true,
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
	return co.New(widget.Tabbar, func() {
		co.WithLayoutData(props.LayoutData())

		co.WithChild("car", co.New(widget.TabbarTab, func() {
			co.WithData(widget.TabbarTabData{
				Text: "car",
			})
		}))

		co.WithChild("tree", co.New(widget.TabbarTab, func() {
			co.WithData(widget.TabbarTabData{
				Text:     "tree",
				Selected: true,
			})
		}))
	})
}))
