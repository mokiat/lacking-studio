package view

import (
	"strings"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/preview/model"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Registry = mvc.EventListener(co.Define(&registryComponent{}))

type RegistryData struct {
	AppModel *model.AppModel
}

type registryComponent struct {
	co.BaseComponent

	appModel *model.AppModel

	searchText string
}

func (c *registryComponent) OnUpsert() {
	data := co.GetData[RegistryData](c.Properties())
	c.appModel = data.AppModel
}

func (c *registryComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout: layout.Vertical(layout.VerticalSettings{
				ContentAlignment: layout.HorizontalAlignmentCenter,
				ContentSpacing:   20,
			}),
			Padding: ui.SymmetricSpacing(100, 20),
		})

		co.WithChild("search", co.New(std.EditBox, func() {
			co.WithLayoutData(layout.Data{
				Width: opt.V(600),
			})
			co.WithData(std.EditBoxData{
				Text: c.searchText,
			})
			co.WithCallbackData(std.EditBoxCallbackData{
				OnChange: c.handleSearchChange,
				OnReject: c.handleSearchCancel,
			})
		}))

		co.WithChild("separator", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				Width:  opt.V(600),
				Height: opt.V(1),
			})
			co.WithData(std.ContainerData{
				BackgroundColor: opt.V(std.OutlineColor),
			})
		}))

		co.WithChild("list", co.New(RegistryList, func() {
			co.WithLayoutData(layout.Data{
				Width: opt.V(600),
			})

			c.eachResource(func(resource *asset.Resource) {
				co.WithChild(resource.ID(), co.New(RegistryItem, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(RegistryItemData{
						Resource: resource,
					})
					co.WithCallbackData(RegistryItemCallbackData{
						OnSelected: c.handleResourceSelected,
					})
				}))
			})
		}))
	})
}

func (c *registryComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.RefreshEvent:
		c.Invalidate()
	}
}

func (c *registryComponent) handleSearchChange(text string) {
	c.searchText = text
	c.Invalidate()
}

func (c *registryComponent) handleSearchCancel() {
	c.searchText = ""
	c.Invalidate()
}

func (c *registryComponent) handleResourceSelected(resource *asset.Resource) {
	c.appModel.SetSelectedResource(resource)
}

func (c *registryComponent) eachResource(callback func(resource *asset.Resource)) {
	resources := c.appModel.Resources()
	for _, resource := range resources {
		if c.showResource(resource) {
			callback(resource)
		}
	}
}

func (c *registryComponent) showResource(resource *asset.Resource) bool {
	if c.searchText == "" {
		return true
	}
	return strings.Contains(resource.ID(), c.searchText) || strings.Contains(resource.Name(), c.searchText)
}
