package editor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	glgame "github.com/mokiat/lacking-native/game"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking-studio/internal/view/common"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/graphics/shading"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/lacking/util/async"
	"github.com/mokiat/lacking/util/blob"
)

var Workbench = co.Define(&workbenchComponent{})

type WorkbenchData struct {
	EditorModel *editormodel.Model
}

type workbenchComponent struct {
	co.BaseComponent

	renderAPI render.API

	gfxEngine *graphics.Engine
	gfxScene  *graphics.Scene
	gfxCamera *graphics.Camera
}

func (c *workbenchComponent) OnCreate() {
	window := co.Window(c.Scope())
	c.renderAPI = window.RenderAPI()

	c.gfxEngine = graphics.NewEngine(c.renderAPI, glgame.NewShaderCollection())
	c.gfxEngine.Create()

	c.gfxScene = c.gfxEngine.CreateScene()
	c.gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.01, 0.01, 0.02))

	c.gfxCamera = c.gfxScene.CreateCamera()
	c.gfxCamera.SetMatrix(
		dprec.Mat4MultiProd(
			dprec.RotationMat4(dprec.Degrees(15), 0.0, 1.0, 0.0),
			dprec.RotationMat4(dprec.Degrees(15), -1.0, 0.0, 0.0),
			dprec.TranslationMat4(0.0, 0.0, 10.0),
		),
	)
	c.gfxCamera.SetExposure(1.0)
	c.gfxCamera.SetAutoExposure(false)
	c.gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)

	gridMeshDef := createGridMeshDefinition(c.gfxEngine)

	gridMesh := c.gfxScene.CreateMesh(graphics.MeshInfo{
		Definition: gridMeshDef,
		Armature:   nil,
	})
	gridMesh.SetMatrix(dprec.IdentityMat4())
}

func (c *workbenchComponent) OnDelete() {
	c.gfxScene.Delete()
	c.gfxEngine.Destroy()
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

func (c *workbenchComponent) handleViewportKeyboardEvent(event ui.KeyboardEvent) bool {
	return false
}

func (c *workbenchComponent) handleViewportMouseEvent(event std.ViewportMouseEvent) bool {
	return false
}

func (c *workbenchComponent) handleViewportRender(framebuffer render.Framebuffer, size ui.Size) {
	// if false {
	// 	c.gfxEngine.Debug().Reset()
	// 	c.gfxEngine.Debug().Line(
	// 		dprec.ZeroVec3(),
	// 		dprec.NewVec3(1000.0, 0.0, 0.0),
	// 		dprec.NewVec3(1.0, 0.0, 0.0),
	// 	)
	// 	c.gfxEngine.Debug().Line(
	// 		dprec.ZeroVec3(),
	// 		dprec.NewVec3(-1000.0, 0.0, 0.0),
	// 		dprec.NewVec3(0.3, 0.0, 0.0),
	// 	)
	// 	c.gfxEngine.Debug().Line(
	// 		dprec.ZeroVec3(),
	// 		dprec.NewVec3(0.0, 0.0, 1000.0),
	// 		dprec.NewVec3(0.0, 0.0, 1.0),
	// 	)
	// 	c.gfxEngine.Debug().Line(
	// 		dprec.ZeroVec3(),
	// 		dprec.NewVec3(0.0, 0.0, -1000.0),
	// 		dprec.NewVec3(0.0, 0.0, 0.3),
	// 	)
	// 	const distance = 10
	// 	for i := 1; i <= distance; i++ {
	// 		c.gfxEngine.Debug().Line(
	// 			dprec.NewVec3(-float64(i), 0.0, -distance),
	// 			dprec.NewVec3(-float64(i), 0.0, distance),
	// 			dprec.NewVec3(0.3, 0.3, 0.3),
	// 		)
	// 		c.gfxEngine.Debug().Line(
	// 			dprec.NewVec3(float64(i), 0.0, -distance),
	// 			dprec.NewVec3(float64(i), 0.0, distance),
	// 			dprec.NewVec3(0.3, 0.3, 0.3),
	// 		)
	// 		c.gfxEngine.Debug().Line(
	// 			dprec.NewVec3(-distance, 0.0, -float64(i)),
	// 			dprec.NewVec3(distance, 0.0, -float64(i)),
	// 			dprec.NewVec3(0.3, 0.3, 0.3),
	// 		)
	// 		c.gfxEngine.Debug().Line(
	// 			dprec.NewVec3(-distance, 0.0, float64(i)),
	// 			dprec.NewVec3(distance, 0.0, float64(i)),
	// 			dprec.NewVec3(0.3, 0.3, 0.3),
	// 		)
	// 	}
	// }

	c.gfxScene.RenderFramebuffer(framebuffer, graphics.Viewport{
		X:      0,
		Y:      0,
		Width:  size.Width,
		Height: size.Height,
	})
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

func createGridMeshDefinition(gfxEngine *graphics.Engine) *graphics.MeshDefinition {
	gridShading := gfxEngine.CreateShading(graphics.ShadingInfo{
		ShadowFunc:   nil,
		GeometryFunc: nil,
		ForwardFunc: func(palette *shading.ForwardPalette) {
			color := palette.ConstVec4(0.0, 0.0, 1.0, 1.0)
			color = palette.MulVec4(color, 0.1)
			palette.OutputColor(color)
		},
		EmissiveFunc: nil,
		LightingFunc: nil,
	})

	gridMaterialDef := gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		BackfaceCulling: false,
		AlphaTesting:    false,
		AlphaBlending:   false,
		AlphaThreshold:  0.0,
		Vectors:         []sprec.Vec4{},
		TwoDTextures:    []render.Texture{},
		CubeTextures:    []render.Texture{},
		Shading:         gridShading,
	})

	// gridMeshBuilder := graphics.NewMeshBuilder(
	// 	graphics.MeshBuilderWithCoords(),
	// )
	// gridMeshBuilder.BuildInfo() // TODO

	// TODO: Draw actual lines, not a quad.
	vertexData := blob.NewPlotter(make([]byte, 4*3*4))
	vertexData.PlotSPVec3(sprec.NewVec3(-5.0, 5.0, 0.0))
	vertexData.PlotSPVec3(sprec.NewVec3(-5.0, -5.0, 0.0))
	vertexData.PlotSPVec3(sprec.NewVec3(5.0, -5.0, 0.0))
	vertexData.PlotSPVec3(sprec.NewVec3(5.0, 5.0, 0.0))

	indexData := blob.NewPlotter(make([]byte, 6*2))
	indexData.PlotUint16(0)
	indexData.PlotUint16(1)
	indexData.PlotUint16(2)
	indexData.PlotUint16(0)
	indexData.PlotUint16(2)
	indexData.PlotUint16(3)

	return gfxEngine.CreateMeshDefinition(graphics.MeshDefinitionInfo{
		VertexData: vertexData.Data(),
		VertexFormat: graphics.VertexFormat{
			HasCoord: true,
		},
		IndexData:   indexData.Data(),
		IndexFormat: graphics.IndexFormatU16,
		Fragments: []graphics.MeshFragmentDefinitionInfo{
			{
				Primitive:   graphics.PrimitiveTriangles,
				IndexOffset: 0,
				IndexCount:  8,
				Material:    gridMaterialDef,
			},
		},
	})
}
