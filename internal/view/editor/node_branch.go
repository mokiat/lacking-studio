package editor

import (
	"fmt"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/model/editor"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var (
	AccordionHeaderPadding        = 2
	AccordionHeaderContentSpacing = 5
	AccordionHeaderIconSize       = 24
	AccordionHeaderFontSize       = float32(20)
	AccordionHeaderFontFile       = "ui:///roboto-regular.ttf"
	AccordionExpandedIconFile     = "ui:///expanded.png"
	AccordionCollapsedIconFile    = "ui:///collapsed.png"
)

type NodeBranchData struct {
	Node         editor.Node
	SelectedNode editor.Node // TODO: Just pass the editor model maybe
}

type NodeBranchCallbackData struct {
	OnSelect func(node editor.Node)
}

var NodeBranch = co.Define(&nodeBranchComponent{})

type nodeBranchComponent struct {
	co.BaseComponent
	std.BaseButtonComponent

	node         editor.Node
	selectedNode editor.Node

	image *ui.Image
	title string

	isSelected bool
	isExpanded bool
	icon       *ui.Image

	onSelect func(node editor.Node)
}

func (c *nodeBranchComponent) OnUpsert() {
	data := co.GetData[NodeBranchData](c.Properties())
	c.node = data.Node
	c.selectedNode = data.SelectedNode

	c.title = c.node.Name()
	switch c.node.Kind() {
	case editor.NodeKindNode:
		c.image = co.OpenImage(c.Scope(), "icons/model.png")
	case editor.NodeKindPointLight:
		c.image = co.OpenImage(c.Scope(), "icons/point-light.png")
	default:
		c.image = nil
	}

	c.isSelected = c.node == data.SelectedNode
	if _, ok := c.node.(editor.ExtendableNode); !ok {
		c.isExpanded = false
	}
	if c.isExpanded {
		c.icon = co.OpenImage(c.Scope(), AccordionExpandedIconFile)
	} else {
		c.icon = co.OpenImage(c.Scope(), AccordionCollapsedIconFile)
	}
	c.SetOnClickFunc(c.onToggle)

	callbackData := co.GetCallbackData[NodeBranchCallbackData](c.Properties())
	c.onSelect = callbackData.OnSelect
}

func (c *nodeBranchComponent) Render() co.Instance {
	return co.New(co.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(co.ElementData{
			// Layout: layout.Frame(layout.FrameSettings{}),
			Layout: layout.Vertical(layout.VerticalSettings{
				ContentAlignment: layout.HorizontalAlignmentLeft,
			}),
		})

		co.WithChild("header", co.New(co.Element, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
				// VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(co.ElementData{
				Essence: c,
				Padding: ui.Spacing{
					Left:   AccordionHeaderPadding,
					Right:  AccordionHeaderPadding,
					Top:    AccordionHeaderPadding,
					Bottom: AccordionHeaderPadding,
				},
				Layout: layout.Horizontal(layout.HorizontalSettings{
					ContentAlignment: layout.VerticalAlignmentCenter,
					ContentSpacing:   AccordionHeaderContentSpacing,
				}),
			})

			co.WithChild("icon", co.New(std.Picture, func() {
				co.WithData(std.PictureData{
					Image:      c.icon,
					ImageColor: opt.V(std.OnPrimaryLightColor),
					Mode:       std.ImageModeFit,
				})
				co.WithLayoutData(layout.Data{
					Width:  opt.V(AccordionHeaderIconSize),
					Height: opt.V(AccordionHeaderIconSize),
				})
			}))

			if c.image != nil {
				co.WithChild("image", co.New(std.Picture, func() {
					co.WithData(std.PictureData{
						Image:      c.image,
						ImageColor: opt.V(std.OnPrimaryLightColor),
						Mode:       std.ImageModeFit,
					})
					co.WithLayoutData(layout.Data{
						Width:  opt.V(AccordionHeaderIconSize),
						Height: opt.V(AccordionHeaderIconSize),
					})
				}))
			}

			co.WithChild("title", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), AccordionHeaderFontFile),
					FontSize:  opt.V(AccordionHeaderFontSize),
					FontColor: opt.V(std.OnPrimaryLightColor),
					Text:      c.title,
				})
			}))
		}))

		co.WithChild("content", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				// VerticalAlignment: layout.VerticalAlignmentBottom,
				GrowHorizontally: true,
			})
			co.WithData(std.ContainerData{
				Padding: ui.Spacing{
					Left: 10,
				},
				Layout: layout.Vertical(layout.VerticalSettings{
					ContentAlignment: layout.HorizontalAlignmentLeft,
				}),
			})
			if c.isExpanded {
				if extendable, ok := c.node.(editor.ExtendableNode); ok {
					for i, childNode := range extendable.Children() {
						co.WithChild(fmt.Sprintf("child-%d", i), co.New(NodeBranch, func() {
							co.WithLayoutData(layout.Data{
								GrowHorizontally: true,
							})
							co.WithData(NodeBranchData{
								Node:         childNode,
								SelectedNode: c.selectedNode,
							})
							co.WithCallbackData(NodeBranchCallbackData{
								OnSelect: c.onSelect,
							})
						}))
					}
				}
			}
		}))
	})
}

func (c *nodeBranchComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var backgroundColor ui.Color
	switch c.State() {
	case std.ButtonStateOver:
		backgroundColor = std.PrimaryLightColor.Overlay(std.HoverOverlayColor)
	case std.ButtonStateDown:
		backgroundColor = std.PrimaryLightColor.Overlay(std.PressOverlayColor)
	default:
		if c.isSelected {
			backgroundColor = ui.Red() // FIXME
		} else {
			backgroundColor = std.PrimaryLightColor
		}
	}

	drawBounds := canvas.DrawBounds(element, false)

	canvas.Reset()
	canvas.SetStrokeSize(1.0)
	canvas.SetStrokeColor(std.OutlineColor)
	canvas.Rectangle(
		drawBounds.Position,
		drawBounds.Size,
	)
	if !backgroundColor.Transparent() {
		canvas.Fill(ui.Fill{
			Color: backgroundColor,
		})
	}
	canvas.Stroke()
}

func (c *nodeBranchComponent) onToggle() {
	c.isExpanded = !c.isExpanded
	c.Invalidate()
	c.onSelect(c.node)
}
