package model

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/graphics"
)

type ModelEditor interface {
	Editor

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)

	Update()
	Scene() *graphics.Scene
	Camera() *graphics.Camera
	OnViewportMouseEvent(event widget.ViewportMouseEvent) bool
}
