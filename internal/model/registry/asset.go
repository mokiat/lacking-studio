package registry

import "github.com/mokiat/lacking/ui"

type Asset struct {
}

func (a *Asset) ID() string {
	return "guid-guid-guid"
}

func (a *Asset) Name() string {
	return "Funny Game Character"
}

func (a *Asset) Image() *ui.Image {
	return nil
}
