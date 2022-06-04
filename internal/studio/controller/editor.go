package controller

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type Editor interface {
	Save() error
	Render(scope co.Scope, layoutData mat.LayoutData) co.Instance
	Destroy()
}
