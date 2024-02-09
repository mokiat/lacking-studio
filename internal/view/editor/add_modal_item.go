package editor

import (
	"github.com/mokiat/gog/opt"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type NodeKindItem struct {
	Title string
	Kind  editormodel.NodeKind
}

var AddNodeModalItem = co.Define(&addModalItemComponent{})

type AddNodeModalItemData struct {
	Item     NodeKindItem
	Selected bool
}

type AddNodeModalItemCallbackData struct {
	OnSelected func(item NodeKindItem)
}

type addModalItemComponent struct {
	co.BaseComponent

	item     NodeKindItem
	selected bool

	onSelected func(item NodeKindItem)
}

func (c *addModalItemComponent) OnUpsert() {
	data := co.GetData[AddNodeModalItemData](c.Properties())
	c.item = data.Item
	c.selected = data.Selected

	callbackData := co.GetCallbackData[AddNodeModalItemCallbackData](c.Properties())
	c.onSelected = callbackData.OnSelected
}

func (c *addModalItemComponent) Render() co.Instance {
	return co.New(std.ListItem, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ListItemData{
			Selected: c.selected,
		})
		co.WithCallbackData(std.ListItemCallbackData{
			OnSelected: c.handleSelected,
		})

		co.WithChild("title", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
				FontSize:  opt.V(float32(24)),
				FontColor: opt.V(ui.Black()),
				Text:      c.item.Title,
			})
		}))
	})
}

func (c *addModalItemComponent) handleSelected() {
	c.onSelected(c.item)
}
