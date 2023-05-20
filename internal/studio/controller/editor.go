package controller

import (
	co "github.com/mokiat/lacking/ui/component"
)

type Editor interface {
	Save() error
	Render(scope co.Scope, layoutData any) co.Instance
	Destroy()
}
