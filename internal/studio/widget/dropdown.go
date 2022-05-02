package widget

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type DropdownData struct {
	Items       []DropdownItem
	SelectedKey interface{}
}

type DropdownItem struct {
	Key   interface{}
	Label string
}

type DropdownCallbackData struct {
	OnItemSelected func(key interface{})
}

var defaultDropdownCallbackData = DropdownCallbackData{
	OnItemSelected: func(key interface{}) {},
}

var Dropdown = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data DropdownData
	props.InjectOptionalData(&data, DropdownData{})

	var callbackData DropdownCallbackData
	props.InjectOptionalCallbackData(&callbackData, defaultDropdownCallbackData)

	itemsOverlay := co.UseState(func() interface{} {
		return co.Overlay{}
	})

	onOpen := func() {
		itemsOverlay.Set(co.OpenOverlay(co.New(DropdownItemList, func() {
			co.WithData(props.Data())
			co.WithCallbackData(DropdownListCallbackData{
				OnSelected: func(key interface{}) {
					overlay := itemsOverlay.Get().(co.Overlay)
					overlay.Close()

					callbackData.OnItemSelected(key)
				},
			})
		})))
	}

	essence := co.UseState(func() *dropdownEssence {
		return &dropdownEssence{}
	}).Get()
	essence.clickListener = onOpen

	label := ""
	for _, item := range data.Items {
		if item.Key == data.SelectedKey {
			label = item.Label
		}
	}

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Layout:  mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.GetFont("roboto", "bold"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      label,
			})
			co.WithLayoutData(mat.LayoutData{
				Left:           optional.Value(0),
				Right:          optional.Value(24),
				Height:         optional.Value(24),
				VerticalCenter: optional.Value(0),
			})
		}))

		co.WithChild("button", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				Image:      co.OpenImage("resources/icons/dropdown.png"),
				ImageColor: optional.Value(ui.Black()),
				Mode:       mat.ImageModeFit,
			})
			co.WithLayoutData(mat.LayoutData{
				Width:          optional.Value(24),
				Height:         optional.Value(24),
				Right:          optional.Value(0),
				VerticalCenter: optional.Value(0),
			})
		}))
	})
}))

type DropdownListCallbackData struct {
	OnSelected func(key interface{})
}

var DropdownItemList = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data DropdownData
	props.InjectOptionalData(&data, DropdownData{})

	var callbackData DropdownListCallbackData
	props.InjectCallbackData(&callbackData)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(ui.RGBA(0x00, 0x00, 0x00, 0xF0)),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("content", co.New(Paper, func() {
			co.WithData(PaperData{
				Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
					ContentSpacing: 5,
				}),
			})
			co.WithLayoutData(mat.LayoutData{
				HorizontalCenter: optional.Value(0),
				VerticalCenter:   optional.Value(0),
			})

			for i, item := range data.Items {
				func(i int, item DropdownItem) {
					co.WithChild(fmt.Sprintf("item-%d", i), co.New(ListItem, func() {
						co.WithData(ListItemData{
							Selected: item.Key == data.SelectedKey,
						})
						co.WithLayoutData(mat.LayoutData{
							GrowHorizontally: true,
						})
						co.WithCallbackData(ListItemCallbackData{
							OnSelected: func() {
								callbackData.OnSelected(item.Key)
							},
						})
						co.WithChild("label", co.New(mat.Label, func() {
							co.WithData(mat.LabelData{
								Font:      co.GetFont("roboto", "bold"),
								FontSize:  optional.Value(float32(18)),
								FontColor: optional.Value(ui.Black()),
								Text:      item.Label,
							})
						}))
					}))
				}(i, item)
			}
		}))
	})
}))

var _ ui.ElementRenderHandler = (*dropdownEssence)(nil)
var _ ui.ElementMouseHandler = (*dropdownEssence)(nil)

type dropdownEssence struct {
	state         buttonState
	clickListener mat.ClickListener
}

func (e *dropdownEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var outlineColor ui.Color
	switch e.state {
	case buttonStateOver:
		outlineColor = Gray
	case buttonStateDown:
		outlineColor = DarkGray
	default:
		outlineColor = LightGray
	}

	size := element.Bounds().Size
	canvas.Reset()
	canvas.Rectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
	)
	canvas.Fill(ui.Fill{
		Color: SurfaceColor,
	})

	canvas.Reset()
	canvas.SetStrokeSize(2.0)
	canvas.SetStrokeColor(outlineColor)
	canvas.RoundRectangle(
		sprec.ZeroVec2(),
		sprec.NewVec2(float32(size.Width), float32(size.Height)),
		sprec.NewVec4(5, 5, 5, 5),
	)
	canvas.Stroke()
}

func (e *dropdownEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	context := element.Context()
	switch event.Type {
	case ui.MouseEventTypeEnter:
		e.state = buttonStateOver
		context.Window().Invalidate()
	case ui.MouseEventTypeLeave:
		e.state = buttonStateUp
		context.Window().Invalidate()
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonLeft {
			if e.state == buttonStateDown {
				e.onClick()
			}
			e.state = buttonStateOver
			context.Window().Invalidate()
		}
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonLeft {
			e.state = buttonStateDown
			context.Window().Invalidate()
		}
	}
	return true
}

func (e *dropdownEssence) onClick() {
	if e.clickListener != nil {
		e.clickListener()
	}
}
