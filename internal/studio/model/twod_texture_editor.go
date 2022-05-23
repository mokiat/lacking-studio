package model

import (
	"github.com/mokiat/lacking/game/asset"
)

type TwoDTextureEditor interface {
	Editor

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	Wrapping() asset.WrapMode
	Filtering() asset.FilterMode
	DataFormat() asset.TexelFormat

	ChangeName(name string)
	ChangeContent(path string)
	ChangeWrapping(wrap asset.WrapMode)
	ChangeFiltering(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	Visualization() Visualization
}
