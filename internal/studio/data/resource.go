package data

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/mokiat/lacking/game/asset"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const (
	ResourceKindTwoDTexture ResourceKind = "twod_texture"
	ResourceKindCubeTexture ResourceKind = "cube_texture"
	ResourceKindModel       ResourceKind = "model"
	ResourceKindScene       ResourceKind = "scene"
)

type ResourceKind string

func newResource(registry *Registry, kind ResourceKind) *Resource {
	return &Resource{
		id:   uuid.Must(uuid.NewRandom()).String(),
		kind: kind,
		name: "Unnamed",
	}
}

type Resource struct {
	registry *Registry

	id            string
	kind          ResourceKind
	name          string
	resourceDirty bool

	dependencies      []*Resource
	dependenciesDirty bool

	previewImage image.Image
}

func (r *Resource) ID() string {
	return r.id
}

func (r *Resource) Kind() ResourceKind {
	return r.kind
}

func (r *Resource) Name() string {
	return r.name
}

func (r *Resource) SetName(name string) {
	r.name = name
	r.resourceDirty = true
}

func (r *Resource) Save() error {
	if r.resourceDirty {
		if err := r.registry.saveResources(); err != nil {
			return fmt.Errorf("error saving resources: %w", err)
		}
		r.resourceDirty = false
	}
	if r.dependenciesDirty {
		if err := r.registry.saveDependencies(); err != nil {
			return fmt.Errorf("error saving dependencies: %w", err)
		}
		r.dependenciesDirty = false
	}
	return nil
}

func (r *Resource) LoadPreview() (image.Image, error) {
	if r.previewImage != nil {
		return r.previewImage, nil
	}
	img, err := r.registry.readPreview(r.id)
	if err != nil {
		return nil, fmt.Errorf("error reading preview: %w", err)
	}
	r.previewImage = img
	return img, nil
}

func (r *Resource) SavePreview(img image.Image) error {
	if img == r.previewImage {
		return nil
	}
	bounds := img.Bounds()
	if width := bounds.Dx(); width > PreviewSize {
		return fmt.Errorf("width (%d) is larger than maximum (%d)", width, PreviewSize)
	}
	if height := bounds.Dy(); height > PreviewSize {
		return fmt.Errorf("height (%d) is larger than maximum (%d)", height, PreviewSize)
	}
	if err := r.registry.writePreview(r.id, img); err != nil {
		return fmt.Errorf("error writing preview: %w", err)
	}
	r.previewImage = img
	return nil
}

func (r *Resource) LoadContent(target asset.Decodable) error {
	if err := r.registry.readContent(r.id, target); err != nil {
		return fmt.Errorf("error reading content: %w", err)
	}
	return nil
}

func (r *Resource) SaveContent(target asset.Encodable) error {
	if err := r.registry.writeContent(r.id, target); err != nil {
		return fmt.Errorf("error writing content: %w", err)
	}
	return nil
}

func (r *Resource) Clone() (*Resource, error) {
	newResource := &Resource{
		registry:          r.registry,
		id:                uuid.NewString(),
		kind:              r.kind,
		name:              fmt.Sprintf("%s copy", r.name),
		resourceDirty:     true,
		dependencies:      slices.Clone(r.dependencies),
		dependenciesDirty: true,
		previewImage:      r.previewImage,
	}
	r.registry.resources = append(r.registry.resources, newResource)
	r.registry.resourceFromID[newResource.id] = newResource
	if err := newResource.Save(); err != nil {
		return nil, fmt.Errorf("error saving resource: %w", err)
	}
	if r.previewImage != nil {
		if err := r.registry.writePreview(r.id, r.previewImage); err != nil {
			return nil, fmt.Errorf("error writing preview: %w", err)
		}
	}
	var tmp rawResource
	if err := r.registry.readContent(r.id, &tmp); err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}
	if err := r.registry.writeContent(newResource.id, &tmp); err != nil {
		return nil, fmt.Errorf("error writing content: %w", err)
	}
	return newResource, nil
}

func (r *Resource) Delete() error {
	if err := r.registry.delegate.DeleteContent(r.id); err != nil {
		if !errors.Is(err, asset.ErrNotFound) {
			return fmt.Errorf("error deleting content: %w", err)
		}
	}
	if err := r.registry.delegate.DeletePreview(r.id); err != nil {
		if !errors.Is(err, asset.ErrNotFound) {
			return fmt.Errorf("error deleting preview: %w", err)
		}
	}
	delete(r.registry.resourceFromID, r.id)
	index := slices.Index(r.registry.resources, r)
	r.registry.resources = slices.Delete(r.registry.resources, index, index+1)
	if err := r.registry.saveResources(); err != nil {
		return fmt.Errorf("error saving resources: %w", err)
	}
	if err := r.registry.saveDependencies(); err != nil {
		return fmt.Errorf("error saving dependencies: %w", err)
	}
	return nil
}

type rawResource struct {
	data bytes.Buffer
}

func (b *rawResource) DecodeFrom(in io.Reader) error {
	_, err := io.Copy(&b.data, in)
	return err
}

func (b *rawResource) EncodeTo(out io.Writer) error {
	_, err := io.Copy(out, &b.data)
	return err
}
