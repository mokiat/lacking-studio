package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/global"
	"github.com/mokiat/lacking-studio/internal/preview/model"
	"github.com/mokiat/lacking-studio/internal/viewport"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/hierarchy"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Viewport = mvc.EventListener(co.Define(&viewportComponent{}))

type ViewportData struct {
	AppModel *model.AppModel
	Resource *asset.Resource
}

type viewportComponent struct {
	co.BaseComponent

	appModel *model.AppModel
	resource *asset.Resource

	renderAPI render.API

	gameEngine *game.Engine
	gameScene  *game.Scene

	commonData  *viewport.CommonData
	cameraGizmo *viewport.CameraGizmo

	currentResourceSet *game.ResourceSet
	newResourceSet     *game.ResourceSet

	gfxCamera           *graphics.Camera
	gfxGrid             *graphics.Mesh
	gfxAmbientLight     *graphics.AmbientLight
	gfxDirectionalLight *graphics.DirectionalLight
	gfxSky              *graphics.Sky

	modelNode     *hierarchy.Node
	modelPlayback *game.Playback
}

func (c *viewportComponent) OnCreate() {
	data := co.GetData[ViewportData](c.Properties())
	c.appModel = data.AppModel
	c.resource = data.Resource

	window := co.Window(c.Scope())
	c.renderAPI = window.RenderAPI()

	ctx := co.TypedValue[*global.Context](c.Scope())
	c.commonData = ctx.CommonData
	c.gameEngine = ctx.GameEngine

	c.gameScene = c.gameEngine.CreateScene()
	gfxScene := c.gameScene.Graphics()

	c.gfxCamera = gfxScene.CreateCamera()
	c.gfxCamera.SetExposure(1.0)
	c.gfxCamera.SetAutoExposure(false)
	c.gfxCamera.SetFoV(sprec.Degrees(60))
	c.gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	c.refreshAutoExposure()

	c.gfxGrid = gfxScene.CreateMesh(graphics.MeshInfo{
		Definition: c.commonData.GridMeshDefinition(),
	})
	c.gfxGrid.SetMatrix(dprec.IdentityMat4())
	c.refreshShowGrid()

	c.gfxAmbientLight = gfxScene.CreateAmbientLight(graphics.AmbientLightInfo{
		Position:          dprec.ZeroVec3(),
		InnerRadius:       20000.0,
		OuterRadius:       20000.0,
		ReflectionTexture: c.commonData.SkyTexture(),
		RefractionTexture: c.commonData.SkyTexture(),
		CastShadow:        false,
	})
	c.refreshShowAmbientLight()

	c.gfxDirectionalLight = gfxScene.CreateDirectionalLight(graphics.DirectionalLightInfo{
		Position:   c.commonData.SkyColor().VecXYZ(),
		Rotation:   dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
		EmitColor:  dprec.NewVec3(1.0, 1.0, 1.0),
		EmitRange:  20000.0,
		CastShadow: true,
	})
	c.refreshShowDirectionalLight()

	c.gfxSky = gfxScene.CreateSky(graphics.SkyInfo{
		Definition: c.commonData.SkyDefinition(),
	})
	c.refreshShowSky()

	c.cameraGizmo = viewport.NewCameraGizmo(c.gfxCamera)

	c.loadResource()

}

func (c *viewportComponent) OnDelete() {
	c.gameScene.Delete()
	if c.currentResourceSet != nil {
		c.currentResourceSet.Delete()
	}
}

