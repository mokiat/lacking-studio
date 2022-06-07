package view

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/optional"
)

type BinaryInfoPropertiesSectionData struct {
	Binary *model.Binary
}

var BinaryInfoPropertiesSection = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data   = co.GetData[BinaryInfoPropertiesSectionData](props)
		binary = data.Binary
	)

	mvc.UseBinding(binary, func(change mvc.Change) bool {
		return true // TODO
	})

	return co.New(mat.Element, func() {
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
		co.WithLayoutData(mat.LayoutData{
			GrowHorizontally: true,
		})

		co.WithChild("size-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      "Size:",
			})
		}))

		co.WithChild("size-value-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      fmt.Sprintf("%d bytes", binary.Size()),
			})
		}))

		co.WithChild("sha-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      "Digest:",
			})
		}))

		co.WithChild("sha-value-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
				FontSize:  optional.Value(float32(18)),
				FontColor: optional.Value(ui.Black()),
				Text:      binary.Digest(),
			})
		}))
	})
})
