package editor

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/mokiat/gomath/dprec"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking-studio/internal/view/common"
	"github.com/mokiat/lacking-studio/internal/visualization"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/lacking/util/async"
)

var Workbench = co.Define(&workbenchComponent{})

type WorkbenchData struct {
	EditorModel *editormodel.Model
}

type workbenchComponent struct {
	co.BaseComponent

	editorModel *editormodel.Model
	vis         *visualization.Fragment

	renderAPI render.API
	// commonData *viewport.CommonData

	// gfxEngine *graphics.Engine
	// gfxScene  *graphics.Scene
	// gfxCamera *graphics.Camera

	// cameraGizmo *viewport.CameraGizmo
}

func (c *workbenchComponent) OnCreate() {
	data := co.GetData[WorkbenchData](c.Properties())
	c.editorModel = data.EditorModel
	c.vis = c.editorModel.Visualization()

	window := co.Window(c.Scope())
	c.renderAPI = window.RenderAPI()

	// c.commonData = co.TypedValue[*viewport.CommonData](c.Scope())
	// c.gfxEngine = co.TypedValue[*graphics.Engine](c.Scope())

	// c.gfxScene = c.gfxEngine.CreateScene()
	// c.gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.01, 0.01, 0.02))

	// c.gfxCamera = c.gfxScene.CreateCamera()
	// c.gfxCamera.SetExposure(1.0)
	// c.gfxCamera.SetAutoExposure(false)
	// c.gfxCamera.SetFoV(sprec.Degrees(60))
	// c.gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	// c.cameraGizmo = viewport.NewCameraGizmo(c.gfxCamera)

	// gridMeshDef := c.commonData.GridMeshDefinition()
	// gridMesh := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: gridMeshDef,
	// })
	// gridMesh.SetMatrix(dprec.IdentityMat4())

	// cameraMeshDef := c.commonData.CameraMeshDefinition()
	// cameraMesh := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: cameraMeshDef,
	// })
	// cameraMesh.SetMatrix(dprec.TranslationMat4(-3.0, 0.1, 0.0))

	// nodeMeshDef := c.commonData.NodeMeshDefinition()
	// nodeMesh := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: nodeMeshDef,
	// })
	// nodeMesh.SetMatrix(dprec.TranslationMat4(-2.0, 0.1, 0.0))

	// ambientLightMeshDef := c.commonData.AmbientLightMeshDefinition()
	// ambientLight := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: ambientLightMeshDef,
	// })
	// ambientLight.SetMatrix(dprec.TranslationMat4(-1.0, 0.1, 0.0))

	// pointLightMeshDef := c.commonData.PointLightMeshDefinition()
	// pointLight := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: pointLightMeshDef,
	// })
	// pointLight.SetMatrix(dprec.TranslationMat4(0.0, 0.1, 0.0))

	// spotLightMeshDef := c.commonData.SpotLightMeshDefinition()
	// spotLight := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: spotLightMeshDef,
	// })
	// spotLight.SetMatrix(dprec.TranslationMat4(1.0, 0.1, 0.0))

	// directionalLightMeshDef := c.commonData.DirectionalLightMeshDefinition()
	// directionalLight := c.gfxScene.CreateMesh(graphics.MeshInfo{
	// 	Definition: directionalLightMeshDef,
	// })
	// directionalLight.SetMatrix(dprec.TranslationMat4(2.0, 0.1, 0.0))
}

func (c *workbenchComponent) OnDelete() {
	// c.gfxScene.Delete()
}

func (c *workbenchComponent) Render() co.Instance {
	return co.New(std.DropZone, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithCallbackData(std.DropZoneCallbackData{
			OnDrop: c.handleDrop,
		})

		co.WithChild("viewport", co.New(std.Viewport, func() {
			co.WithData(std.ViewportData{
				API: c.renderAPI,
			})
			co.WithCallbackData(std.ViewportCallbackData{
				OnKeyboardEvent: c.handleViewportKeyboardEvent,
				OnMouseEvent:    c.handleViewportMouseEvent,
				OnRender:        c.handleViewportRender,
			})
		}))
	})
}

