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
	// NOTE: This allows for example the studio to subscribe for name changes only
	// for an editor.
	NameChange = observer.StringChange("name")

	TwoDTextureChange                = observer.StringChange("twod_texture")
	TwoDTextureNameChange            = observer.ExtendChange(TwoDTextureChange, NameChange)
	TwoDTextureWrappingChange        = observer.ExtendChange(TwoDTextureChange, observer.StringChange("wrapping"))
	TwoDTextureFilteringChange       = observer.ExtendChange(TwoDTextureChange, observer.StringChange("filtering"))
	TwoDTextureWidthChange           = observer.ExtendChange(TwoDTextureChange, observer.StringChange("width"))
	TwoDTextureHeightChange          = observer.ExtendChange(TwoDTextureChange, observer.StringChange("height"))
	TwoDTextureFormatChange          = observer.ExtendChange(TwoDTextureChange, observer.StringChange("format"))
	TwoDTextureMipmappingChange      = observer.ExtendChange(TwoDTextureChange, observer.StringChange("mipmapping"))
	TwoDTextureGammaCorrectionChange = observer.ExtendChange(TwoDTextureChange, observer.StringChange("gamma_correction"))
	TwoDTextureDataChange            = observer.ExtendChange(TwoDTextureChange, observer.StringChange("data"))
	TwoDTexturePreviewChange         = observer.ExtendChange(TwoDTextureChange, observer.StringChange("preview"))
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
		target:     observer.NewTarget(),
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
		target:     observer.NewTarget(),
		resource:   resource,
		texAsset:   texAsset,
		previewImg: previewImg,
	}, nil
}

type TwoDTexture struct {
	target     *observer.Target
	resource   *data.Resource
	texAsset   *asset.TwoDTexture
	previewImg image.Image
}

func (t *TwoDTexture) Target() *observer.Target {
	return t.target
}

func (t *TwoDTexture) ID() string {
	return t.resource.ID()
}

func (t *TwoDTexture) Name() string {
	return t.resource.Name()
}

func (t *TwoDTexture) SetName(name string) {
	t.resource.SetName(name)
	t.target.SignalChange(TwoDTextureNameChange)
}

func (t *TwoDTexture) Wrapping() asset.WrapMode {
	return t.texAsset.Wrapping
}

func (t *TwoDTexture) SetWrapping(wrapping asset.WrapMode) {
	t.texAsset.Wrapping = wrapping
	t.target.SignalChange(TwoDTextureWrappingChange)
}

func (t *TwoDTexture) Filtering() asset.FilterMode {
	return t.texAsset.Filtering
}

func (t *TwoDTexture) SetFiltering(filtering asset.FilterMode) {
	t.texAsset.Filtering = filtering
	t.target.SignalChange(TwoDTextureFilteringChange)
}

func (t *TwoDTexture) Width() int {
	return int(t.texAsset.Width)
}

func (t *TwoDTexture) SetWidth(width int) {
	t.texAsset.Width = uint16(width)
	t.target.SignalChange(TwoDTextureWidthChange)
}

func (t *TwoDTexture) Height() int {
	return int(t.texAsset.Height)
}

func (t *TwoDTexture) SetHeight(height int) {
	t.texAsset.Height = uint16(height)
	t.target.SignalChange(TwoDTextureHeightChange)
}

func (t *TwoDTexture) Format() asset.TexelFormat {
	return t.texAsset.Format
}

func (t *TwoDTexture) SetFormat(format asset.TexelFormat) {
	t.texAsset.Format = format
	t.target.SignalChange(TwoDTextureFormatChange)
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
	t.target.SignalChange(TwoDTextureMipmappingChange)
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
	t.target.SignalChange(TwoDTextureGammaCorrectionChange)
}

func (t *TwoDTexture) Data() []byte {
	return t.texAsset.Data
}

func (t *TwoDTexture) SetData(data []byte) {
	t.texAsset.Data = data
	t.target.SignalChange(TwoDTextureDataChange)
}

func (t *TwoDTexture) PreviewImage() image.Image {
	return t.previewImg
}

func (t *TwoDTexture) SetPreviewImage(img image.Image) {
	t.previewImg = img
	t.target.SignalChange(TwoDTexturePreviewChange)
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
