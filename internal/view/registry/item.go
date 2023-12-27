package registry

import (
	"github.com/mokiat/gog/opt"
	registrymodel "github.com/mokiat/lacking-studio/internal/model/registry"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Item = co.Define(&itemComponent{})

type ItemData struct {
	Asset    *registrymodel.Asset
	Selected bool
}

type ItemCallbackData struct {
	OnSelected func(asset *registrymodel.Asset)
}

type itemComponent struct {
	co.BaseComponent

	defaultPreviewImage *ui.Image

	asset    *registrymodel.Asset
	selected bool

	onSelected func(asset *registrymodel.Asset)
}

func (c *itemComponent) OnCreate() {
	c.defaultPreviewImage = co.OpenImage(c.Scope(), "icons/broken-image.png")
}

func (c *itemComponent) OnUpsert() {
	data := co.GetData[ItemData](c.Properties())
	c.asset = data.Asset
	c.selected = data.Selected

	callbackData := co.GetCallbackData[ItemCallbackData](c.Properties())
	c.onSelected = callbackData.OnSelected
}

func (c *itemComponent) Render() co.Instance {
	previewImage := c.asset.Image()
	if previewImage == nil {
		previewImage = c.defaultPreviewImage
	}

	return co.New(std.ListItem, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ListItemData{
			Selected: c.selected,
		})
		co.WithCallbackData(std.ListItemCallbackData{
			OnSelected: c.handleSelected,
		})

		co.WithChild("item", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("preview", co.New(std.Picture, func() {
				co.WithData(std.PictureData{
					Image:           previewImage,
					BackgroundColor: opt.V(ui.Black()),
					ImageColor:      opt.V(ui.White()),
					Mode:            std.ImageModeFit,
				})
				co.WithLayoutData(layout.Data{
					Width:  opt.V(64),
					Height: opt.V(64),
				})
			}))

			co.WithChild("info", co.New(std.Element, func() {
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentAlignment: layout.HorizontalAlignmentLeft,
						ContentSpacing:   5,
					}),
				})

				co.WithChild("name", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      c.asset.Name(),
					})
				}))

				co.WithChild("id", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(16)),
						FontColor: opt.V(ui.Black()),
						Text:      c.asset.ID(),
					})
				}))
			}))
		}))
	})
}

func (c *itemComponent) handleSelected() {
	c.onSelected(c.asset)
}