func (c *viewportComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout: layout.Frame(),
		})

		co.WithChild("canvas", co.New(std.Viewport, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ViewportData{
				API: c.renderAPI,
			})
			co.WithCallbackData(std.ViewportCallbackData{
				OnKeyboardEvent: c.handleViewportKeyboardEvent,
				OnMouseEvent:    c.handleViewportMouseEvent,
				OnRender:        c.handleViewportRender,
			})
		}))

		co.WithChild("sidebar", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
				Width:               opt.V(300),
			})
			co.WithData(std.ContainerData{
				Padding:     ui.UniformSpacing(5),
				BorderColor: opt.V(std.OutlineColor),
				BorderSize: ui.Spacing{
					Left: 1,
				},
				Layout: layout.Vertical(layout.VerticalSettings{
					ContentAlignment: layout.HorizontalAlignmentLeft,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("camera-settings", co.New(std.Accordion, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})
				co.WithData(std.AccordionData{
					Title:    "Camera",
					Expanded: c.appModel.CameraSectionExpanded(),
				})
				co.WithCallbackData(std.AccordionCallbackData{
					OnToggle: c.handleCameraSectionExpandedToggle,
				})

				co.WithChild("panel", co.New(std.Container, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ContainerData{
						BorderColor: opt.V(std.OutlineColor),
						BorderSize: ui.Spacing{
							Left:   1,
							Right:  1,
							Bottom: 1,
						},
						Padding: ui.UniformSpacing(2),
						Layout: layout.Vertical(layout.VerticalSettings{
							ContentAlignment: layout.HorizontalAlignmentLeft,
							ContentSpacing:   10,
						}),
					})

					co.WithChild("auto-exposure", co.New(std.Checkbox, func() {
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithData(std.CheckboxData{
							Label:   "Auto Exposure",
							Checked: c.appModel.AutoExposure(),
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: c.handleAutoExposureToggle,
						})
					}))
				}))
			}))

			co.WithChild("scene-settings", co.New(std.Accordion, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})
				co.WithData(std.AccordionData{
					Title:    "Scene",
					Expanded: c.appModel.SceneSectionExpanded(),
				})
				co.WithCallbackData(std.AccordionCallbackData{
					OnToggle: c.handleSceneSectionExpandedToggle,
				})

				co.WithChild("panel", co.New(std.Container, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ContainerData{
						BorderColor: opt.V(std.OutlineColor),
						BorderSize: ui.Spacing{
							Left:   1,
							Right:  1,
							Bottom: 1,
						},
						Padding: ui.UniformSpacing(2),
						Layout: layout.Vertical(layout.VerticalSettings{
							ContentAlignment: layout.HorizontalAlignmentLeft,
							ContentSpacing:   10,
						}),
					})

					co.WithChild("show-grid", co.New(std.Checkbox, func() {
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithData(std.CheckboxData{
							Label:   "Grid",
							Checked: c.appModel.ShowGrid(),
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: c.handleShowGridToggle,
						})
					}))

					co.WithChild("show-ambient-light", co.New(std.Checkbox, func() {
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithData(std.CheckboxData{
							Label:   "Default Ambient Light",
							Checked: c.appModel.ShowAmbientLight(),
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: c.handleShowAmbientLightToggle,
						})
					}))

					co.WithChild("show-directional-light", co.New(std.Checkbox, func() {
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithData(std.CheckboxData{
							Label:   "Default Directional Light",
							Checked: c.appModel.ShowDirectionalLight(),
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: c.handleShowDirectionalLightToggle,
						})
					}))

					co.WithChild("show-sky", co.New(std.Checkbox, func() {
						co.WithLayoutData(layout.Data{
							GrowHorizontally: true,
						})
						co.WithData(std.CheckboxData{
							Label:   "Default Sky",
							Checked: c.appModel.ShowSky(),
						})
						co.WithCallbackData(std.CheckboxCallbackData{
							OnToggle: c.handleShowSkyToggle,
						})
					}))
				}))
			}))
		}))
	})
}

func (c *viewportComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.RefreshEvent:
		c.loadResource()
		c.Invalidate()
	case model.CameraSectionExpandedChangedEvent:
		c.Invalidate()
	case model.AutoExposureChangedEvent:
		c.refreshAutoExposure()
		c.Invalidate()
	case model.SceneSectionExpandedChangedEvent:
		c.Invalidate()
	case model.ShowGridChangedEvent:
		c.refreshShowGrid()
		c.Invalidate()
	case model.ShowAmbientLightChangedEvent:
		c.refreshShowAmbientLight()
		c.Invalidate()
	case model.ShowDirectionalLightChangedEvent:
		c.refreshShowDirectionalLight()
		c.Invalidate()
	case model.ShowSkyChangedEvent:
		c.refreshShowSky()
		c.Invalidate()
	}
}

