package view

import (
	"fmt"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

type BinaryInfoPropertiesSectionData struct {
	Binary *model.Binary
}

var BinaryInfoPropertiesSection = mvc.Wrap(co.Define(&binaryInfoPropertiesSectionComponent{}))

type binaryInfoPropertiesSectionComponent struct {
	co.BaseComponent

	binary *model.Binary
}

func (c *binaryInfoPropertiesSectionComponent) OnUpsert() {
	data := co.GetData[BinaryInfoPropertiesSectionData](c.Properties())
	c.binary = data.Binary

	mvc.UseBinding(c.Scope(), c.binary, func(change mvc.Change) bool {
		return true // TODO
	})
}

func (c *binaryInfoPropertiesSectionComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    5,
				Bottom: 5,
			},
			Layout: layout.Vertical(layout.VerticalSettings{
				ContentAlignment: layout.HorizontalAlignmentLeft,
				ContentSpacing:   5,
			}),
		})
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})

		co.WithChild("size-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Size:",
			})
		}))

		co.WithChild("size-value-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      fmt.Sprintf("%d bytes", c.binary.Size()),
			})
		}))

		co.WithChild("sha-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      "Digest:",
			})
		}))

		co.WithChild("sha-value-label", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
				FontSize:  opt.V(float32(18)),
				FontColor: opt.V(ui.Black()),
				Text:      c.binary.Digest(),
			})
		}))
	})
}
