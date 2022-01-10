package controller

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/studio/change"
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/widget"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewCubeTextureEditor(studio *Studio, resource *asset.Resource) (*CubeTextureEditor, error) {
	gfxScene := studio.GraphicsEngine().CreateScene()
	gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.0, 0.0, 0.0))

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetPosition(sprec.NewVec3(0.0, 0.0, 0.0))
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	gfxCamera.SetFoV(sprec.Degrees(66))
	gfxCamera.SetAutoExposure(true)
	gfxCamera.SetExposure(1.0)
	gfxCamera.SetAutoFocus(false)

	var assetImage asset.CubeTexture
	if err := studio.Registry().ReadContent(resource.GUID, &assetImage); err != nil {
		return nil, fmt.Errorf("failed to open asset %q: %w", resource.GUID, err)
	}
	result := &CubeTextureEditor{
		BaseEditor: NewBaseEditor(),

		studio:   studio,
		resource: resource,

		propsAssetExpanded:  false,
		propsConfigExpanded: true,

		gfxEngine:      studio.GraphicsEngine(),
		gfxScene:       gfxScene,
		gfxCamera:      gfxCamera,
		gfxCameraPitch: sprec.Degrees(0),
		gfxCameraYaw:   sprec.Degrees(0),
		gfxCameraFoV:   sprec.Degrees(66),

		assetImage: assetImage,
	}
	result.savedChange = &change.Combined{
		Changes: []history.Change{
			&change.CubeTextureData{
				Controller: result,
				ToAsset:    assetImage,
			},
			&change.CubeTextureMinFilter{
				Controller: result,
				ToFilter:   assetImage.MinFilter,
			},
			&change.CubeTextureMagFilter{
				Controller: result,
				ToFilter:   assetImage.MagFilter,
			},
		},
	}
	if err := result.changes.Push(result.savedChange); err != nil {
		return nil, fmt.Errorf("failed to init editor: %w", err)
	}
	return result, nil
}

var _ model.Editor = (*CubeTextureEditor)(nil)
var _ model.CubeTextureEditor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	BaseEditor

	studio      *Studio
	resource    *asset.Resource
	savedChange history.Change

	propsAssetExpanded  bool
	propsConfigExpanded bool

	gfxEngine      graphics.Engine
	gfxScene       graphics.Scene
	gfxCamera      graphics.Camera
	gfxCameraPitch sprec.Angle
	gfxCameraYaw   sprec.Angle
	gfxCameraFoV   sprec.Angle
	gfxImage       graphics.CubeTexture

	assetImage asset.CubeTexture

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
}

func (e *CubeTextureEditor) IsPropertiesVisible() bool {
	return e.studio.IsPropertiesVisible()
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
	previewImage := image.NewRGBA(image.Rect(0, 0, 128, 128)) // TODO: Use snapshot
	if err := e.studio.Registry().WritePreview(e.ID(), previewImage); err != nil {
		return fmt.Errorf("failed to write preview image: %w", err)
	}
	if err := e.studio.Registry().WriteContent(e.ID(), &e.assetImage); err != nil {
		return fmt.Errorf("failed to write content image: %w", err)
	}
	e.savedChange = e.changes.LastChange()
	return nil
}

func (e *CubeTextureEditor) Update() {
	e.gfxCamera.SetFoV(e.gfxCameraFoV)
	e.gfxCamera.SetRotation(sprec.QuatProd(
		sprec.RotationQuat(e.gfxCameraYaw, sprec.BasisYVec3()),
		sprec.RotationQuat(e.gfxCameraPitch, sprec.BasisXVec3()),
	))
}

