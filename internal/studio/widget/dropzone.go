package widget

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var defaultDropZoneCallbackData = DropZoneCallbackData{
	OnDrop: func(paths []string) {},
}

type DropZoneCallbackData struct {
	OnDrop func(paths []string)
}

var DropZone = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var callbackData DropZoneCallbackData
	props.InjectOptionalCallbackData(&callbackData, defaultDropZoneCallbackData)

	essence := co.UseState(func() *dropZoneEssence {
		return &dropZoneEssence{}
	}).Get()
	essence.onDrop = callbackData.OnDrop

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Layout:  mat.NewFillLayout(),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
}))

var _ ui.ElementMouseHandler = (*dropZoneEssence)(nil)

type dropZoneEssence struct {
	onDrop func(paths []string)
}

func (e *dropZoneEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	if event.Type != ui.MouseEventTypeDrop {
		return false
	}

	filePayload, ok := event.Payload.(ui.FilepathPayload)
	if !ok {
		return false
	}

	e.onDrop(filePayload.Paths)
	return true
}
