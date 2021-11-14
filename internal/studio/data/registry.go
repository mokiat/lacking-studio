package data

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/mokiat/lacking/game/asset"
)

const (
	previewSize = 128
)

func NewRegistry(delegate asset.Registry) *Registry {
	return &Registry{
		delegate: delegate,
	}
}

type Registry struct {
	delegate asset.Registry
}

func (r *Registry) PreparePreview(img image.Image) image.Image {
	bounds := img.Bounds()

	var scaleFactor float64
	switch {
	case bounds.Dx() > previewSize && bounds.Dy() > previewSize:
		if bounds.Dx() > bounds.Dy() {
			scaleFactor = float64(previewSize) / float64(bounds.Dx())
		} else {
			scaleFactor = float64(previewSize) / float64(bounds.Dy())
		}
	case bounds.Dx() < previewSize && bounds.Dy() < previewSize:
		if bounds.Dx() > bounds.Dy() {
			scaleFactor = float64(previewSize) / float64(bounds.Dx())
		} else {
			scaleFactor = float64(previewSize) / float64(bounds.Dy())
		}
	case bounds.Dx() > previewSize:
		scaleFactor = float64(previewSize) / float64(bounds.Dx())
	case bounds.Dy() > previewSize:
		scaleFactor = float64(previewSize) / float64(bounds.Dy())
	default:
		return img
	}

	dstRect := image.Rect(
		0,
		0,
		int(float64(bounds.Dx())*scaleFactor),
		int(float64(bounds.Dy())*scaleFactor),
	)
	dst := image.NewNRGBA(dstRect)
	draw.ApproxBiLinear.Scale(dst, dstRect, img, img.Bounds(), draw.Src, nil)
	return dst
}

func (r *Registry) WritePreview(guid string, img image.Image) error {
	return r.delegate.WritePreview(guid, r.PreparePreview(img))
}

func (r *Registry) WriteContent(guid string, target interface{}) error {
	return r.delegate.WriteContent(guid, target)
}