func (e *CubeTextureEditor) OnViewportMouseEvent(event widget.ViewportMouseEvent) bool {
	switch event.Type {
	case ui.MouseEventTypeDown:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = true
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
		return true
	case ui.MouseEventTypeMove:
		if e.rotatingCamera {
			e.gfxCameraPitch += sprec.Degrees(float32(event.Position.Y-e.oldMouseY) / 5)
			e.gfxCameraYaw += sprec.Degrees(float32(event.Position.X-e.oldMouseX) / 5)
			e.oldMouseX = event.Position.X
			e.oldMouseY = event.Position.Y
		}
		return true
	case ui.MouseEventTypeUp:
		if event.Button == ui.MouseButtonMiddle {
			e.rotatingCamera = false
		}
		return true
	case ui.MouseEventTypeScroll:
		fov := e.gfxCameraFoV.Degrees()
		fov -= 2 * float32(event.ScrollY)
		fov = sprec.Clamp(fov, 0.1, 179.0)
		e.gfxCameraFoV = sprec.Degrees(fov)
		return true
	default:
		return false
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

func (e *CubeTextureEditor) IsConfigAccordionExpanded() bool {
	return e.propsConfigExpanded
}

func (e *CubeTextureEditor) SetConfigAccordionExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.NotifyChanged()
}

func (e *CubeTextureEditor) ChangeSourcePath(path string) {
	err := e.Alter(func() error {
		if !filepath.IsAbs(path) {
			path = filepath.Join(e.studio.ProjectDir(), path)
		}

		img, err := e.openImage(path)
		if err != nil {
			return fmt.Errorf("failed to open source image: %w", err)
		}

		// TODO: Use GPU for all of this!
		twodImg := pack.BuildImageResource(img)
		dimension := twodImg.Height / 2 // TODO: Allow user to configure
		frontPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideFront, dimension)
		rearPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideRear, dimension)
		leftPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideLeft, dimension)
		rightPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideRight, dimension)
		topPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideTop, dimension)
		bottomPackImg := pack.BuildCubeSideFromEquirectangularScaled(twodImg, pack.CubeSideBottom, dimension)
		cubeImg, err := pack.BuildCube(frontPackImg, rearPackImg, leftPackImg, rightPackImg, topPackImg, bottomPackImg, 0)
		if err != nil {
			return fmt.Errorf("failed to build cube image: %w", err)
		}

		ch := &change.CubeTextureData{
			Controller: e,
			FromAsset:  e.assetImage,
			ToAsset: asset.CubeTexture{
				Dimension: uint16(cubeImg.Dimension),
				Format:    asset.TexelFormatRGBA32F,
				FrontSide: asset.CubeTextureSide{
					Data: cubeImg.RGBA32FData(pack.CubeSideFront),
				},
				BackSide: asset.CubeTextureSide{
					Data: cubeImg.RGBA32FData(pack.CubeSideRear),
				},
				LeftSide: asset.CubeTextureSide{
					Data: cubeImg.RGBA32FData(pack.CubeSideLeft),
				},
				RightSide: asset.CubeTextureSide{
					Data: cubeImg.RGBA32FData(pack.CubeSideRight),
				},
				TopSide: asset.CubeTextureSide{
					Data: cubeImg.RGBA32FData(pack.CubeSideTop),
				},
				BottomSide: asset.CubeTextureSide{
					Data: cubeImg.RGBA32FData(pack.CubeSideBottom),
				},
			},
		}
		if err := e.changes.Push(ch); err != nil {
			return fmt.Errorf("failed to apply change: %w", err)
		}
		e.studio.NotifyChanged()
		return nil
	})
	if err != nil {
		panic(err) // TODO
	}
}

func (e *CubeTextureEditor) SetAssetData(data asset.CubeTexture) {
	e.assetImage.Dimension = data.Dimension
	e.assetImage.Format = data.Format
	e.assetImage.FrontSide = data.FrontSide
	e.assetImage.BackSide = data.BackSide
	e.assetImage.LeftSide = data.LeftSide
	e.assetImage.RightSide = data.RightSide
	e.assetImage.TopSide = data.TopSide
	e.assetImage.BottomSide = data.BottomSide
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *CubeTextureEditor) MinFilter() asset.FilterMode {
	return e.assetImage.MinFilter
}

func (e *CubeTextureEditor) SetMinFilter(filter asset.FilterMode) {
	e.assetImage.MinFilter = filter
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *CubeTextureEditor) ChangeMinFilter(filter asset.FilterMode) {
	e.changes.Push(&change.CubeTextureMinFilter{
		Controller: e,
		FromFilter: e.assetImage.MinFilter,
		ToFilter:   filter,
	})
}

