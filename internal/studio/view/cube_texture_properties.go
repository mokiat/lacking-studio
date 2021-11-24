package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

var CubeTextureProperties = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.CubeTextureEditor)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
				ContentSpacing:   5,
			}),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("asset", co.New(AssetAccordion, func() {
			co.WithData(AssetAccordionData{
				AssetID:   editor.ID(),
				AssetName: editor.Name(),
				AssetType: "Cube Texture",
				Expanded:  editor.IsAssetAccordionExpanded(),
			})
			co.WithCallbackData(AssetAccordionCallbackData{
				OnToggleExpanded: func() {
					editor.SetAssetAccordionExpanded(!editor.IsAssetAccordionExpanded())
				},
			})
		}))

		co.WithChild("source", co.New(CubeTextureSourceAccordion, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(CubeTextureSourceAccordionData{
				Expanded: editor.IsSourceAccordionExpanded(),
				Filename: editor.SourceFilename(),
				Image:    editor.SourcePreview(),
			})
			co.WithCallbackData(CubeTextureSourceAccordionCallbackData{
				OnToggle: func() {
					editor.SetSourceAccordionExpanded(!editor.IsSourceAccordionExpanded())
				},
				OnDrop: func(paths []string) {
					editor.ChangeSourcePath(paths[0])
				},
				OnReload: func() {
					editor.ReloadSource()
				},
			})
		}))

		co.WithChild("config", co.New(widget.Accordion, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(widget.AccordionData{
				Title:    "Config",
				Expanded: editor.IsConfigAccordionExpanded(),
			})
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
		}))
	})
}))
