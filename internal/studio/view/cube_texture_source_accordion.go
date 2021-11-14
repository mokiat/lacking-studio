package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type CubeTextureSourceAccordionData struct {
	Expanded bool
	Filename string
	Image    ui.Image
}

type CubeTextureSourceAccordionCallbackData struct {
	OnToggle func()
	OnDrop   func(paths []string)
	OnReload func()
}

var CubeTextureSourceAccordion = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data CubeTextureSourceAccordionData
	props.InjectData(&data)

	var callbackData CubeTextureSourceAccordionCallbackData
	props.InjectCallbackData(&callbackData)

	return co.New(widget.Accordion, func() {
		co.WithLayoutData(mat.LayoutData{
			GrowHorizontally: true,
		})
		co.WithData(widget.AccordionData{
			Title:    "Source",
			Expanded: data.Expanded,
		})
		co.WithCallbackData(widget.AccordionCallbackData{
			OnToggle: callbackData.OnToggle,
		})

		co.WithChild("content", co.New(mat.Container, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(mat.ContainerData{
				Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
					ContentAlignment: mat.AlignmentCenter,
					ContentSpacing:   5,
				}),
				Padding: ui.Spacing{
					Left:   5,
					Right:  5,
					Top:    5,
					Bottom: 5,
				},
			})

			co.WithChild("label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "regular"),
					FontSize:  optional.NewInt(18),
					FontColor: optional.NewColor(ui.Black()),
					Text:      data.Filename,
				})
			}))

			co.WithChild("dropzone", co.New(widget.DropZone, func() {
				co.WithCallbackData(widget.DropZoneCallbackData{
					OnDrop: callbackData.OnDrop,
				})
				co.WithChild("image", co.New(mat.Picture, func() {
					co.WithData(mat.PictureData{
						BackgroundColor: optional.NewColor(ui.Gray()),
						Image:           data.Image,
						ImageColor:      optional.NewColor(ui.White()),
						Mode:            mat.ImageModeFit,
					})
					co.WithLayoutData(mat.LayoutData{
						Width:  optional.NewInt(200),
						Height: optional.NewInt(200),
					})
				}))
			}))

			// TODO: Add reload button
		}))
	})
}))
