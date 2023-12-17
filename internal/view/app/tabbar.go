package app

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Tabbar = co.Define(&tabbarComponent{})

type TabbarData struct {
}

type tabbarComponent struct {
	co.BaseComponent
}

func (c *tabbarComponent) Render() co.Instance {
	return co.New(std.Tabbar, func() {
		co.WithLayoutData(c.Properties().LayoutData())

		// TODO: Children
	})
}
