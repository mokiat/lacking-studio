package model

import (
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui/mat"
)

type CubeTextureEditor interface {
	Editor

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	Filtering() asset.FilterMode
	SetFiltering(filter asset.FilterMode)

	DataFormat() asset.TexelFormat
	SetAssetData(data asset.CubeTexture)

	ChangeSourcePath(path string)
	ChangeFiltering(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	OnViewportMouseEvent(event mat.ViewportMouseEvent) bool
	Update()
	Scene() *graphics.Scene
	Camera() *graphics.Camera
}
