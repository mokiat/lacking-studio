package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/optional"
)

type EditorController interface {
	OnRenameResource(name string)
}

type AssetPropertiesSectionData struct {
	Model            *model.Resource
	StudioController StudioController
	EditorController EditorController
}

var AssetPropertiesSection = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data             = co.GetData[AssetPropertiesSectionData](props)
		resource         = data.Model
		controller       = data.StudioController
		editorController = data.EditorController
	)

	mvc.UseBinding(resource, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeResourceName)
	})

	return co.New(mat.Element, func() {
		co.WithLayoutData(mat.LayoutData{
			GrowHorizontally: true,
		})
		co.WithData(mat.ElementData{
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    5,
				Bottom: 5,
			},
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
				ContentSpacing:   5,
			}),
		})

		// TODO: Use GridLayout (2 columns)

		co.WithChild("id", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
					ContentAlignment: mat.AlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("label", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(mat.OnSurfaceColor),
					Text:      "ID:",
				})
			}))

			co.WithChild("value", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(mat.OnSurfaceColor),
					Text:      resource.ID(),
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
					Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(mat.OnSurfaceColor),
					Text:      "Type:",
				})
			}))

			co.WithChild("value", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(mat.OnSurfaceColor),
					Text:      string(resource.Kind()),
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
					Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
					FontSize:  optional.Value(float32(18)),
					FontColor: optional.Value(mat.OnSurfaceColor),
					Text:      "Name:",
				})
			}))

			co.WithChild("value", co.New(mat.Editbox, func() {
				co.WithData(mat.EditboxData{
					Text: resource.Name(),
				})
				co.WithLayoutData(mat.LayoutData{
					Width: optional.Value(300),
				})
				co.WithCallbackData(mat.EditboxCallbackData{
					OnChanged: func(text string) {
						editorController.OnRenameResource(text)
					},
				})
			}))
		}))

		co.WithChild("actions", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
					ContentAlignment: mat.AlignmentCenter,
					ContentSpacing:   10,
				}),
			})

			co.WithChild("delete", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Icon: co.OpenImage(scope, "icons/delete.png"),
					Text: "Delete",
				})

				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: func() {
						controller.OnDeleteResource(resource.ID())
					},
				})
			}))

			co.WithChild("clone", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Icon: co.OpenImage(scope, "icons/file-copy.png"),
					Text: "Clone",
				})

				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: func() {
						newResource := controller.OnCloneResource(resource.ID())
						if newResource != nil {
							controller.OnOpenResource(newResource.ID())
						}
					},
				})
			}))
		}))
	})
})
