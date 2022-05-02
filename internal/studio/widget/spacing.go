package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type SpacingData struct {
	Width  int
	Height int
}

var Spacing = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var data SpacingData
	props.InjectOptionalData(&data, SpacingData{})

	return co.New(mat.Element, func() {
		co.WithLayoutData(props.LayoutData())
		co.WithData(mat.ElementData{
			IdealSize: optional.Value(ui.Size{
				Width:  data.Width,
				Height: data.Height,
			}),
		})
	})
}))
