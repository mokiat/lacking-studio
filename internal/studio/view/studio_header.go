package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

type StudioHeaderData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var StudioHeader = co.Define(&studioHeaderComponent{})

type studioHeaderComponent struct {
	co.BaseComponent

	studio     *model.Studio
	controller StudioController
}

func (c *studioHeaderComponent) OnUpsert() {
	data := co.GetData[StudioHeaderData](c.Properties())
	c.studio = data.StudioModel
	c.controller = data.StudioController
}

func (c *studioHeaderComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Layout: layout.Vertical(layout.VerticalSettings{
				ContentAlignment: layout.HorizontalAlignmentLeft,
			}),
		})
		co.WithLayoutData(c.Properties().LayoutData())

		co.WithChild("toolbar", co.New(StudioToolbar, func() {
			co.WithData(StudioToolbarData{
				StudioModel:      c.studio,
				StudioController: c.controller,
			})
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
		}))

		co.WithChild("tabbar", co.New(StudioTabbar, func() {
			co.WithData(StudioTabbarData{
				StudioModel:      c.studio,
				StudioController: c.controller,
			})
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
		}))
	})
}