func (c *viewportComponent) loadResource() {
	c.newResourceSet = c.gameEngine.CreateResourceSet()
	promise := c.newResourceSet.OpenModelByID(c.resource.ID())
	promise.OnSuccess(func(modelDefinition *game.ModelDefinition) {
		co.Schedule(c.Scope(), func() {
			c.handleModelLoaded(modelDefinition)
		})
	})
	promise.OnError(func(err error) {
		co.Schedule(c.Scope(), func() {
			c.handleModelLoadError(err)
		})
	})
}

func (c *viewportComponent) handleViewportKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	return c.cameraGizmo.OnKeyboardEvent(element, event)
}

func (c *viewportComponent) handleViewportMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return c.cameraGizmo.OnMouseEvent(element, event)
}

func (c *viewportComponent) handleViewportRender(framebuffer render.Framebuffer, size ui.Size) {
	c.gameEngine.Update()
	c.gameEngine.Render(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  uint32(size.Width),
		Height: uint32(size.Height),
	})
}

func (c *viewportComponent) handleModelLoaded(modelDefinition *game.ModelDefinition) {
	if c.modelPlayback != nil {
		c.modelPlayback.Stop()
		c.modelPlayback = nil
	}
	if c.modelNode != nil {
		c.modelNode.Delete()
		c.modelNode = nil
	}
	if c.currentResourceSet != nil {
		// FIXME: This panics! WHY?!?!?!
		// resourceSet := c.currentResourceSet
		// co.Schedule(c.Scope(), func() {
		// 	resourceSet.Delete()
		// })
	}
	c.currentResourceSet = c.newResourceSet

	model := c.gameScene.CreateModel(game.ModelInfo{
		Name:       "Model",
		Definition: modelDefinition,
		Position:   dprec.ZeroVec3(),
		Rotation:   dprec.IdentityQuat(),
		Scale:      dprec.NewVec3(1.0, 1.0, 1.0),
		IsDynamic:  false, // FIXME: Setting this to true kills large scenes
	})
	c.modelNode = model.Root()
	if len(model.Animations()) > 0 {
		animation := model.Animations()[0]
		c.modelPlayback = c.gameScene.PlayAnimation(animation)
	}
	// TODO: Find camera and light nodes and attach indicator gizmos to them
	// from the common data.
}

func (c *viewportComponent) handleModelLoadError(err error) {
}

func (c *viewportComponent) handleCameraSectionExpandedToggle(expanded bool) {
	c.appModel.SetCameraSectionExpanded(expanded)
}

func (c *viewportComponent) handleAutoExposureToggle(checked bool) {
	c.appModel.SetAutoExposure(checked)
}

func (c *viewportComponent) refreshAutoExposure() {
	if c.appModel.AutoExposure() {
		c.gfxCamera.SetAutoExposure(true)
	} else {
		c.gfxCamera.SetExposure(1.0)
		c.gfxCamera.SetAutoExposure(false)
	}
}

func (c *viewportComponent) handleSceneSectionExpandedToggle(expanded bool) {
	c.appModel.SetSceneSectionExpanded(expanded)
}

func (c *viewportComponent) handleShowGridToggle(checked bool) {
	c.appModel.SetShowGrid(checked)
}

func (c *viewportComponent) refreshShowGrid() {
	c.gfxGrid.SetActive(c.appModel.ShowGrid())
}

func (c *viewportComponent) handleShowAmbientLightToggle(checked bool) {
	c.appModel.SetShowAmbientLight(checked)
}

func (c *viewportComponent) refreshShowAmbientLight() {
	c.gfxAmbientLight.SetActive(c.appModel.ShowAmbientLight())
}

func (c *viewportComponent) handleShowDirectionalLightToggle(checked bool) {
	c.appModel.SetShowDirectionalLight(checked)
}

func (c *viewportComponent) refreshShowDirectionalLight() {
	c.gfxDirectionalLight.SetActive(c.appModel.ShowDirectionalLight())
}

func (c *viewportComponent) handleShowSkyToggle(checked bool) {
	c.appModel.SetShowSky(checked)
}

func (c *viewportComponent) refreshShowSky() {
	c.gfxSky.SetActive(c.appModel.ShowSky())
}
