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

		sourcePath: "<none>",
	}
}

var _ Editor = (*CubeTextureEditor)(nil)
var _ model.CubeTextureEditor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	BaseEditor

	studio      *Studio
	resource    *gameasset.Resource
	savedChange history.Change

	propsAssetExpanded  bool
	propsSourceExpanded bool
	propsConfigExpanded bool

	gfxEngine      graphics.Engine
	gfxScene       graphics.Scene
	gfxCamera      graphics.Camera
	gfxCameraPitch sprec.Angle
	gfxCameraYaw   sprec.Angle

	sourceImage    image.Image
	sourcePath     string
	previewImage   image.Image
	previewUIImage ui.Image
	convertedImage *pack.CubeImage
	graphicsImage  graphics.CubeTexture

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
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

func (e *CubeTextureEditor) Save() error {
	assetTexture := e.buildAssetCubeTexture()
	if err := e.studio.Registry().WritePreview(e.ID(), e.previewImage); err != nil {
		return fmt.Errorf("failed to write preview image: %w", err)
	}
	if err := e.studio.Registry().WriteContent(e.ID(), assetTexture); err != nil {
		return fmt.Errorf("failed to write content image: %w", err)
	}
	e.savedChange = e.changes.LastChange()
	return nil
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
	return e.previewUIImage
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
	return e.Alter(func() error {
		path := e.sourcePath
		if !filepath.IsAbs(path) {
			path = filepath.Join(e.studio.ProjectDir(), path)
		}

		img, err := e.openImage(path)
		if err != nil {
			return fmt.Errorf("failed to open source image: %w", err)
		}
		e.setSourceImage(img)

		e.rebuildPreviewImage()
		e.rebuildConvertedImage()
		e.rebuildGraphicsImage()
		return nil
	})
}

func (e *CubeTextureEditor) RenderProperties() co.Instance {
	return co.New(view.CubeTextureProperties, func() {
		co.WithData(e)
	})
}

func (e *CubeTextureEditor) Destroy() {
	// TODO: Delete other images
	if e.graphicsImage != nil {
		e.graphicsImage.Delete()
	}
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

func (e *CubeTextureEditor) setSourceImage(img image.Image) {
	e.sourceImage = img
}

func (e *CubeTextureEditor) rebuildPreviewImage() {
	previewImg := e.studio.Registry().PreparePreview(e.sourceImage)
	e.setPreviewImage(previewImg)
	previewUIImage := co.CreateImage(e.previewImage)
	e.setPreviewUIImage(previewUIImage)
}

func (e *CubeTextureEditor) setPreviewImage(img image.Image) {
	e.previewImage = img
}

func (e *CubeTextureEditor) setPreviewUIImage(img ui.Image) {
	// TODO: Erase old image
	e.previewUIImage = img
	e.NotifyChanged()
}

func (e *CubeTextureEditor) rebuildConvertedImage() {
	// TODO: Do all of this on the GPU
	twodImg := pack.BuildImageResource(e.sourceImage)
	dimension := twodImg.Height / 2 // TODO: Allow user to configure
	frontPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideFront, dimension)
	rearPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideRear, dimension)
	leftPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideLeft, dimension)
	rightPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideRight, dimension)
	topPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideTop, dimension)
	bottomPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideBottom, dimension)
	cubeImg, err := pack.BuildCube(frontPackImg, rearPackImg, leftPackImg, rightPackImg, topPackImg, bottomPackImg, 0)
	if err != nil {
		panic(fmt.Errorf("failed to build cube image: %w", err))
	}
	e.setConvertedImage(cubeImg)
}

func (e *CubeTextureEditor) setConvertedImage(img *pack.CubeImage) {
	e.convertedImage = img
}

func (e *CubeTextureEditor) rebuildGraphicsImage() {
	definition := e.buildGraphicsDefinition()
	graphicsImg := e.gfxEngine.CreateCubeTexture(definition)
	e.setGraphicsImage(graphicsImg)
}

func (e *CubeTextureEditor) setGraphicsImage(img graphics.CubeTexture) {
	if e.graphicsImage != nil {
		e.graphicsImage.Delete()
	}
	e.graphicsImage = img
	e.gfxScene.Sky().SetSkybox(img)
}

func (e *CubeTextureEditor) buildGraphicsDefinition() graphics.CubeTextureDefinition {
	// TODO: Take filtering and internal format settings from configuration
	// accordion
	return graphics.CubeTextureDefinition{
		Dimension:      e.convertedImage.Dimension,
		WrapS:          graphics.WrapClampToEdge,
		WrapT:          graphics.WrapClampToEdge,
		MinFilter:      graphics.FilterNearest,
		MagFilter:      graphics.FilterNearest,
		InternalFormat: graphics.InternalFormatRGBA32F,
		DataFormat:     graphics.DataFormatRGBA32F,
		FrontSideData:  e.convertedImage.RGBA32FData(pack.CubeSideFront),
		BackSideData:   e.convertedImage.RGBA32FData(pack.CubeSideRear),
		LeftSideData:   e.convertedImage.RGBA32FData(pack.CubeSideLeft),
		RightSideData:  e.convertedImage.RGBA32FData(pack.CubeSideRight),
		TopSideData:    e.convertedImage.RGBA32FData(pack.CubeSideTop),
		BottomSideData: e.convertedImage.RGBA32FData(pack.CubeSideBottom),
	}
}

func (e *CubeTextureEditor) buildAssetCubeTexture() *asset.CubeTexture {
	definition := e.buildGraphicsDefinition()

	texOut := &asset.CubeTexture{
		Dimension: uint16(definition.Dimension),
		Format: e.calculateAssetFormatFromGraphics(
			definition.DataFormat,
			definition.InternalFormat,
		),
	}
	texOut.Sides[asset.TextureSideFront].Data = definition.FrontSideData
	texOut.Sides[asset.TextureSideBack].Data = definition.BackSideData
	texOut.Sides[asset.TextureSideLeft].Data = definition.LeftSideData
	texOut.Sides[asset.TextureSideRight].Data = definition.RightSideData
	texOut.Sides[asset.TextureSideTop].Data = definition.TopSideData
	texOut.Sides[asset.TextureSideBottom].Data = definition.BottomSideData
	return texOut
}

func (e *CubeTextureEditor) calculateAssetFormatFromGraphics(dataFormat graphics.DataFormat, internalFormat graphics.InternalFormat) asset.DataFormat {
	switch dataFormat {
	case graphics.DataFormatRGBA8:
		return asset.DataFormatRGBA8
	case graphics.DataFormatRGBA32F:
		return asset.DataFormatRGBA32F
	default:
		panic(fmt.Errorf("unknown data format: %#v", dataFormat))
	}
}
