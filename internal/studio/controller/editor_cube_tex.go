package controller

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/studio/change"
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/data/pack"
	gameasset "github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
)

func NewCubeTextureEditor(studio *Studio, resource *gameasset.Resource) *CubeTextureEditor {
	gfxScene := studio.GraphicsEngine().CreateScene()
	gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.0, 0.3, 1.0))

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetPosition(sprec.NewVec3(0.0, 0.0, 0.0))
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	gfxCamera.SetFoV(sprec.Degrees(66))
	gfxCamera.SetAutoExposure(true)
	gfxCamera.SetExposure(1.0)
	gfxCamera.SetAutoFocus(false)

	return &CubeTextureEditor{
		BaseEditor: NewBaseEditor(),

		studio:   studio,
		resource: resource,

		propsAssetExpanded:  true,
		propsSourceExpanded: true,
		propsConfigExpanded: true,

		gfxEngine: studio.GraphicsEngine(),
		gfxScene:  gfxScene,
		gfxCamera: gfxCamera,

		sourcePath: "---",
	}
}

var _ Editor = (*CubeTextureEditor)(nil)
var _ model.CubeTextureEditor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	BaseEditor

	studio   *Studio
	resource *gameasset.Resource

	propsAssetExpanded  bool
	propsSourceExpanded bool
	propsConfigExpanded bool

	gfxEngine      graphics.Engine
	gfxScene       graphics.Scene
	gfxCamera      graphics.Camera
	gfxCameraPitch sprec.Angle
	gfxCameraYaw   sprec.Angle
	gfxImage       graphics.CubeTexture

	definition  graphics.CubeTextureDefinition
	sourcePath  string
	sourceImg   image.Image
	sourceImage ui.Image

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int

	savedChange history.Change
}

func (e *CubeTextureEditor) ID() string {
	return e.resource.GUID
}

func (e *CubeTextureEditor) Name() string {
	return e.resource.Name
}

func (e *CubeTextureEditor) Icon() ui.Image {
	return co.OpenImage("resources/icons/texture.png")
}

func (e *CubeTextureEditor) CanSave() bool {
	return e.savedChange != e.changes.LastChange()
}

func (e *CubeTextureEditor) Save() {
	texOut := &asset.CubeTexture{
		Dimension: uint16(e.definition.Dimension),
		Format:    asset.DataFormatRGBA32F,
		Sides:     [6]asset.CubeTextureSide{
			// TODO
		},
	}
	if err := e.studio.Registry().WritePreview(e.ID(), e.sourceImg); err != nil {
		panic(err)
	}
	if err := e.studio.Registry().WriteContent(e.ID(), texOut); err != nil {
		panic(err)
	}
	e.savedChange = e.changes.LastChange()
}

func (e *CubeTextureEditor) Update() {}

func (e *CubeTextureEditor) OnViewportMouseEvent(event widget.ViewportMouseEvent) {
	switch event.Type {
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = true
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
	case ui.MouseEventTypeMove:
		if e.rotatingCamera {
			e.gfxCameraPitch += sprec.Degrees(float32(event.Position.Y-e.oldMouseY) / 5)
			e.gfxCameraYaw += sprec.Degrees(float32(event.Position.X-e.oldMouseX) / 5)
			e.gfxCamera.SetRotation(sprec.QuatProd(
				sprec.RotationQuat(e.gfxCameraYaw, sprec.BasisYVec3()),
				sprec.RotationQuat(e.gfxCameraPitch, sprec.BasisXVec3()),
			))
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = false
		}
	}
}

func (e *CubeTextureEditor) Scene() graphics.Scene {
	return e.gfxScene
}

func (e *CubeTextureEditor) Camera() graphics.Camera {
	return e.gfxCamera
}

func (e *CubeTextureEditor) IsAssetAccordionExpanded() bool {
	return e.propsAssetExpanded
}

func (e *CubeTextureEditor) SetAssetAccordionExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) IsSourceAccordionExpanded() bool {
	return e.propsSourceExpanded
}

func (e *CubeTextureEditor) SetSourceAccordionExpanded(expanded bool) {
	e.propsSourceExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) IsConfigAccordionExpanded() bool {
	return e.propsConfigExpanded
}

func (e *CubeTextureEditor) SetConfigAccordionExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) SourceFilename() string {
	return filepath.Base(e.sourcePath)
}

func (e *CubeTextureEditor) SourcePreview() ui.Image {
	return e.sourceImage
}

func (e *CubeTextureEditor) ChangeSourcePath(path string) {
	ch := &change.CubeTextureChangeSource{
		Controller: e,
		FromURI:    e.sourcePath,
		ToURI:      path,
	}
	if err := e.changes.Push(ch); err != nil {
		panic(err)
	}
	e.studio.NotifyChanged()
}

func (e *CubeTextureEditor) SetSourcePath(path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(e.studio.ProjectDir(), path)
	}

	relativePath, err := filepath.Rel(e.studio.ProjectDir(), path)
	if err != nil {
		log.Printf("cannot convert to relative dir: %v", err)
	} else {
		path = relativePath
	}
	e.sourcePath = path
	e.NotifyChanged()
}

func (e *CubeTextureEditor) ReloadSource() error {
	path := e.sourcePath
	if !filepath.IsAbs(path) {
		path = filepath.Join(e.studio.ProjectDir(), path)
	}

	img, err := e.openImage(path)
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

	e.Alter(func() {
		e.changePreviewImage(e.studio.Registry().PreparePreview(img))
		e.changeGraphicsImage(graphics.CubeTextureDefinition{
			Dimension:      cubeImg.Dimension,
			WrapS:          graphics.WrapClampToEdge,
			WrapT:          graphics.WrapClampToEdge,
			MinFilter:      graphics.FilterNearest,
			MagFilter:      graphics.FilterNearest,
			InternalFormat: graphics.InternalFormatRGBA32F,
			DataFormat:     graphics.DataFormatRGBA32F,
			FrontSideData:  cubeImg.RGBA32FData(pack.CubeSideFront),
			BackSideData:   cubeImg.RGBA32FData(pack.CubeSideRear),
			LeftSideData:   cubeImg.RGBA32FData(pack.CubeSideLeft),
			RightSideData:  cubeImg.RGBA32FData(pack.CubeSideRight),
			TopSideData:    cubeImg.RGBA32FData(pack.CubeSideTop),
			BottomSideData: cubeImg.RGBA32FData(pack.CubeSideBottom),
		})
	})
	return nil
}

func (e *CubeTextureEditor) RenderProperties() co.Instance {
	return co.New(view.CubeTextureProperties, func() {
		co.WithData(e)
	})
}

func (e *CubeTextureEditor) Destroy() {
	if e.gfxImage != nil {
		e.gfxImage.Delete()
	}
}

func (e *CubeTextureEditor) changePreviewImage(img image.Image) {
	// TODO: Erase old image
	e.sourceImg = img
	e.sourceImage = co.CreateImage(img)
	e.NotifyChanged()
}

func (e *CubeTextureEditor) changeGraphicsImage(definition graphics.CubeTextureDefinition) {
	if e.gfxImage != nil {
		e.gfxImage.Delete()
	}
	e.definition = definition
	e.gfxImage = e.gfxEngine.CreateCubeTexture(definition)
	e.gfxScene.Sky().SetSkybox(e.gfxImage)
}

func (e *CubeTextureEditor) openImage(path string) (image.Image, error) {
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
