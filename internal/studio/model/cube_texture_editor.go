package model

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking/game/asset"
)

type CubeTextureEditor interface {
	Editor

	Target() *observer.Target

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	Filtering() asset.FilterMode
	DataFormat() asset.TexelFormat

	ChangeName(name string)
	ChangeContent(path string)
	ChangeFiltering(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	Visualization() Visualization
}
