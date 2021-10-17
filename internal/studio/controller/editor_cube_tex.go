package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

func NewCubeTextureEditor() *CubeTextureEditor {
	return &CubeTextureEditor{
		Controller: co.NewBaseController(),

		propsAssetExpanded:  true,
		propsSourceExpanded: true,
		propsConfigExpanded: true,
	}
}

var _ Editor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	co.Controller

	propsAssetExpanded  bool
	propsSourceExpanded bool
	propsConfigExpanded bool
}

func (e *CubeTextureEditor) ID() string {
	return "bab99e80-ded1-459a-b00b-6a17afa44046"
}

func (e *CubeTextureEditor) Name() string {
	return "Night-Sky"
}

func (e *CubeTextureEditor) Icon() ui.Image {
	return co.OpenImage("resources/icons/texture.png")
}

func (e *CubeTextureEditor) Update() {

}

func (e *CubeTextureEditor) IsPropsAssetExpanded() bool {
	return e.propsAssetExpanded
}

func (e *CubeTextureEditor) SetPropsAssetExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) IsPropsSourceExpanded() bool {
	return e.propsSourceExpanded
}

func (e *CubeTextureEditor) SetPropsSourceExpanded(expanded bool) {
	e.propsSourceExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) IsPropsConfigExpanded() bool {
	return e.propsConfigExpanded
}

func (e *CubeTextureEditor) SetPropsConfigExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) RenderProperties() co.Instance {
	return co.New(CubeTexturePropertiesView, func() {
		co.WithData(e)
	})
}

var CubeTexturePropertiesView = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(*CubeTextureEditor)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
				ContentSpacing:   5,
			}),
		})

		co.WithChild("asset", co.New(AssetAccordion, func() {
			co.WithData(AssetAccordionData{
				AssetID:   editor.ID(),
				AssetName: editor.Name(),
				AssetType: "Cube Texture",
				Expanded:  editor.IsPropsAssetExpanded(),
			})
			co.WithCallbackData(AssetAccordionCallbackData{
				OnToggleExpanded: func() {
					editor.SetPropsAssetExpanded(!editor.IsPropsAssetExpanded())
				},
			})
		}))

		co.WithChild("source", co.New(widget.Accordion, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(widget.AccordionData{
				Title:    "Source",
				Expanded: editor.IsPropsSourceExpanded(),
			})
			co.WithCallbackData(widget.AccordionCallbackData{
				OnToggle: func() {
					editor.SetPropsSourceExpanded(!editor.IsPropsSourceExpanded())
				},
			})

			co.WithChild("content", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.GetFont("roboto", "regular"),
					FontSize:  optional.NewInt(20),
					FontColor: optional.NewColor(ui.Black()),
					Text:      "TODO: Source image here...",
				})
			}))
		}))

		co.WithChild("config", co.New(widget.Accordion, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(widget.AccordionData{
				Title:    "Config",
				Expanded: editor.IsPropsConfigExpanded(),
			})
			co.WithCallbackData(widget.AccordionCallbackData{
				OnToggle: func() {
					editor.SetPropsConfigExpanded(!editor.IsPropsConfigExpanded())
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

type AssetAccordionData struct {
	AssetID   string
	AssetName string
	AssetType string

	Expanded bool
}

type AssetAccordionCallbackData struct {
	OnToggleExpanded func()
}

var AssetAccordion = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data AssetAccordionData
	props.InjectData(&data)

	var callbackData AssetAccordionCallbackData
	props.InjectCallbackData(&callbackData)

	return co.New(widget.Accordion, func() {
		co.WithLayoutData(mat.LayoutData{
			GrowHorizontally: true,
		})
		co.WithData(widget.AccordionData{
			Title:    "Asset",
			Expanded: data.Expanded,
		})
		co.WithCallbackData(widget.AccordionCallbackData{
			OnToggle: callbackData.OnToggleExpanded,
		})

		co.WithChild("content", co.New(mat.Container, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(mat.ContainerData{
				Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
					ContentAlignment: mat.AlignmentLeft,
					ContentSpacing:   5,
				}),
				Padding: ui.Spacing{
					Left:   5,
					Right:  5,
					Top:    5,
					Bottom: 5,
				},
			})

			co.WithChild("id", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   10,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "bold"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      "ID:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      data.AssetID,
					})
				}))
			}))

			co.WithChild("type", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   10,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "bold"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      "Type:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      data.AssetType,
					})
				}))
			}))

			co.WithChild("name", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
						ContentAlignment: mat.AlignmentCenter,
						ContentSpacing:   10,
					}),
				})

				co.WithChild("label", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "bold"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      "Name:",
					})
				}))

				co.WithChild("value", co.New(mat.Label, func() {
					co.WithData(mat.LabelData{
						Font:      co.GetFont("roboto", "regular"),
						FontSize:  optional.NewInt(18),
						FontColor: optional.NewColor(ui.Black()),
						Text:      data.AssetName,
					})
				}))
			}))
		}))
	})
}))
