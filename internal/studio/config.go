package studio

import (
	"fmt"

	"github.com/mokiat/lacking/ui"
)

type Config struct{}

func (c Config) CreateView(window ui.Window) (ui.View, error) {
	template, err := window.OpenTemplate("resources/studio/editor/view.xml")
	if err != nil {
		return nil, fmt.Errorf("failed to open template: %w", err)
	}
	return window.CreateTemplatedView(template, &Handler{})
}
