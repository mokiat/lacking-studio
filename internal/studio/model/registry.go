package model

import (
	"fmt"
	"strings"

	"github.com/mokiat/gog/filter"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewRegistry(delegate asset.Registry) (*Registry, error) {
	result := &Registry{
		Observable:  mvc.NewObservable(),
		delegate:    delegate,
		modelsCache: make(map[asset.Resource]*Resource),
	}
	for _, resource := range delegate.Resources() {
		resourceModel, err := openResource(result, resource)
		if err != nil {
			return nil, fmt.Errorf("error preparing resource %q: %w", resource.ID(), err)
		}
		result.modelsCache[resource] = resourceModel
	}
	return result, nil
}

type Registry struct {
	mvc.Observable
	delegate    asset.Registry
	modelsCache map[asset.Resource]*Resource
}

func (r *Registry) Save() error {
	return r.delegate.Save()
}

func (r *Registry) CreateResource(kind ResourceKind, name string) *Resource {
	resource := r.delegate.CreateResource(kind, name)
	resourceModel := newResource(r, resource)
	r.modelsCache[resource] = resourceModel
	return resourceModel
}

func (r *Registry) RemoveResource(resourceModel *Resource) {
	delete(r.modelsCache, resourceModel.Raw())
	resourceModel.Raw().Delete()
}

func (r *Registry) ResourceByID(id string) *Resource {
	resource := r.delegate.ResourceByID(id)
	if resource == nil {
		return nil
	}
	return r.modelsCache[resource]
}

func (r *Registry) IterateResources(cb func(*Resource), fltrs ...filter.Func[*Resource]) {
	fltr := filter.And(fltrs...)
	for _, resource := range r.delegate.Resources() {
		resourceModel := r.modelsCache[resource]
		if fltr(resourceModel) {
			cb(resourceModel)
		}
	}
}

func ResourcesWithKind(kind ResourceKind) filter.Func[*Resource] {
	return func(resource *Resource) bool {
		return resource.Kind() == kind
	}
}

func ResourcesWithSimilarName(name string) filter.Func[*Resource] {
	if name == "" {
		return filter.True[*Resource]()
	}
	name = strings.ToLower(name)
	return func(resource *Resource) bool {
		return strings.Contains(strings.ToLower(resource.Name()), name)
	}
}
