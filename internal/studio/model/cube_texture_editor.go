package model

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/game/graphics"
)

type CubeTextureEditor interface {
	Editor

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	SetAssetData(data asset.CubeTexture)
	SetAssetMinFilter(filter asset.FilterMode)
	SetAssetMagFilter(filter asset.FilterMode)

	ChangeSourcePath(path string)

	OnViewportMouseEvent(event widget.ViewportMouseEvent) bool
	Update()
	Scene() graphics.Scene
	Camera() graphics.Camera
}
