package data

import (
	"fmt"
	"image"

	"github.com/mokiat/lacking/game/asset"

	"github.com/google/uuid"
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
