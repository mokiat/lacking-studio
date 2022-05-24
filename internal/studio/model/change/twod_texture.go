package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
)

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
