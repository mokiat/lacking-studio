package model

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

const (
	PreviewSize = 64
)

var (
	ChangeResource        = mvc.NewChange("resource")
	ChangeResourceName    = mvc.SubChange(ChangeResource, "name")
	ChangeResourcePreview = mvc.SubChange(ChangeTwoDTexture, "preview")
)

const (
	ResourceKindTwoDTexture ResourceKind = "twod_texture"
	ResourceKindCubeTexture ResourceKind = "cube_texture"
	ResourceKindModel       ResourceKind = "model"
	ResourceKindScene       ResourceKind = "scene"
)

type ResourceKind = string

func openResource(registry *Registry, resource asset.Resource) (*Resource, error) {
	result := newResource(registry, resource)
	previewImg, err := resource.ReadPreview()
	if err != nil {
		if !errors.Is(err, asset.ErrNotFound) {
			return nil, fmt.Errorf("error writing preview: %w", err)
		}
		previewImg = nil
	}
	result.SetPreviewImage(previewImg)
	return result, nil
}

func newResource(registry *Registry, resource asset.Resource) *Resource {
	return &Resource{
		Observable: mvc.NewObservable(),
		registry:   registry,
		resource:   resource,
	}
}

type Resource struct {
	mvc.Observable
	registry   *Registry
	resource   asset.Resource
	previewImg image.Image
}

func (r *Resource) ID() string {
	return r.resource.ID()
}

func (r *Resource) Kind() ResourceKind {
	return r.resource.Kind()
}

func (r *Resource) Name() string {
	return r.resource.Name()
}

func (r *Resource) SetName(name string) {
	r.resource.SetName(name)
	r.SignalChange(ChangeResourceName)
}

func (r *Resource) PreviewImage() image.Image {
	return r.previewImg
}

func (r *Resource) SetPreviewImage(img image.Image) {
	r.previewImg = img
	r.SignalChange(ChangeResourcePreview)
}

func (r *Resource) Raw() asset.Resource {
	return r.resource
}

func (r *Resource) Save() error {
	if err := r.registry.Save(); err != nil {
		return fmt.Errorf("error saving resources: %w", err)
	}
	if r.previewImg != nil {
		if err := r.resource.WritePreview(r.previewImg); err != nil {
			return fmt.Errorf("error writing preview: %w", err)
		}
	}
	return nil
}

func (r *Resource) SaveContent(content asset.Encodable) error {
	return r.resource.WriteContent(content)
}

func (r *Resource) LoadContent(content asset.Decodable) error {
	return r.resource.ReadContent(content)
}

func (r *Resource) Clone() (*Resource, error) {
	newResource := r.registry.CreateResource(r.Kind(), fmt.Sprintf("%s copy", r.Name()))
	newResource.SetPreviewImage(r.PreviewImage())
	if err := newResource.Save(); err != nil {
		return nil, fmt.Errorf("error saving resource: %w", err)
	}
	var blob rawContent
	if err := r.LoadContent(&blob); err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}
	if err := newResource.SaveContent(&blob); err != nil {
		return nil, fmt.Errorf("error writing content: %w", err)
	}
	return newResource, nil
}

func (r *Resource) Delete() error {
	if err := r.resource.DeleteContent(); err != nil {
		return fmt.Errorf("error deleting content: %w", err)
	}
	if err := r.resource.DeletePreview(); err != nil {
		return fmt.Errorf("error deleting preview: %w", err)
	}
	r.registry.RemoveResource(r)
	if err := r.registry.Save(); err != nil {
		return fmt.Errorf("error saving resources: %w", err)
	}
	return nil
}

type rawContent struct {
	data bytes.Buffer
}

func (b *rawContent) DecodeFrom(in io.Reader) error {
	_, err := io.Copy(&b.data, in)
	return err
}

func (b *rawContent) EncodeTo(out io.Writer) error {
	_, err := io.Copy(out, &b.data)
	return err
}
