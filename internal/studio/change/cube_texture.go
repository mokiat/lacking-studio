package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/data/asset"
)

var _ history.Change = (*CubeTextureData)(nil)

type CubeTextureData struct {
	Controller model.CubeTextureEditor

	FromAsset asset.CubeTexture
	ToAsset   asset.CubeTexture
}

func (ch *CubeTextureData) Apply() error {
	ch.Controller.SetAssetData(ch.ToAsset)
	return nil
}

func (ch *CubeTextureData) Revert() error {
	ch.Controller.SetAssetData(ch.FromAsset)
	return nil
}

var _ history.Change = (*CubeTextureMinFilter)(nil)

type CubeTextureMinFilter struct {
	Controller model.CubeTextureEditor

	FromFilter asset.FilterMode
	ToFilter   asset.FilterMode
}

func (ch *CubeTextureMinFilter) Apply() error {
	ch.Controller.SetMinFilter(ch.ToFilter)
	return nil
}

func (ch *CubeTextureMinFilter) Revert() error {
	ch.Controller.SetMinFilter(ch.FromFilter)
	return nil
}

var _ history.Change = (*CubeTextureMagFilter)(nil)

type CubeTextureMagFilter struct {
	Controller model.CubeTextureEditor

	FromFilter asset.FilterMode
	ToFilter   asset.FilterMode
}

func (ch *CubeTextureMagFilter) Apply() error {
	ch.Controller.SetMagFilter(ch.ToFilter)
	return nil
}

func (ch *CubeTextureMagFilter) Revert() error {
	ch.Controller.SetMagFilter(ch.FromFilter)
	return nil
}
