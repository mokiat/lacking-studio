package change

import (
	"fmt"
	"image"
	"os"

	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/graphics"
)

type CubeTextureChangeSourceController interface {
	Alter(func())
	ChangePreviewImage(img image.Image)
	ChangeGraphicsImage(definition graphics.CubeTextureDefinition)
	ChangeSourceFilename(uri string)
}

var _ history.Change = (*CubeTextureChangeSource)(nil)

type CubeTextureChangeSource struct {
	Controller CubeTextureChangeSourceController
	FromURI    string
	ToURI      string
}

func (ch *CubeTextureChangeSource) Apply() error {
	return ch.run(ch.ToURI)
}

func (ch *CubeTextureChangeSource) Revert() error {
	return ch.run(ch.FromURI)
}

func (ch *CubeTextureChangeSource) run(uri string) error {
	img, err := ch.openImage(uri)
	if err != nil {
		return fmt.Errorf("failed to open source image: %w", err)
	}

	packImg := pack.BuildImageResource(img)
	frontPackImg := pack.BuildCubeSideFromEquirectangular(packImg, pack.CubeSideFront)
	rearPackImg := pack.BuildCubeSideFromEquirectangular(packImg, pack.CubeSideRear)
	leftPackImg := pack.BuildCubeSideFromEquirectangular(packImg, pack.CubeSideLeft)
	rightPackImg := pack.BuildCubeSideFromEquirectangular(packImg, pack.CubeSideRight)
	topPackImg := pack.BuildCubeSideFromEquirectangular(packImg, pack.CubeSideTop)
	bottomPackImg := pack.BuildCubeSideFromEquirectangular(packImg, pack.CubeSideBottom)
	cubeImg, err := pack.BuildCube(frontPackImg, rearPackImg, leftPackImg, rightPackImg, topPackImg, bottomPackImg, 0)
	if err != nil {
		return fmt.Errorf("failed to build cube image: %w", err)
	}

	ch.Controller.Alter(func() {
		ch.Controller.ChangePreviewImage(img)
		ch.Controller.ChangeGraphicsImage(graphics.CubeTextureDefinition{
			Dimension:      cubeImg.Dimension,
			WrapS:          graphics.WrapClampToEdge,
			WrapT:          graphics.WrapClampToEdge,
			MinFilter:      graphics.FilterNearest,
			MagFilter:      graphics.FilterNearest,
			InternalFormat: graphics.InternalFormatRGBA8,
			DataFormat:     graphics.DataFormatRGBA8,
			FrontSideData:  cubeImg.RGBA8Data(pack.CubeSideFront),
			BackSideData:   cubeImg.RGBA8Data(pack.CubeSideRear),
			LeftSideData:   cubeImg.RGBA8Data(pack.CubeSideLeft),
			RightSideData:  cubeImg.RGBA8Data(pack.CubeSideRight),
			TopSideData:    cubeImg.RGBA8Data(pack.CubeSideTop),
			BottomSideData: cubeImg.RGBA8Data(pack.CubeSideBottom),
		})
		ch.Controller.ChangeSourceFilename(uri)
	})
	return nil
}

func (ch *CubeTextureChangeSource) openImage(path string) (image.Image, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image resource: %w", err)
	}
	defer in.Close()

	img, _, err := image.Decode(in)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}
