package action

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
)

type ChangeTwoDTextureWrapping struct {
	Texture  *model.TwoDTexture
	Wrapping asset.WrapMode
}

type ChangeTwoDTextureFiltering struct {
	Texture   *model.TwoDTexture
	Filtering asset.FilterMode
}

type ChangeTwoDTextureFormat struct {
	Texture *model.TwoDTexture
	Format  asset.TexelFormat
}

type ChangeTwoDTextureContentFromPath struct {
	Texture *model.TwoDTexture
	Path    string
}

type ChangeCubeTextureFiltering struct {
	Texture   *model.CubeTexture
	Filtering asset.FilterMode
}

type ChangeCubeTextureFormat struct {
	Texture *model.CubeTexture
	Format  asset.TexelFormat
}

type ChangeCubeTextureContentFromPath struct {
	Texture *model.CubeTexture
	Path    string
}
