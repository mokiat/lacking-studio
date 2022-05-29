package model

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/data"
)

var (
	ChangeResource     = observer.NewChange("resource")
	ChangeResourceName = observer.ExtChange(ChangeResource, "name")
)

func NewResource(resource *data.Resource) *Resource {
	return &Resource{
		Target:   observer.NewTarget(),
		resource: resource,
	}
}

type Resource struct {
	observer.Target
	resource *data.Resource
}

func (r *Resource) ID() string {
	return r.resource.ID()
}

func (r *Resource) Name() string {
	return r.resource.Name()
}

func (r *Resource) SetName(name string) {
	r.resource.SetName(name)
	r.SignalChange(ChangeResourceName)
}

func (r *Resource) Kind() data.ResourceKind {
	return r.resource.Kind()
}