func (e *CubeTextureEditor) MagFilter() asset.FilterMode {
	return e.assetImage.MagFilter
}

func (e *CubeTextureEditor) SetMagFilter(filter asset.FilterMode) {
	e.assetImage.MagFilter = filter
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *CubeTextureEditor) ChangeMagFilter(filter asset.FilterMode) {
	e.changes.Push(&change.CubeTextureMagFilter{
		Controller: e,
		FromFilter: e.assetImage.MagFilter,
		ToFilter:   filter,
	})
}

func (e *CubeTextureEditor) DataFormat() asset.TexelFormat {
	return e.assetImage.Format
}

func (e *CubeTextureEditor) ChangeDataFormat(format asset.TexelFormat) {
	// TODO
}

func (e *CubeTextureEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(view.CubeTexture, func() {
		co.WithData(e)
		co.WithLayoutData(layoutData)
	})
}

func (e *CubeTextureEditor) Destroy() {
	if e.gfxImage != nil {
		e.gfxImage.Delete()
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

func (e *CubeTextureEditor) rebuildGraphicsImage() {
	definition := e.buildGraphicsDefinition(e.assetImage)

	oldImage := e.gfxImage
	e.gfxImage = e.gfxEngine.CreateCubeTexture(definition)
	e.gfxScene.Sky().SetSkybox(e.gfxImage)
	if oldImage != nil {
		oldImage.Delete()
	}
}

func (e *CubeTextureEditor) buildGraphicsDefinition(src asset.CubeTexture) graphics.CubeTextureDefinition {
	return graphics.CubeTextureDefinition{
		Dimension:      int(src.Dimension),
		MinFilter:      e.assetToGraphicsFilter(src.MinFilter),
		MagFilter:      e.assetToGraphicsFilter(src.MagFilter),
		InternalFormat: e.assetFormatToInternalFormat(src.Format),
		DataFormat:     e.assetFormatToDataFormat(src.Format),
		FrontSideData:  src.FrontSide.Data,
		BackSideData:   src.BackSide.Data,
		LeftSideData:   src.LeftSide.Data,
		RightSideData:  src.RightSide.Data,
		TopSideData:    src.TopSide.Data,
		BottomSideData: src.BottomSide.Data,
	}
}

func (e *CubeTextureEditor) assetToGraphicsFilter(filter asset.FilterMode) graphics.Filter {
	switch filter {
	case asset.FilterModeUnspecified:
		fallthrough
	case asset.FilterModeNearest:
		return graphics.FilterNearest
	case asset.FilterModeLinear:
		return graphics.FilterLinear
	case asset.FilterModeNearestMipmapNearest:
		return graphics.FilterNearestMipmapNearest
	case asset.FilterModeNearestMipmapLinear:
		return graphics.FilterNearestMipmapLinear
	case asset.FilterModeLinearMipmapNearest:
		return graphics.FilterLinearMipmapNearest
	case asset.FilterModeLinearMipmapLinear:
		return graphics.FilterLinearMipmapLinear
	default:
		panic(fmt.Errorf("unsupported filter: %v", filter))
	}
}

func (e *CubeTextureEditor) assetFormatToInternalFormat(format asset.TexelFormat) graphics.InternalFormat {
	switch format {
	case asset.TexelFormatRGBA8:
		return graphics.InternalFormatRGBA8
	case asset.TexelFormatRGBA32F:
		return graphics.InternalFormatRGBA32F
	default:
		panic(fmt.Errorf("unsupported format: %v", format))
	}
}

func (e *CubeTextureEditor) assetFormatToDataFormat(format asset.TexelFormat) graphics.DataFormat {
	switch format {
	case asset.TexelFormatRGBA8:
		return graphics.DataFormatRGBA8
	case asset.TexelFormatRGBA32F:
		return graphics.DataFormatRGBA32F
	default:
		panic(fmt.Errorf("unsupported format: %v", format))
	}
}
