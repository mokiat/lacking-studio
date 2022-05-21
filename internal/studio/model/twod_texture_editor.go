package model

import (
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui/mat"
)

type TwoDTextureEditor interface {
	Editor

	SetName(name string)

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	Wrapping() asset.WrapMode
	SetWrapping(wrap asset.WrapMode)

	Filtering() asset.FilterMode
	SetFiltering(filter asset.FilterMode)

	DataFormat() asset.TexelFormat
	SetAssetData(data asset.TwoDTexture)

	ChangeName(name string)
	ChangeSourcePath(path string)
	ChangeWrapping(wrap asset.WrapMode)
	ChangeFiltering(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	Update()
	Scene() *graphics.Scene
	Camera() *graphics.Camera
	OnViewportMouseEvent(event mat.ViewportMouseEvent) bool
}
