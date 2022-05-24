package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
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
