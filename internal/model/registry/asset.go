package registry

import (
	asset "github.com/mokiat/lacking/game/newasset"
	"github.com/mokiat/lacking/ui"
)

type Asset struct {
	delegate     *asset.Resource
	previewImage *ui.Image
}

func (a *Asset) ID() string {
	return a.delegate.ID()
}

func (a *Asset) Name() string {
	return a.delegate.Name()
}

func (a *Asset) Image() *ui.Image {
	return a.previewImage
}
