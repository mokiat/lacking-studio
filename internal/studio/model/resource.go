package model

import (
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeResource     = mvc.NewChange("resource")
	ChangeResourceName = mvc.SubChange(ChangeResource, "name")
)

func NewResource(resource *data.Resource) *Resource {
	return &Resource{
		Observable: mvc.NewObservable(),
		resource:   resource,
	}
}

type Resource struct {
	mvc.Observable
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

func (r *Resource) Raw() *data.Resource {
	return r.resource
}
