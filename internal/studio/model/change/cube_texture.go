package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
)

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

func CubeTextureContent(target *model.CubeTexture, from, to CubeTextureContentState) history.Change {
	return history.FuncChange(
		func() error {
			return target.Target().AccumulateChanges(func() error {
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
			return target.Target().AccumulateChanges(func() error {
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
