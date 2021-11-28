package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/data/asset"
)

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

var _ history.Change = (*TwoDTextureWrapS)(nil)

type TwoDTextureWrapS struct {
	Controller model.TwoDTextureEditor

	FromWrap asset.WrapMode
	ToWrap   asset.WrapMode
}

func (ch *TwoDTextureWrapS) Apply() error {
	ch.Controller.SetWrapS(ch.ToWrap)
	return nil
}

func (ch *TwoDTextureWrapS) Revert() error {
	ch.Controller.SetWrapS(ch.FromWrap)
	return nil
}

var _ history.Change = (*TwoDTextureWrapT)(nil)

type TwoDTextureWrapT struct {
	Controller model.TwoDTextureEditor

	FromWrap asset.WrapMode
	ToWrap   asset.WrapMode
}

func (ch *TwoDTextureWrapT) Apply() error {
	ch.Controller.SetWrapT(ch.ToWrap)
	return nil
}

func (ch *TwoDTextureWrapT) Revert() error {
	ch.Controller.SetWrapT(ch.FromWrap)
	return nil
}

var _ history.Change = (*TwoDTextureMinFilter)(nil)

type TwoDTextureMinFilter struct {
	Controller model.TwoDTextureEditor

	FromFilter asset.FilterMode
	ToFilter   asset.FilterMode
}

func (ch *TwoDTextureMinFilter) Apply() error {
	ch.Controller.SetMinFilter(ch.ToFilter)
	return nil
}

func (ch *TwoDTextureMinFilter) Revert() error {
	ch.Controller.SetMinFilter(ch.FromFilter)
	return nil
}

var _ history.Change = (*TwoDTextureMagFilter)(nil)

type TwoDTextureMagFilter struct {
	Controller model.TwoDTextureEditor

	FromFilter asset.FilterMode
	ToFilter   asset.FilterMode
}

func (ch *TwoDTextureMagFilter) Apply() error {
	ch.Controller.SetMagFilter(ch.ToFilter)
	return nil
}

func (ch *TwoDTextureMagFilter) Revert() error {
	ch.Controller.SetMagFilter(ch.FromFilter)
	return nil
}
