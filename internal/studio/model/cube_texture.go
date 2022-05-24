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
	CubeTextureChange                = observer.StringChange("twod_texture")
	CubeTextureNameChange            = observer.ExtendChange(CubeTextureChange, NameChange)
	CubeTextureFilteringChange       = observer.ExtendChange(CubeTextureChange, observer.StringChange("filtering"))
	CubeTextureDimensionChange       = observer.ExtendChange(CubeTextureChange, observer.StringChange("dimension"))
	CubeTextureFormatChange          = observer.ExtendChange(CubeTextureChange, observer.StringChange("format"))
	CubeTextureMipmappingChange      = observer.ExtendChange(CubeTextureChange, observer.StringChange("mipmapping"))
	CubeTextureGammaCorrectionChange = observer.ExtendChange(CubeTextureChange, observer.StringChange("gamma_correction"))
	CubeTextureDataChange            = observer.ExtendChange(CubeTextureChange, observer.StringChange("data"))
	CubeTexturePreviewChange         = observer.ExtendChange(CubeTextureChange, observer.StringChange("preview"))
)

func CreateCubeTexture(registry *data.Registry) (*CubeTexture, error) {
	resource := registry.CreateResource(data.ResourceKindCubeTexture)
	if err := resource.Save(); err != nil {
		return nil, fmt.Errorf("error saving resource: %w", err)
	}
	previewImg := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	if err := resource.SavePreview(previewImg); err != nil {
		return nil, fmt.Errorf("error saving preview: %w", err)
	}
	texAsset := &asset.CubeTexture{
		Dimension: 1,
		Filtering: asset.FilterModeNearest,
		Flags:     asset.TextureFlagNone,
		Format:    asset.TexelFormatRGBA8,
		FrontSide: asset.CubeTextureSide{
			Data: []byte{0xFF, 0x00, 0x00, 0xFF},
		},
		BackSide: asset.CubeTextureSide{
			Data: []byte{0x00, 0xFF, 0x00, 0xFF},
		},
		LeftSide: asset.CubeTextureSide{
			Data: []byte{0x00, 0x00, 0xFF, 0xFF},
		},
		RightSide: asset.CubeTextureSide{
			Data: []byte{0xFF, 0xFF, 0x00, 0xFF},
		},
		TopSide: asset.CubeTextureSide{
			Data: []byte{0xFF, 0x00, 0xFF, 0xFF},
		},
		BottomSide: asset.CubeTextureSide{
			Data: []byte{0x00, 0xFF, 0xFF, 0xFF},
		},
	}
	if err := resource.SaveContent(texAsset); err != nil {
		return nil, fmt.Errorf("error saving content: %w", err)
	}
	return &CubeTexture{
		target:     observer.NewTarget(),
		resource:   resource,
		texAsset:   texAsset,
		previewImg: previewImg,
	}, nil
}

func OpenCubeTexture(registry *data.Registry, id string) (*CubeTexture, error) {
	resource := registry.GetResourceByID(id)
	texAsset := new(asset.CubeTexture)
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
	return &CubeTexture{
		target:     observer.NewTarget(),
		resource:   resource,
		texAsset:   texAsset,
		previewImg: previewImg,
	}, nil
}

type CubeTexture struct {
	target     *observer.Target
	resource   *data.Resource
	texAsset   *asset.CubeTexture
	previewImg image.Image
}

func (t *CubeTexture) Target() *observer.Target {
	return t.target
}

func (t *CubeTexture) ID() string {
	return t.resource.ID()
}

func (t *CubeTexture) Name() string {
	return t.resource.Name()
}

func (t *CubeTexture) SetName(name string) {
	t.resource.SetName(name)
	t.target.SignalChange(CubeTextureNameChange)
}

func (t *CubeTexture) Kind() data.ResourceKind {
	return t.resource.Kind()
}

func (t *CubeTexture) Filtering() asset.FilterMode {
	return t.texAsset.Filtering
}

func (t *CubeTexture) SetFiltering(filtering asset.FilterMode) {
	t.texAsset.Filtering = filtering
	t.target.SignalChange(CubeTextureFilteringChange)
}

func (t *CubeTexture) Dimension() int {
	return int(t.texAsset.Dimension)
}

func (t *CubeTexture) SetDimension(dimesion int) {
	t.texAsset.Dimension = uint16(dimesion)
	t.target.SignalChange(CubeTextureDimensionChange)
}

func (t *CubeTexture) Format() asset.TexelFormat {
	return t.texAsset.Format
}

func (t *CubeTexture) SetFormat(format asset.TexelFormat) {
	t.texAsset.Format = format
	t.target.SignalChange(CubeTextureFormatChange)
}

func (t *CubeTexture) Mipmapping() bool {
	return t.texAsset.Flags.Has(asset.TextureFlagMipmapping)
}

func (t *CubeTexture) SetMipmapping(mipmapping bool) {
	if mipmapping {
		t.texAsset.Flags |= asset.TextureFlagMipmapping
	} else {
		t.texAsset.Flags &= ^asset.TextureFlagMipmapping
	}
	t.target.SignalChange(CubeTextureMipmappingChange)
}

func (t *CubeTexture) GammaCorrection() bool {
	return !t.texAsset.Flags.Has(asset.TextureFlagLinear)
}

func (t *CubeTexture) SetGammaCorrection(correction bool) {
	if correction {
		t.texAsset.Flags &= ^asset.TextureFlagLinear
	} else {
		t.texAsset.Flags |= asset.TextureFlagLinear
	}
	t.target.SignalChange(CubeTextureGammaCorrectionChange)
}

func (t *CubeTexture) FrontData() []byte {
	return t.texAsset.FrontSide.Data
}

func (t *CubeTexture) SetFrontData(data []byte) {
	t.texAsset.FrontSide.Data = data
	t.target.SignalChange(CubeTextureDataChange)
}

func (t *CubeTexture) BackData() []byte {
	return t.texAsset.BackSide.Data
}

func (t *CubeTexture) SetBackData(data []byte) {
	t.texAsset.BackSide.Data = data
	t.target.SignalChange(CubeTextureDataChange)
}

func (t *CubeTexture) LeftData() []byte {
	return t.texAsset.LeftSide.Data
}

func (t *CubeTexture) SetLeftData(data []byte) {
	t.texAsset.LeftSide.Data = data
	t.target.SignalChange(CubeTextureDataChange)
}

func (t *CubeTexture) RightData() []byte {
	return t.texAsset.RightSide.Data
}

func (t *CubeTexture) SetRightData(data []byte) {
	t.texAsset.RightSide.Data = data
	t.target.SignalChange(CubeTextureDataChange)
}

func (t *CubeTexture) TopData() []byte {
	return t.texAsset.TopSide.Data
}

func (t *CubeTexture) SetTopData(data []byte) {
	t.texAsset.TopSide.Data = data
	t.target.SignalChange(CubeTextureDataChange)
}

func (t *CubeTexture) BottomData() []byte {
	return t.texAsset.BottomSide.Data
}

func (t *CubeTexture) SetBottomData(data []byte) {
	t.texAsset.BottomSide.Data = data
	t.target.SignalChange(CubeTextureDataChange)
}

func (t *CubeTexture) PreviewImage() image.Image {
	return t.previewImg
}

func (t *CubeTexture) SetPreviewImage(img image.Image) {
	t.previewImg = img
	t.target.SignalChange(CubeTexturePreviewChange)
}

func (t *CubeTexture) Save() error {
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
