package model

import (
	"github.com/mokiat/lacking/data/asset"
)

type CubeTextureEditor interface {
	Editor

	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	SetAssetData(data asset.CubeTexture)
	SetAssetMinFilter(filter asset.FilterMode)
	SetAssetMagFilter(filter asset.FilterMode)

	ChangeSourcePath(path string)
}
