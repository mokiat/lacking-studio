package model

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/game/graphics"
)

type TwoDTextureEditor interface {
	Editor

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	SetWrapS(wrap asset.WrapMode)
	WrapS() asset.WrapMode
	SetWrapT(wrap asset.WrapMode)
	WrapT() asset.WrapMode
	SetAssetData(data asset.TwoDTexture)
	MinFilter() asset.FilterMode
	SetMinFilter(filter asset.FilterMode)
	MagFilter() asset.FilterMode
	SetMagFilter(filter asset.FilterMode)
	DataFormat() asset.TexelFormat

	ChangeSourcePath(path string)
	ChangeWrapS(wrap asset.WrapMode)
	ChangeWrapT(wrap asset.WrapMode)
	ChangeMinFilter(filter asset.FilterMode)
	ChangeMagFilter(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	Update()
	Scene() graphics.Scene
	Camera() graphics.Camera
	OnViewportMouseEvent(event widget.ViewportMouseEvent) bool
}
