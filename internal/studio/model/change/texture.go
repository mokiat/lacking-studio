package change

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking/game/asset"
)

type Wrappable interface {
	SetWrapping(asset.WrapMode)
}

type WrappingState struct {
	Value asset.WrapMode
}

func Wrapping(target Wrappable, from, to WrappingState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetWrapping(to.Value)
			return nil
		},
		func() error {
			target.SetWrapping(from.Value)
			return nil
		},
	)
}

type Filterable interface {
	SetFiltering(asset.FilterMode)
}

type FilteringState struct {
	Value asset.FilterMode
}

func Filtering(target Filterable, from, to FilteringState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetFiltering(to.Value)
			return nil
		},
		func() error {
			target.SetFiltering(from.Value)
			return nil
		},
	)
}

type Mipmappable interface {
	SetMipmapping(bool)
}

type GammaCorrectable interface {
	SetGammaCorrection(bool)
}

type TwoDTexture interface {
	observer.Target
	SetWidth(int)
	SetHeight(int)
	SetFormat(asset.TexelFormat)
	SetData([]byte)
}

type TwoDTextureContentState struct {
	Width  int
	Height int
	Format asset.TexelFormat
	Data   []byte
}

func TwoDTextureContent(target TwoDTexture, from, to TwoDTextureContentState) history.Change {
	return history.FuncChange(
		func() error {
			return target.AccumulateChanges(func() error {
				target.SetWidth(to.Width)
				target.SetHeight(to.Height)
				target.SetFormat(to.Format)
				target.SetData(to.Data)
				return nil
			})
		},
		func() error {
			return target.AccumulateChanges(func() error {
				target.SetWidth(from.Width)
				target.SetHeight(from.Height)
				target.SetFormat(from.Format)
				target.SetData(from.Data)
				return nil
			})
		},
	)
}

type CubeTexture interface {
	observer.Target
	SetDimension(int)
	SetFormat(asset.TexelFormat)
	SetFrontData([]byte)
	SetBackData([]byte)
	SetLeftData([]byte)
	SetRightData([]byte)
	SetTopData([]byte)
	SetBottomData([]byte)
}

type CubeTextureContentState struct {
	Dimension  int
	Format     asset.TexelFormat
	FrontData  []byte
	BackData   []byte
	LeftData   []byte
	RightData  []byte
	TopData    []byte
	BottomData []byte
}

func CubeTextureContent(target CubeTexture, from, to CubeTextureContentState) history.Change {
	return history.FuncChange(
		func() error {
			return target.AccumulateChanges(func() error {
				target.SetDimension(to.Dimension)
				target.SetFormat(to.Format)
				target.SetFrontData(to.FrontData)
				target.SetBackData(to.BackData)
				target.SetLeftData(to.LeftData)
				target.SetRightData(to.RightData)
				target.SetTopData(to.TopData)
				target.SetBottomData(to.BottomData)
				return nil
			})
		},
		func() error {
			return target.AccumulateChanges(func() error {
				target.SetDimension(from.Dimension)
				target.SetFormat(from.Format)
				target.SetFrontData(from.FrontData)
				target.SetBackData(from.BackData)
				target.SetLeftData(from.LeftData)
				target.SetRightData(from.RightData)
				target.SetTopData(from.TopData)
				target.SetBottomData(from.BottomData)
				return nil
			})
		},
	)
}
