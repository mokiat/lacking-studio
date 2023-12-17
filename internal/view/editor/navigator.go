package editor

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Navigator = co.Define(&navigatorComponent{})

type NavigatorData struct{}

type navigatorComponent struct {
	co.BaseComponent
}

func (c *navigatorComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Red()),
		})
	})
}
