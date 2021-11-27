package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

var CubeTextureConfig = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.CubeTextureEditor)

	return co.New(widget.Accordion, func() {
		co.WithData(widget.AccordionData{
			Title:    "Config",
			Expanded: editor.IsConfigAccordionExpanded(),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithCallbackData(widget.AccordionCallbackData{
			OnToggle: func() {
				editor.SetConfigAccordionExpanded(!editor.IsConfigAccordionExpanded())
			},
		})

		co.WithChild("content", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.GetFont("roboto", "regular"),
				FontSize:  optional.NewInt(20),
				FontColor: optional.NewColor(ui.Black()),
				Text:      "TODO: Asset config here...",
			})
		}))
	})
}))
