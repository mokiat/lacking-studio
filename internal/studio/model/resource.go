package model

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking/game/asset"
)

var (
	// NOTE: This allows for example the studio to subscribe for name changes only
	// for an editor.
	NameChange = observer.StringChange("name")
)

type Resource interface {
	ID() string
	Name() string
	SetName(name string)
	Kind() data.ResourceKind
}

type Wrappable interface {
	Wrapping() asset.WrapMode
	SetWrapping(asset.WrapMode)
}

type Filterable interface {
	Filtering() asset.FilterMode
	SetFiltering(asset.FilterMode)
}

type Mipmappable interface {
	Mipmapping() bool
	SetMipmapping(bool)
}

type GammaCorrectable interface {
	GammaCorrection() bool
	SetGammaCorrection(bool)
}
