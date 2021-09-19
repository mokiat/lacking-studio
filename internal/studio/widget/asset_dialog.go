package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

type AssetDialogCallbackData struct {
	OnAssetSelected func( /*TODO*/ )
	OnClose         func()
}

var AssetDialog = co.Define(func(props co.Properties) co.Instance {
	var callbackData AssetDialogCallbackData
	props.InjectOptionalCallbackData(&callbackData, AssetDialogCallbackData{})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(ui.RGBA(0x00, 0x00, 0x00, 0xF0)),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("content", co.New(Paper, func() {
			co.WithData(PaperData{
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Width:            optional.NewInt(600),
				Height:           optional.NewInt(600),
				HorizontalCenter: optional.NewInt(0),
				VerticalCenter:   optional.NewInt(0),
			})

			co.WithChild("header", co.New(Toolbar, func() {
				co.WithLayoutData(mat.LayoutData{
					Top:   optional.NewInt(0),
					Left:  optional.NewInt(0),
					Right: optional.NewInt(0),
				})

				co.WithChild("tex_2d", co.New(ToolbarButton, func() {
					co.WithData(ToolbarButtonData{
						Text:     "Tex2D",
						Selected: true,
					})
				}))
				co.WithChild("tex_3d", co.New(ToolbarButton, func() {
					co.WithData(ToolbarButtonData{
						Text:     "Tex3D",
						Selected: false,
					})
				}))
				co.WithChild("model", co.New(ToolbarButton, func() {
					co.WithData(ToolbarButtonData{
						Text:     "Model",
						Selected: false,
					})
				}))
				co.WithChild("scene", co.New(ToolbarButton, func() {
					co.WithData(ToolbarButtonData{
						Text:     "Scene",
						Selected: false,
					})
				}))
			}))

			co.WithChild("footer", co.New(Toolbar, func() {
				co.WithData(ToolbarData{
					Flipped: true,
				})
				co.WithLayoutData(mat.LayoutData{
					Left:   optional.NewInt(0),
					Right:  optional.NewInt(0),
					Bottom: optional.NewInt(0),
					Height: optional.NewInt(50),
				})

				co.WithChild("open", co.New(ToolbarButton, func() {
					co.WithData(ToolbarButtonData{
						Text:     "Open",
						Disabled: true,
					})
				}))

				co.WithChild("cancel", co.New(ToolbarButton, func() {
					co.WithData(ToolbarButtonData{
						Text: "Cancel",
					})
					co.WithCallbackData(ToolbarButtonCallbackData{
						ClickListener: func() {
							callbackData.OnClose()
						},
					})
				}))
			}))
		}))
	})
})
