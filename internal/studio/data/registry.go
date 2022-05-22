package data

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"

	"github.com/mokiat/lacking/game/asset"
)

const (
	PreviewSize = 64
)

func NewRegistry(delegate asset.Registry) *Registry {
	return &Registry{
		delegate: delegate,
	}
}

type Registry struct {
	delegate        asset.Registry
	resources       []*Resource
	resourcesFromID map[string]*Resource
}

func (r *Registry) Init() error {
	assetResources, err := r.delegate.ReadResources()
	if err != nil {
		return fmt.Errorf("error reading resources: %w", err)
	}

	r.resources = make([]*Resource, len(assetResources))
	r.resourcesFromID = make(map[string]*Resource)
	for i, assetResource := range assetResources {
		r.resources[i] = &Resource{
			registry: r,
			id:       assetResource.GUID,
			kind:     ResourceKind(assetResource.Kind),
			name:     assetResource.Name,
		}
		r.resourcesFromID[assetResource.GUID] = r.resources[i]
	}

	assetDependencies, err := r.delegate.ReadDependencies()
	if err != nil {
		return fmt.Errorf("error reading dependencies: %w", err)
	}

	for _, assetDependency := range assetDependencies {
		sourceResource := r.resourcesFromID[assetDependency.SourceGUID]
		if sourceResource == nil {
			return fmt.Errorf("cannot find resource with id %q", assetDependency.SourceGUID)
		}
		targetResource := r.resourcesFromID[assetDependency.TargetGUID]
		if targetResource == nil {
			return fmt.Errorf("cannot find resource with id %q", assetDependency.TargetGUID)
		}
		sourceResource.dependencies = append(sourceResource.dependencies, targetResource)
	}

	return nil
}

func (r *Registry) GetResourceByID(id string) *Resource {
	return r.resourcesFromID[id]
}

func (r *Registry) EachResource(filter Filter[*Resource], fn func(*Resource)) {
	for _, resource := range r.resources {
		if filter(resource) {
			fn(resource)
		}
	}
}

func (r *Registry) PreparePreview(img image.Image) image.Image {
	bounds := img.Bounds()

	var scaleFactor float64
	switch {
	case bounds.Dx() > PreviewSize && bounds.Dy() > PreviewSize:
		if bounds.Dx() > bounds.Dy() {
			scaleFactor = float64(PreviewSize) / float64(bounds.Dx())
		} else {
			scaleFactor = float64(PreviewSize) / float64(bounds.Dy())
		}
	case bounds.Dx() < PreviewSize && bounds.Dy() < PreviewSize:
		if bounds.Dx() > bounds.Dy() {
			scaleFactor = float64(PreviewSize) / float64(bounds.Dx())
		} else {
			scaleFactor = float64(PreviewSize) / float64(bounds.Dy())
		}
	case bounds.Dx() > PreviewSize:
		scaleFactor = float64(PreviewSize) / float64(bounds.Dx())
	case bounds.Dy() > PreviewSize:
		scaleFactor = float64(PreviewSize) / float64(bounds.Dy())
	default:
		return img
	}

	dstRect := image.Rect(
		0,
		0,
		int(float64(bounds.Dx())*scaleFactor),
		int(float64(bounds.Dy())*scaleFactor),
	)
	dst := image.NewNRGBA(dstRect)
	draw.ApproxBiLinear.Scale(dst, dstRect, img, img.Bounds(), draw.Src, nil)
	return dst
}

func (r *Registry) saveResources() error {
	assetResources := make([]asset.Resource, len(r.resources))
	for i, resource := range r.resources {
		assetResources[i] = asset.Resource{
			GUID: resource.id,
			Kind: string(resource.kind),
			Name: resource.name,
		}
	}
	if err := r.delegate.WriteResources(assetResources); err != nil {
		return fmt.Errorf("error writing resources: %w", err)
	}
	return nil
}

func (r *Registry) saveDependencies() error {
	var assetDependencies []asset.Dependency
	for _, resource := range r.resources {
		for _, dependency := range resource.dependencies {
			assetDependencies = append(assetDependencies, asset.Dependency{
				SourceGUID: resource.id,
				TargetGUID: dependency.id,
			})
		}
	}
	if err := r.delegate.WriteDependencies(assetDependencies); err != nil {
		return fmt.Errorf("error writing dependencies: %w", err)
	}
	return nil
}

func (r *Registry) readPreview(id string) (image.Image, error) {
	return r.delegate.ReadPreview(id)
}

func (r *Registry) writePreview(id string, img image.Image) error {
	return r.delegate.WritePreview(id, r.PreparePreview(img))
}

func (r *Registry) readContent(guid string, target asset.Decodable) error {
	return r.delegate.ReadContent(guid, target)
}

func (r *Registry) writeContent(guid string, target asset.Encodable) error {
	return r.delegate.WriteContent(guid, target)
}
