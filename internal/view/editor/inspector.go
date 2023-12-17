package editor

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Inspector = co.Define(&inspectorComponent{})

type InspectorData struct{}

type inspectorComponent struct {
	co.BaseComponent
}

func (c *inspectorComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Green()),
		})
	})
}
