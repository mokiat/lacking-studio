package editor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mokiat/lacking-studio/internal/view/common"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/lacking/util/async"
)

var Workbench = co.Define(&workbenchComponent{})

type WorkbenchData struct{}

type workbenchComponent struct {
	co.BaseComponent

	renderAPI render.API
}

func (c *workbenchComponent) OnCreate() {
	window := co.Window(c.Scope())
	c.renderAPI = window.RenderAPI()
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
		c.loadGLB(path)
		return true
	case ".hdr":
		c.loadHDR(path)
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
	c.renderAPI.BeginRenderPass(render.RenderPassInfo{
		Framebuffer: framebuffer,
		Viewport: render.Area{
			X:      0,
			Y:      0,
			Width:  size.Width,
			Height: size.Height,
		},
		DepthLoadOp:     render.LoadOperationClear,
		DepthStoreOp:    render.StoreOperationStore,
		DepthClearValue: 1.0,
		StencilLoadOp:   render.LoadOperationDontCare,
		StencilStoreOp:  render.StoreOperationDontCare,
		Colors: [4]render.ColorAttachmentInfo{
			{
				LoadOp:     render.LoadOperationClear,
				StoreOp:    render.StoreOperationStore,
				ClearValue: [4]float32{0.1, 0.3, 0.5, 1.0},
			},
		},
	})

	c.renderAPI.EndRenderPass()
}

func (c *workbenchComponent) loadGLB(path string) {
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
			log.Info("Textures: %d", len(model.Textures))
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

func (c *workbenchComponent) loadHDR(path string) {

}
