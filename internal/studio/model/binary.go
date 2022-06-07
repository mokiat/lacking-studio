package model

import (
	"crypto/sha1"
	"fmt"

	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeBinary     = mvc.NewChange("binary")
	ChangeBinaryData = mvc.SubChange(ChangeBinary, "data")
)

func CreateBinary(registry *Registry) (*Binary, error) {
	binAsset := &asset.Binary{
		Data: []byte{},
	}
	resourceModel := registry.CreateResource(ResourceKindBinary, "Unnamed")
	if err := resourceModel.SaveContent(binAsset); err != nil {
		return nil, fmt.Errorf("error saving content: %w", err)
	}
	if err := resourceModel.Save(); err != nil {
		return nil, err
	}
	return &Binary{
		Observable:    mvc.NewObservable(),
		resourceModel: resourceModel,
		binAsset:      binAsset,
	}, nil
}

func OpenBinary(resourceModel *Resource) (*Binary, error) {
	binAsset := new(asset.Binary)
	if err := resourceModel.LoadContent(binAsset); err != nil {
		return nil, fmt.Errorf("error loading content: %w", err)
	}
	return &Binary{
		Observable:    mvc.NewObservable(),
		resourceModel: resourceModel,
		binAsset:      binAsset,
	}, nil
}

type Binary struct {
	mvc.Observable
	resourceModel *Resource
	binAsset      *asset.Binary
	digest        string
}

func (b *Binary) Resource() *Resource {
	return b.resourceModel
}

func (b *Binary) Data() []byte {
	return b.binAsset.Data
}

func (b *Binary) SetData(data []byte) {
	b.binAsset.Data = data
	b.digest = ""
	b.SignalChange(ChangeBinaryData)
}

func (b *Binary) Size() int {
	return len(b.Data())
}

func (b *Binary) Digest() string {
	if b.digest == "" {
		b.digest = fmt.Sprintf("%x", sha1.Sum(b.Data()))
	}
	return b.digest
}

func (b *Binary) Save() error {
	if err := b.resourceModel.SaveContent(b.binAsset); err != nil {
		return fmt.Errorf("error saving content: %w", err)
	}
	return b.resourceModel.Save()
}
