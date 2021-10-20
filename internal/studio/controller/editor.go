package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
)

type Editor interface {
	ID() string
	Name() string
	Icon() ui.Image

	// CanUndo() bool
	// Undo()
	// CanRedo() bool
	// Redo()

	// CanSave() bool
	// Save()

	Update()
	OnViewportMouseEvent(event widget.ViewportMouseEvent)

	Scene() graphics.Scene
	Camera() graphics.Camera

	RenderProperties() co.Instance

	Destroy()
}
