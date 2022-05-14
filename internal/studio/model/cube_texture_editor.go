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

	SetAssetData(data asset.CubeTexture)
	MinFilter() asset.FilterMode
	SetMinFilter(filter asset.FilterMode)
	MagFilter() asset.FilterMode
	SetMagFilter(filter asset.FilterMode)
	DataFormat() asset.TexelFormat

	ChangeSourcePath(path string)
	ChangeMinFilter(filter asset.FilterMode)
	ChangeMagFilter(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	OnViewportMouseEvent(event mat.ViewportMouseEvent) bool
	Update()
	Scene() *graphics.Scene
	Camera() *graphics.Camera
}
