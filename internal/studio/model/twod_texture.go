package model

import (
	"errors"
	"fmt"
	"image"

	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking/game/asset"
)

var (
	TwoDTextureChangeToName            observer.Change = "twod_texture:name"
	TwoDTextureChangeToWrapping        observer.Change = "twod_texture:wrapping"
	TwoDTextureChangeToFiltering       observer.Change = "twod_texture:filtering"
	TwoDTextureChangeToWidth           observer.Change = "twod_texture:width"
	TwoDTextureChangeToHeight          observer.Change = "twod_texture:height"
	TwoDTextureChangeToFormat          observer.Change = "twod_texture:format"
	TwoDTextureChangeToMipmapping      observer.Change = "twod_texture:mipmapping"
	TwoDTextureChangeToGammaCorrection observer.Change = "twod_texture:gamme_correction"
	TwoDTextureChangeToData            observer.Change = "twod_texture:data"
	TwoDTextureChangeToPreview         observer.Change = "twod_texture:preview"
)

func CreateTwoDTexture(registry *data.Registry) (*TwoDTexture, error) {
	resource := registry.CreateResource(data.ResourceKindTwoDTexture)
	if err := resource.Save(); err != nil {
		return nil, fmt.Errorf("error saving resource: %w", err)
	}
	content := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	if err := resource.SavePreview(content); err != nil {
		return nil, fmt.Errorf("error saving preview: %w", err)
	}
	texAsset := &asset.TwoDTexture{
		Width:     1,
		Height:    1,
		Wrapping:  asset.WrapModeClampToEdge,
		Filtering: asset.FilterModeNearest,
		Flags:     asset.TextureFlagMipmapping,
		Format:    asset.TexelFormatRGBA8,
		Data:      content.Pix,
	}
	if err := resource.SaveContent(texAsset); err != nil {
		return nil, fmt.Errorf("error saving content: %w", err)
	}
	return &TwoDTexture{
		resource:   resource,
		texAsset:   texAsset,
		previewImg: content,
	}, nil
}

func OpenTwoDTexture(registry *data.Registry, id string) (*TwoDTexture, error) {
	resource := registry.GetResourceByID(id)
	texAsset := new(asset.TwoDTexture)
	if err := resource.LoadContent(texAsset); err != nil {
		return nil, fmt.Errorf("error loading content: %w", err)
	}
	previewImg, err := resource.LoadPreview()
	if err != nil {
		if !errors.Is(err, asset.ErrNotFound) {
			return nil, fmt.Errorf("error loading preview: %w", err)
		}
		previewImg = nil
	}
	return &TwoDTexture{
		resource:   resource,
		texAsset:   texAsset,
		previewImg: previewImg,
	}, nil
}

type TwoDTexture struct {
	observer.Target
	resource   *data.Resource
	texAsset   *asset.TwoDTexture
	previewImg image.Image
}

func (t *TwoDTexture) ID() string {
	return t.resource.ID()
}

func (t *TwoDTexture) Name() string {
	return t.resource.Name()
}

func (t *TwoDTexture) SetName(name string) {
	t.resource.SetName(name)
	t.SignalChange(TwoDTextureChangeToName)
}

func (t *TwoDTexture) Wrapping() asset.WrapMode {
	return t.texAsset.Wrapping
}

func (t *TwoDTexture) SetWrapping(wrapping asset.WrapMode) {
	t.texAsset.Wrapping = wrapping
	t.SignalChange(TwoDTextureChangeToWrapping)
}

func (t *TwoDTexture) Filtering() asset.FilterMode {
	return t.texAsset.Filtering
}

func (t *TwoDTexture) SetFiltering(filtering asset.FilterMode) {
	t.texAsset.Filtering = filtering
	t.SignalChange(TwoDTextureChangeToFiltering)
}

func (t *TwoDTexture) Width() int {
	return int(t.texAsset.Width)
}

func (t *TwoDTexture) SetWidth(width int) {
	t.texAsset.Width = uint16(width)
	t.SignalChange(TwoDTextureChangeToWidth)
}

func (t *TwoDTexture) Height() int {
	return int(t.texAsset.Height)
}

func (t *TwoDTexture) SetHeight(height int) {
	t.texAsset.Height = uint16(height)
	t.SignalChange(TwoDTextureChangeToHeight)
}

func (t *TwoDTexture) Format() asset.TexelFormat {
	return t.texAsset.Format
}

func (t *TwoDTexture) SetFormat(format asset.TexelFormat) {
	t.texAsset.Format = format
	t.SignalChange(TwoDTextureChangeToFormat)
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
	t.SignalChange(TwoDTextureChangeToMipmapping)
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
	t.SignalChange(TwoDTextureChangeToGammaCorrection)
}

func (t *TwoDTexture) Data() []byte {
	return t.texAsset.Data
}

func (t *TwoDTexture) SetData(data []byte) {
	t.texAsset.Data = data
	t.SignalChange(TwoDTextureChangeToData)
}

func (t *TwoDTexture) PreviewImage() image.Image {
	return t.previewImg
}

func (t *TwoDTexture) SetPreviewImage(img image.Image) {
	t.previewImg = img
	t.SignalChange(TwoDTextureChangeToPreview)
}

func (t *TwoDTexture) Save() error {
	if err := t.resource.Save(); err != nil {
		return fmt.Errorf("error saving resource: %w", err)
	}
	if err := t.resource.SavePreview(t.previewImg); err != nil {
		return fmt.Errorf("error saving preview: %w", err)
	}
	if err := t.resource.SaveContent(t.texAsset); err != nil {
		return fmt.Errorf("error saving content: %w", err)
	}
	return nil
}