func (c *workbenchComponent) handleDrop(paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	path := paths[0]
	switch ext := filepath.Ext(path); ext {
	case ".glb":
		c.handleDropGLB(path)
		return true
	case ".hdr":
		c.handleDropHDR(path)
		return true
	default:
		common.OpenWarning(c.Scope(), fmt.Sprintf("Unsupported file extension %q", ext))
		return false
	}
}

func (c *workbenchComponent) handleViewportKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	if event.Modifiers.Contains(ui.KeyModifierShift) && event.Code == ui.KeyCodeA {
		if event.Action == ui.KeyboardActionUp {
			c.openAddNodeModal()
		}
		return true
	}

	// return c.cameraGizmo.OnKeyboardEvent(element, event)

	return c.vis.OnKeyboardEvent(element, event)
}

func (c *workbenchComponent) handleViewportMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	// TODO: Do camera motion. Have a "Gadget" concept and pass control initially to it.
	// Then trickle down until you get to here. If no gadget is interested, then do camera motion.

	// return c.cameraGizmo.OnMouseEvent(element, event)

	return c.vis.OnMouseEvent(element, event)
}

func (c *workbenchComponent) handleViewportRender(framebuffer render.Framebuffer, size ui.Size) {
	// c.gfxScene.RenderFramebuffer(framebuffer, graphics.Viewport{
	// 	X:      0,
	// 	Y:      0,
	// 	Width:  size.Width,
	// 	Height: size.Height,
	// })
	c.vis.OnRender(framebuffer, size)
}

func (c *workbenchComponent) openAddNodeModal() {
	co.OpenOverlay(c.Scope(), co.New(AddNodeModal, func() {
		co.WithCallbackData(AddNodeModalCallbackData{
			OnAdd: c.handleAddNode,
		})
	}))
}

func (c *workbenchComponent) handleAddNode(kind editormodel.NodeKind) {
	selectedNode, ok := c.editorModel.Selection().(editormodel.Node)
	if !ok {
		selectedNode = nil
	}

	switch kind {
	case editormodel.NodeKindPointLight:
		c.createPointLight(selectedNode)
	default:
		log.Info("Adding node of kind %v", kind)
	}
}

func (c *workbenchComponent) handleDropGLB(path string) {
	loadingModal := common.OpenLoading(c.Scope())

	promise := async.NewPromise[*pack.Model]()
	go func() {
		if model, err := c.parseGLB(path); err == nil {
			promise.Deliver(model)
		} else {
			promise.Fail(err)
		}
	}()

	promise.OnSuccess(func(model *pack.Model) {
		co.Schedule(c.Scope(), func() {
			loadingModal.Close()
			co.OpenOverlay(c.Scope(), co.New(ModelImport, func() {
				co.WithData(ModelImportData{
					Model: model,
				})
				co.WithCallbackData(ModelImportCallbackData{
					OnImport: c.importModel,
				})
			}))
		})
	})
	promise.OnError(func(err error) {
		co.Schedule(c.Scope(), func() {
			loadingModal.Close()
			common.OpenError(c.Scope(), fmt.Sprintf("Error parsing GLB: %v", err))
		})
	})
}

func (c *workbenchComponent) parseGLB(path string) (*pack.Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	model, err := pack.ParseGLTFResource(file)
	if err != nil {
		return nil, fmt.Errorf("error parsing GLTF: %w", err)
	}

	return model, nil
}

func (c *workbenchComponent) handleDropHDR(path string) {

}

func (c *workbenchComponent) importModel(model *pack.Model) {
	log.Info("Texture count: %d", len(model.Textures))
}

func (c *workbenchComponent) createPointLight(parent editormodel.Node) {
	node := c.editorModel.CreatePointLight(editormodel.PointLightInfo{
		Name: fmt.Sprintf("PointLight %d", rand.Intn(1000)), // FIXME
	})
	node.SetPosition(dprec.Vec3{
		X: rand.Float64()*10.0 - 5.0,
		Y: 2.0, // FIXME
		Z: rand.Float64()*10.0 - 5.0,
	})

	if extendable, ok := parent.(editormodel.ExtendableNode); ok {
		extendable.AppendChild(node)
	} else {
		c.editorModel.AddNode(node)
	}
	// c.editorModel.SetSelection(node)
}
