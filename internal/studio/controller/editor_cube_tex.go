package controller

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/studio/change"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
)

func NewCubeTextureEditor(studio *Studio, gfxEngine graphics.Engine) *CubeTextureEditor {
	gfxScene := gfxEngine.CreateScene()
	gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.0, 0.3, 1.0))

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetPosition(sprec.NewVec3(0.0, 0.0, 0.0))
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	gfxCamera.SetFoV(sprec.Degrees(66))
	gfxCamera.SetAutoExposure(true)
	gfxCamera.SetExposure(1.0)
	gfxCamera.SetAutoFocus(false)

	return &CubeTextureEditor{
		BaseEditor: NewBaseEditor(),

		studio: studio,

		propsAssetExpanded:  true,
		propsSourceExpanded: true,
		propsConfigExpanded: true,

		gfxEngine: gfxEngine,
		gfxScene:  gfxScene,
		gfxCamera: gfxCamera,

		sourceFilename: "---",
	}
}

var _ Editor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	BaseEditor

	studio *Studio

	propsAssetExpanded  bool
	propsSourceExpanded bool
	propsConfigExpanded bool

	gfxEngine      graphics.Engine
	gfxScene       graphics.Scene
	gfxCamera      graphics.Camera
	gfxCameraPitch sprec.Angle
	gfxCameraYaw   sprec.Angle
	gfxImage       graphics.CubeTexture

	sourceFilename string
	sourceImage    ui.Image

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
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
	// e.gfxCamera.SetRotation(
	// 	sprec.QuatProd(
	// 		sprec.RotationQuat(sprec.Degrees(1), sprec.UnitVec3(sprec.NewVec3(1.0, 2.0, 0.5))),
	// 		e.gfxCamera.Rotation(),
	// 	),
	// )
}

func (e *CubeTextureEditor) OnViewportMouseEvent(event widget.ViewportMouseEvent) {
	switch event.Type {
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = true
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
	case ui.MouseEventTypeMove:
		if e.rotatingCamera {
			e.gfxCameraPitch += sprec.Degrees(float32(event.Position.Y-e.oldMouseY) / 5)
			e.gfxCameraYaw += sprec.Degrees(float32(event.Position.X-e.oldMouseX) / 5)
			e.gfxCamera.SetRotation(sprec.QuatProd(
				sprec.RotationQuat(e.gfxCameraYaw, sprec.BasisYVec3()),
				sprec.RotationQuat(e.gfxCameraPitch, sprec.BasisXVec3()),
			))
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = false
		}
	}
}

func (e *CubeTextureEditor) Scene() graphics.Scene {
	return e.gfxScene
}

func (e *CubeTextureEditor) Camera() graphics.Camera {
	return e.gfxCamera
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

func (e *CubeTextureEditor) SourceFilename() string {
	return filepath.Base(e.sourceFilename)
}

func (e *CubeTextureEditor) SourceImage() ui.Image {
	return e.sourceImage
}

func (e *CubeTextureEditor) OnChangeSource(path string) {
	ch := &change.CubeTextureChangeSource{
		Controller: e,
		FromURI:    e.sourceFilename,
		ToURI:      path,
	}
	if err := e.changes.Push(ch); err != nil {
		panic(err)
	}
	e.studio.NotifyChanged()
}

func (e *CubeTextureEditor) ChangePreviewImage(img image.Image) {
	// TODO: Erase old image
	e.sourceImage = co.CreateImage(img)
	e.NotifyChanged()
}

func (e *CubeTextureEditor) ChangeGraphicsImage(definition graphics.CubeTextureDefinition) {
	if e.gfxImage != nil {
		e.gfxImage.Delete()
	}
	e.gfxImage = e.gfxEngine.CreateCubeTexture(definition)
	e.gfxScene.Sky().SetSkybox(e.gfxImage)
}

func (e *CubeTextureEditor) ChangeSourceFilename(uri string) {
	e.sourceFilename = uri
	e.NotifyChanged()
}

func (e *CubeTextureEditor) RenderProperties() co.Instance {
	return co.New(CubeTexturePropertiesView, func() {
		co.WithData(e)
	})
}

func (e *CubeTextureEditor) Destroy() {
	if e.gfxImage != nil {
		e.gfxImage.Delete()
	}
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

		co.WithChild("source", co.New(CubeTextureSourceAccordion, func() {
			co.WithLayoutData(mat.LayoutData{
				GrowHorizontally: true,
			})
			co.WithData(CubeTextureSourceAccordionData{
				Expanded: editor.IsPropsSourceExpanded(),
				Filename: editor.SourceFilename(),
				Image:    editor.SourceImage(),
			})
			co.WithCallbackData(CubeTextureSourceAccordionCallbackData{
				OnToggle: func() {
					editor.SetPropsSourceExpanded(!editor.IsPropsSourceExpanded())
				},
				OnDrop: func(paths []string) {
					editor.OnChangeSource(paths[0])
				},
				OnReload: func() {
					// TODO
					fmt.Println("RELOAD CUBE SOURCE")
				},
			})
		}))

		// co.WithChild("source", co.New(widget.Accordion, func() {
		// 	co.WithLayoutData(mat.LayoutData{
		// 		GrowHorizontally: true,
		// 	})
		// 	co.WithData(widget.AccordionData{
		// 		Title:    "Source",
		// 		Expanded: editor.IsPropsSourceExpanded(),
		// 	})
		// 	co.WithCallbackData(widget.AccordionCallbackData{
		// 		OnToggle: func() {
		// 			editor.SetPropsSourceExpanded(!editor.IsPropsSourceExpanded())
		// 		},
		// 	})

		// 	co.WithChild("content", co.New(mat.Label, func() {
		// 		co.WithData(mat.LabelData{
		// 			Font:      co.GetFont("roboto", "regular"),
		// 			FontSize:  optional.NewInt(20),
		// 			FontColor: optional.NewColor(ui.Black()),
		// 			Text:      "TODO: Source image here...",
		// 		})
		// 	}))
		// }))

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
