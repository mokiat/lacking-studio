package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
)

var _ history.Change = (*TwoDTextureName)(nil)

type TwoDTextureName struct {
	Controller model.TwoDTextureEditor

	From string
	To   string
}

func (ch *TwoDTextureName) Apply() error {
	ch.Controller.SetName(ch.To)
	return nil
}

func (ch *TwoDTextureName) Revert() error {
	ch.Controller.SetName(ch.From)
	return nil
}

var _ history.Change = (*TwoDTextureData)(nil)

type TwoDTextureData struct {
	Controller model.TwoDTextureEditor

	FromAsset asset.TwoDTexture
	ToAsset   asset.TwoDTexture
}

func (ch *TwoDTextureData) Apply() error {
	ch.Controller.SetAssetData(ch.ToAsset)
	return nil
}

func (ch *TwoDTextureData) Revert() error {
	ch.Controller.SetAssetData(ch.FromAsset)
	return nil
}

var _ history.Change = (*TwoDTextureWrapping)(nil)

type TwoDTextureWrapping struct {
	Controller model.TwoDTextureEditor

	FromWrap asset.WrapMode
	ToWrap   asset.WrapMode
}

func (ch *TwoDTextureWrapping) Apply() error {
	ch.Controller.SetWrapping(ch.ToWrap)
	return nil
}

func (ch *TwoDTextureWrapping) Revert() error {
	ch.Controller.SetWrapping(ch.FromWrap)
	return nil
}

var _ history.Change = (*TwoDTextureFiltering)(nil)

type TwoDTextureFiltering struct {
	Controller model.TwoDTextureEditor

	FromFilter asset.FilterMode
	ToFilter   asset.FilterMode
}

func (ch *TwoDTextureFiltering) Apply() error {
	ch.Controller.SetFiltering(ch.ToFilter)
	return nil
}

func (ch *TwoDTextureFiltering) Revert() error {
	ch.Controller.SetFiltering(ch.FromFilter)
	return nil
}
