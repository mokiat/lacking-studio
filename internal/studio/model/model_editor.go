package model

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui/mat"
)

type ModelEditor interface {
	Editor

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)

	Update()
	Scene() *graphics.Scene
	Camera() *graphics.Camera
	OnViewportMouseEvent(event mat.ViewportMouseEvent) bool
}
