package model

import (
	"fmt"
	"image"

	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeTwoDTexture                = mvc.NewChange("twod_texture")
	ChangeTwoDTextureWrapping        = mvc.SubChange(ChangeTwoDTexture, "wrapping")
	ChangeTwoDTextureFiltering       = mvc.SubChange(ChangeTwoDTexture, "filtering")
	ChangeTwoDTextureWidth           = mvc.SubChange(ChangeTwoDTexture, "width")
	ChangeTwoDTextureHeight          = mvc.SubChange(ChangeTwoDTexture, "height")
	ChangeTwoDTextureFormat          = mvc.SubChange(ChangeTwoDTexture, "format")
	ChangeTwoDTextureMipmapping      = mvc.SubChange(ChangeTwoDTexture, "mipmapping")
	ChangeTwoDTextureGammaCorrection = mvc.SubChange(ChangeTwoDTexture, "gamma_correction")
	ChangeTwoDTextureData            = mvc.SubChange(ChangeTwoDTexture, "data")
)

func CreateTwoDTexture(registry *Registry) (*TwoDTexture, error) {
	previewImg := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	previewImg.Pix = []byte{0xFF, 0x00, 0x00, 0xFF}
	texAsset := &asset.TwoDTexture{
		Width:     1,
		Height:    1,
		Wrapping:  asset.WrapModeClampToEdge,
		Filtering: asset.FilterModeNearest,
		Flags:     asset.TextureFlagMipmapping,
		Format:    asset.TexelFormatRGBA8,
		Data:      []byte{0xFF, 0x00, 0x00, 0xFF},
	}

	resourceModel := registry.CreateResource(ResourceKindTwoDTexture, "Unnamed")
	resourceModel.SetPreviewImage(previewImg)
	if err := resourceModel.SaveContent(texAsset); err != nil {
		return nil, fmt.Errorf("error saving content: %w", err)
	}
	if err := resourceModel.Save(); err != nil {
		return nil, err
	}
	return &TwoDTexture{
		Observable:    mvc.NewObservable(),
		resourceModel: resourceModel,
		texAsset:      texAsset,
	}, nil
}

func OpenTwoDTexture(resourceModel *Resource) (*TwoDTexture, error) {
	texAsset := new(asset.TwoDTexture)
	if err := resourceModel.LoadContent(texAsset); err != nil {
		return nil, fmt.Errorf("error loading content: %w", err)
	}
	return &TwoDTexture{
		Observable:    mvc.NewObservable(),
		resourceModel: resourceModel,
		texAsset:      texAsset,
	}, nil
}

type TwoDTexture struct {
	mvc.Observable
	resourceModel *Resource
	texAsset      *asset.TwoDTexture
}

func (t *TwoDTexture) Resource() *Resource {
	return t.resourceModel
}

func (t *TwoDTexture) Wrapping() asset.WrapMode {
	return t.texAsset.Wrapping
}

func (t *TwoDTexture) SetWrapping(wrapping asset.WrapMode) {
	t.texAsset.Wrapping = wrapping
	t.SignalChange(ChangeTwoDTextureWrapping)
}

func (t *TwoDTexture) Filtering() asset.FilterMode {
	return t.texAsset.Filtering
}

func (t *TwoDTexture) SetFiltering(filtering asset.FilterMode) {
	t.texAsset.Filtering = filtering
	t.SignalChange(ChangeTwoDTextureFiltering)
}

func (t *TwoDTexture) Width() int {
	return int(t.texAsset.Width)
}

func (t *TwoDTexture) SetWidth(width int) {
	t.texAsset.Width = uint16(width)
	t.SignalChange(ChangeTwoDTextureWidth)
}

func (t *TwoDTexture) Height() int {
	return int(t.texAsset.Height)
}

func (t *TwoDTexture) SetHeight(height int) {
	t.texAsset.Height = uint16(height)
	t.SignalChange(ChangeTwoDTextureHeight)
}

func (t *TwoDTexture) Format() asset.TexelFormat {
	return t.texAsset.Format
}

func (t *TwoDTexture) SetFormat(format asset.TexelFormat) {
	t.texAsset.Format = format
	t.SignalChange(ChangeTwoDTextureFormat)
}

func (t *TwoDTexture) Mipmapping() bool {
	return t.texAsset.Flags.Has(asset.TextureFlagMipmapping)
}

func (t *TwoDTexture) SetMipmapping(mipmapping bool) {
	if mipmapping {
		t.texAsset.Flags |= asset.TextureFlagMipmapping
	} else {
		t.texAsset.Flags &= ^asset.TextureFlagMipmapping
	}
	t.SignalChange(ChangeTwoDTextureMipmapping)
}

func (t *TwoDTexture) GammaCorrection() bool {
	return !t.texAsset.Flags.Has(asset.TextureFlagLinear)
}

func (t *TwoDTexture) SetGammaCorrection(correction bool) {
	if correction {
		t.texAsset.Flags &= ^asset.TextureFlagLinear
	} else {
		t.texAsset.Flags |= asset.TextureFlagLinear
	}
	t.SignalChange(ChangeTwoDTextureGammaCorrection)
}

func (t *TwoDTexture) Data() []byte {
	return t.texAsset.Data
}

func (t *TwoDTexture) SetData(data []byte) {
	t.texAsset.Data = data
	t.SignalChange(ChangeTwoDTextureData)
}

func (t *TwoDTexture) Save() error {
	if err := t.resourceModel.SaveContent(t.texAsset); err != nil {
		return fmt.Errorf("error saving content: %w", err)
	}
	return t.resourceModel.Save()
}
