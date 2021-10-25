package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewModelEditor() *ModelEditor {
	return &ModelEditor{
		BaseEditor:         NewBaseEditor(),
		propsAssetExpanded: true,
	}
}

var _ Editor = (*ModelEditor)(nil)

type ModelEditor struct {
	BaseEditor
	propsAssetExpanded bool
}

func (e *ModelEditor) ID() string {
	return "2a4ddd33-b284-4d60-91eb-805f8b21a1d1"
}

func (e *ModelEditor) Name() string {
	return "Vehicle"
}

func (e *ModelEditor) Icon() ui.Image {
	return co.OpenImage("resources/icons/model.png")
}

func (e *ModelEditor) Update() {

}

func (e *ModelEditor) OnViewportMouseEvent(event widget.ViewportMouseEvent) {

}

func (e *ModelEditor) Scene() graphics.Scene {
	return nil
}

func (e *ModelEditor) Camera() graphics.Camera {
	return nil
}

func (e *ModelEditor) IsPropsAssetExpanded() bool {
	return e.propsAssetExpanded
}

func (e *ModelEditor) SetPropsAssetExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.NotifyChanged()
}

func (e *ModelEditor) RenderProperties() co.Instance {
	return co.New(ModelPropertiesView, func() {
		co.WithData(e)
	})
}

func (e *ModelEditor) Destroy() {

}

var ModelPropertiesView = co.Controlled(co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(*ModelEditor)

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
				AssetType: "3D Model",
				Expanded:  editor.IsPropsAssetExpanded(),
			})
			co.WithCallbackData(AssetAccordionCallbackData{
				OnToggleExpanded: func() {
					editor.SetPropsAssetExpanded(!editor.IsPropsAssetExpanded())
				},
			})
		}))
	})
}))
