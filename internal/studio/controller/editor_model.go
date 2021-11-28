package controller

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	gameasset "github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewModelEditor(studio *Studio, resource *gameasset.Resource) *ModelEditor {
	return &ModelEditor{
		BaseEditor: NewBaseEditor(),

		studio:   studio,
		resource: resource,

		propsAssetExpanded: true,
	}
}

var _ model.Editor = (*ModelEditor)(nil)

type ModelEditor struct {
	BaseEditor

	studio   *Studio
	resource *gameasset.Resource

	propsAssetExpanded bool
}

func (e *ModelEditor) ID() string {
	return e.resource.GUID
}

func (e *ModelEditor) Name() string {
	return e.resource.Name
}

func (e *ModelEditor) Icon() ui.Image {
	return co.OpenImage("resources/icons/model.png")
}

func (e *ModelEditor) CanSave() bool {
	return false
}

func (e *ModelEditor) Save() error {
	return nil
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

		co.WithChild("asset", co.New(view.AssetAccordion, func() {
			co.WithData(view.AssetAccordionData{
				AssetID:   editor.ID(),
				AssetName: editor.Name(),
				AssetType: "3D Model",
				Expanded:  editor.IsPropsAssetExpanded(),
			})
			co.WithCallbackData(view.AssetAccordionCallbackData{
				OnToggleExpanded: func() {
					editor.SetPropsAssetExpanded(!editor.IsPropsAssetExpanded())
				},
			})
		}))
	})
}))
