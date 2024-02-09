package editor

import (
	"strings"

	"github.com/mokiat/gog/opt"
	editormodel "github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var AddNodeModal = co.Define(&addNodeModalComponent{})

type AddNodeModalCallbackData struct {
	OnAdd func(kind editormodel.NodeKind)
}

type addNodeModalComponent struct {
	co.BaseComponent

	items        []NodeKindItem
	selectedItem opt.T[NodeKindItem]
	searchText   string

	onAdd func(kind editormodel.NodeKind)
}

func (c *addNodeModalComponent) OnCreate() {
	c.items = []NodeKindItem{
		{
			Title: "Node",
			Kind:  editormodel.NodeKindNode,
		},
		{
			Title: "Point Light",
			Kind:  editormodel.NodeKindPointLight,
		},
	}
	c.selectedItem = opt.V(c.items[0])
}

func (c *addNodeModalComponent) OnUpsert() {
	callbackData := co.GetCallbackData[AddNodeModalCallbackData](c.Properties())
	c.onAdd = callbackData.OnAdd

	if c.selectedItem.Specified && !c.assetMatchesSearch(c.selectedItem.Value) {
		c.selectedItem = opt.Unspecified[NodeKindItem]()
	}
}

func (c *addNodeModalComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(600),
			Height:           opt.V(600),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("header", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})

			co.WithChild("search", co.New(std.EditBox, func() {
				co.WithData(std.EditBoxData{
					Text: c.searchText,
				})

				co.WithLayoutData(layout.Data{
					Width: opt.V(200),
				})

				co.WithCallbackData(std.EditBoxCallbackData{
					OnChange: c.handleSearchTextChange,
					OnReject: c.handleSearchReject,
				})
			}))
		}))

		co.WithChild("content", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ElementData{
				Padding: ui.SymmetricSpacing(0, 10),
				Layout:  layout.Fill(),
			})

			co.WithChild("scroll-pane", co.New(std.ScrollPane, func() {
				co.WithData(std.ScrollPaneData{
					DisableHorizontal: true,
					Focused:           false,
				})

				co.WithChild("list", co.New(std.List, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})

					c.eachItem(func(item NodeKindItem) {
						co.WithChild(string(item.Kind), co.New(AddNodeModalItem, func() {
							co.WithLayoutData(layout.Data{
								GrowHorizontally: true,
							})
							co.WithData(AddNodeModalItemData{
								Item:     item,
								Selected: c.selectedItem.Specified && (c.selectedItem.Value == item),
							})
							co.WithCallbackData(AddNodeModalItemCallbackData{
								OnSelected: c.handleItemSelected,
							})
						}))
					})
				}))
			}))

		}))

		co.WithChild("footer", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentBottom,
			})
			co.WithData(std.ToolbarData{
				Positioning: std.ToolbarPositioningBottom,
			})

			co.WithChild("open", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text:    "Open",
					Enabled: opt.V(c.selectedItem.Specified),
				})
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleAdd,
				})
			}))

			co.WithChild("cancel", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Cancel",
				})
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleCancel,
				})
			}))
		}))
	})
}

func (c *addNodeModalComponent) handleSearchTextChange(text string) {
	c.searchText = text
	c.Invalidate()
}

func (c *addNodeModalComponent) handleSearchReject() {
	c.searchText = ""
	c.Invalidate()
}

func (c *addNodeModalComponent) eachItem(cb func(item NodeKindItem)) {
	for _, item := range c.items {
		if c.assetMatchesSearch(item) {
			cb(item)
		}
	}
}

func (c *addNodeModalComponent) assetMatchesSearch(item NodeKindItem) bool {
	if c.searchText == "" {
		return true
	}
	return strings.Contains(item.Title, c.searchText)
}

func (c *addNodeModalComponent) handleItemSelected(item NodeKindItem) {
	c.selectedItem = opt.V(item)
	c.Invalidate()
}

func (c *addNodeModalComponent) handleAdd() {
	co.CloseOverlay(c.Scope())
	c.onAdd(c.selectedItem.Value.Kind)
}

func (c *addNodeModalComponent) handleCancel() {
	co.CloseOverlay(c.Scope())
}
