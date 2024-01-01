package editor

import (
	"fmt"
	"path/filepath"

	"github.com/mokiat/lacking-studio/internal/view/common"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Workbench = co.Define(&workbenchComponent{})

type WorkbenchData struct{}

type workbenchComponent struct {
	co.BaseComponent

	renderAPI render.API
}

func (c *workbenchComponent) OnCreate() {
	c.renderAPI = co.Window(c.Scope()).RenderAPI()
}

func (c *workbenchComponent) Render() co.Instance {
	window := co.Window(c.Scope())
	renderAPI := window.RenderAPI()
	if renderAPI == nil {
		panic("no render API")
	}

	return co.New(std.DropZone, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithCallbackData(std.DropZoneCallbackData{
			OnDrop: c.handleDrop,
		})

		co.WithChild("viewport", co.New(std.Viewport, func() {
			co.WithData(std.ViewportData{
				API: renderAPI,
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
		return true
	case ".hdr":
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
