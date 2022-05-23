package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
)

type TwoDTextureNameState struct {
	Value string
}

func TwoDTextureName(target *model.TwoDTexture, from, to TwoDTextureNameState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetName(to.Value)
			return nil
		},
		func() error {
			target.SetName(from.Value)
			return nil
		},
	)
}

type TwoDTextureWrappingState struct {
	Value asset.WrapMode
}

func TwoDTextureWrapping(target *model.TwoDTexture, from, to TwoDTextureWrappingState) history.Change {
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

type TwoDTextureFilteringState struct {
	Value asset.FilterMode
}

func TwoDTextureFiltering(target *model.TwoDTexture, from, to TwoDTextureFilteringState) history.Change {
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

type TwoDTextureContentState struct {
	Width  int
	Height int
	Format asset.TexelFormat
	Data   []byte
}

func TwoDTextureContent(target *model.TwoDTexture, from, to TwoDTextureContentState) history.Change {
	return history.FuncChange(
		func() error {
			return target.Target().AccumulateChanges(func() error {
				target.SetWidth(to.Width)
				target.SetHeight(to.Height)
				target.SetFormat(to.Format)
				target.SetData(to.Data)
				return nil
			})
		},
		func() error {
			return target.Target().AccumulateChanges(func() error {
				target.SetWidth(from.Width)
				target.SetHeight(from.Height)
				target.SetFormat(from.Format)
				target.SetData(from.Data)
				return nil
			})
		},
	)
}
