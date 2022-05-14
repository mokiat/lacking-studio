package controller

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-studio/internal/studio/change"
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var dirLight *graphics.Light

func NewTwoDTextureEditor(studio *Studio, resource *data.Resource) (*TwoDTextureEditor, error) {
	gfxScene := studio.GraphicsEngine().CreateScene()
	gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.2, 0.2, 0.2))

	dirLight = gfxScene.CreateDirectionalLight()
	dirLight.SetIntensity(sprec.NewVec3(1.0, 1.0, 1.0))
	dirLight.SetRotation(sprec.IdentityQuat())

	gfxCamera := gfxScene.CreateCamera()
	gfxCamera.SetPosition(sprec.NewVec3(0.0, 0.0, 3.0))
	gfxCamera.SetRotation(sprec.IdentityQuat())
	gfxCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	gfxCamera.SetFoV(sprec.Degrees(66))
	gfxCamera.SetAutoExposure(false)
	gfxCamera.SetExposure(3.14)
	gfxCamera.SetAutoFocus(false)

	var assetImage asset.TwoDTexture
	if err := resource.LoadContent(&assetImage); err != nil {
		return nil, fmt.Errorf("failed load content: %w", err)
	}
	result := &TwoDTextureEditor{
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
			&change.TwoDTextureData{
				Controller: result,
				ToAsset:    assetImage,
			},
			&change.TwoDTextureWrapS{
				Controller: result,
				ToWrap:     assetImage.WrapModeS,
			},
			&change.TwoDTextureWrapT{
				Controller: result,
				ToWrap:     assetImage.WrapModeT,
			},
			&change.TwoDTextureMinFilter{
				Controller: result,
				ToFilter:   assetImage.MinFilter,
			},
			&change.TwoDTextureMagFilter{
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

var _ model.Editor = (*TwoDTextureEditor)(nil)
var _ model.TwoDTextureEditor = (*TwoDTextureEditor)(nil)

type TwoDTextureEditor struct {
	BaseEditor

	studio      *Studio
	resource    *data.Resource
	savedChange history.Change

	propsAssetExpanded  bool
	propsConfigExpanded bool

	gfxEngine       *graphics.Engine
	gfxScene        *graphics.Scene
	gfxCamera       *graphics.Camera
	gfxCameraPitch  sprec.Angle
	gfxCameraYaw    sprec.Angle
	gfxCameraFoV    sprec.Angle
	gfxMesh         *graphics.Mesh
	gfxMeshTemplate *graphics.MeshTemplate
	gfxMaterial     *graphics.Material
	gfxImage        *graphics.TwoDTexture

	assetImage asset.TwoDTexture

	rotatingCamera bool
	oldMouseX      int
	oldMouseY      int
}

func (e *TwoDTextureEditor) API() render.API {
	return e.studio.api
}

func (e *TwoDTextureEditor) IsPropertiesVisible() bool {
	return e.studio.IsPropertiesVisible()
}

func (e *TwoDTextureEditor) ID() string {
	return e.resource.ID()
}

func (e *TwoDTextureEditor) ChangeName(newName string) {
	e.changes.Push(&change.TwoDTextureName{
		Controller: e,
		From:       e.resource.Name(),
		To:         newName,
	})
}

func (e *TwoDTextureEditor) SetName(name string) {
	e.resource.SetName(name)
	e.studio.NotifyChanged()
}

func (e *TwoDTextureEditor) Name() string {
	return e.resource.Name()
}

func (e *TwoDTextureEditor) Icon() *ui.Image {
	return co.OpenImage("resources/icons/texture.png")
}

func (e *TwoDTextureEditor) CanSave() bool {
	return e.savedChange != e.changes.LastChange()
}

func (e *TwoDTextureEditor) Save() error {
	previewImage := image.NewRGBA(image.Rect(0, 0, data.PreviewSize, data.PreviewSize)) // TODO: Use snapshot
	if err := e.resource.Save(); err != nil {
		return fmt.Errorf("error saving resource: %w", err)
	}
	if err := e.resource.SavePreview(previewImage); err != nil {
		return fmt.Errorf("error saving preview: %w", err)
	}
	if err := e.resource.SaveContent(&e.assetImage); err != nil {
		return fmt.Errorf("error saving content: %w", err)
	}
	e.savedChange = e.changes.LastChange()
	return nil
}

func (e *TwoDTextureEditor) Update() {
	transform := sprec.Mat4MultiProd(
		sprec.RotationMat4(-e.gfxCameraYaw, 0.0, 1.0, 0.0),
		sprec.RotationMat4(-e.gfxCameraPitch, 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, 3.0),
	)
	e.gfxCamera.SetPosition(transform.Translation())
	e.gfxCamera.SetRotation(matrixToQuat(transform))
	e.gfxCamera.SetFoV(e.gfxCameraFoV)
}

// TODO: Move to gomath library.
// This is calculated by inversing the formulas for
// quat.OrientationX, quat.OrientationY and quat.OrientationZ.
func matrixToQuat(matrix sprec.Mat4) sprec.Quat {
	sqrX := (1.0 + matrix.M11 - matrix.M22 - matrix.M33) / 4.0
	sqrY := (1.0 - matrix.M11 + matrix.M22 - matrix.M33) / 4.0
	sqrZ := (1.0 - matrix.M11 - matrix.M22 + matrix.M33) / 4.0

	var x, y, z, w float32
	if sqrZ > sqrX && sqrZ > sqrY {
		// Z is largest
		if sprec.Abs(sqrZ) < 0.0000001 {
			return sprec.IdentityQuat()
		}
		z = sprec.Sqrt(sqrZ)
		x = (matrix.M31 + matrix.M13) / (4 * z)
		y = (matrix.M32 + matrix.M23) / (4 * z)
		w = (matrix.M21 - matrix.M12) / (4 * z)
	} else if sqrY > sqrX {
		// Y is largest
		if sprec.Abs(sqrY) < 0.0000001 {
			return sprec.IdentityQuat()
		}
		y = sprec.Sqrt(sqrY)
		x = (matrix.M21 + matrix.M12) / (4 * y)
		z = (matrix.M32 + matrix.M23) / (4 * y)
		w = (matrix.M13 - matrix.M31) / (4 * y)
	} else {
		// X is largest
		if sprec.Abs(sqrX) < 0.0000001 {
			return sprec.IdentityQuat()
		}
		x = sprec.Sqrt(sqrX)
		y = (matrix.M21 + matrix.M12) / (4 * x)
		z = (matrix.M31 + matrix.M13) / (4 * x)
		w = (matrix.M32 - matrix.M23) / (4 * x)
	}
	return sprec.UnitQuat(sprec.NewQuat(w, x, y, z))
}

func (e *TwoDTextureEditor) OnViewportMouseEvent(event mat.ViewportMouseEvent) bool {
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

func (e *TwoDTextureEditor) Scene() *graphics.Scene {
	return e.gfxScene
}

func (e *TwoDTextureEditor) Camera() *graphics.Camera {
	return e.gfxCamera
}

func (e *TwoDTextureEditor) IsAssetAccordionExpanded() bool {
	return e.propsAssetExpanded
}

func (e *TwoDTextureEditor) SetAssetAccordionExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) IsConfigAccordionExpanded() bool {
	return e.propsConfigExpanded
}

func (e *TwoDTextureEditor) SetConfigAccordionExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeSourcePath(path string) {
	err := e.Alter(func() error {
		if !filepath.IsAbs(path) {
			path = filepath.Join(e.studio.ProjectDir(), path)
		}

		img, err := e.openImage(path)
		if err != nil {
			return fmt.Errorf("failed to open source image: %w", err)
		}

		twodImg := pack.BuildImageResource(img)

		ch := &change.TwoDTextureData{
			Controller: e,
			FromAsset:  e.assetImage,
			ToAsset: asset.TwoDTexture{
				Width:  uint16(twodImg.Width),
				Height: uint16(twodImg.Height),
				Format: asset.TexelFormatRGBA8,
				Data:   twodImg.RGBA8Data(),
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

func (e *TwoDTextureEditor) SetAssetData(data asset.TwoDTexture) {
	e.assetImage.Width = data.Width
	e.assetImage.Height = data.Height
	e.assetImage.Format = data.Format
	e.assetImage.Data = data.Data
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) WrapS() asset.WrapMode {
	return e.assetImage.WrapModeS
}

func (e *TwoDTextureEditor) SetWrapS(mode asset.WrapMode) {
	e.assetImage.WrapModeS = mode
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeWrapS(wrap asset.WrapMode) {
	e.changes.Push(&change.TwoDTextureWrapS{
		Controller: e,
		FromWrap:   e.assetImage.WrapModeS,
		ToWrap:     wrap,
	})
}

func (e *TwoDTextureEditor) WrapT() asset.WrapMode {
	return e.assetImage.WrapModeT
}

func (e *TwoDTextureEditor) SetWrapT(mode asset.WrapMode) {
	e.assetImage.WrapModeT = mode
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeWrapT(wrap asset.WrapMode) {
	e.changes.Push(&change.TwoDTextureWrapT{
		Controller: e,
		FromWrap:   e.assetImage.WrapModeT,
		ToWrap:     wrap,
	})
}

func (e *TwoDTextureEditor) MinFilter() asset.FilterMode {
	return e.assetImage.MinFilter
}

func (e *TwoDTextureEditor) SetMinFilter(filter asset.FilterMode) {
	e.assetImage.MinFilter = filter
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeMinFilter(filter asset.FilterMode) {
	e.changes.Push(&change.TwoDTextureMinFilter{
		Controller: e,
		FromFilter: e.assetImage.MinFilter,
		ToFilter:   filter,
	})
}

func (e *TwoDTextureEditor) MagFilter() asset.FilterMode {
	return e.assetImage.MagFilter
}

func (e *TwoDTextureEditor) SetMagFilter(filter asset.FilterMode) {
	e.assetImage.MagFilter = filter
	e.rebuildGraphicsImage()
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeMagFilter(filter asset.FilterMode) {
	e.changes.Push(&change.TwoDTextureMagFilter{
		Controller: e,
		FromFilter: e.assetImage.MagFilter,
		ToFilter:   filter,
	})
}

func (e *TwoDTextureEditor) DataFormat() asset.TexelFormat {
	return e.assetImage.Format
}

func (e *TwoDTextureEditor) ChangeDataFormat(format asset.TexelFormat) {
	// TODO
}

func (e *TwoDTextureEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(view.TwoDTexture, func() {
		co.WithData(e)
		co.WithLayoutData(layoutData)
	})
}

func (e *TwoDTextureEditor) Destroy() {
	e.gfxScene.Delete()
	if e.gfxMesh != nil {
		e.gfxMesh.Delete()
	}
	if e.gfxMeshTemplate != nil {
		e.gfxMeshTemplate.Delete()
	}
	if e.gfxMaterial != nil {
		e.gfxMaterial.Delete()
	}
	if e.gfxImage != nil {
		e.gfxImage.Delete()
	}
}

func (e *TwoDTextureEditor) openImage(path string) (image.Image, error) {
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

func (e *TwoDTextureEditor) rebuildGraphicsImage() {
	oldImage := e.gfxImage
	oldMaterial := e.gfxMaterial
	oldMeshTemplate := e.gfxMeshTemplate
	oldMesh := e.gfxMesh

	definition := e.buildGraphicsDefinition(e.assetImage)
	e.gfxImage = e.gfxEngine.CreateTwoDTexture(definition)

	e.gfxMaterial = e.gfxEngine.CreatePBRMaterial(graphics.PBRMaterialDefinition{
		BackfaceCulling: false,
		AlphaBlending:   false,
		AlphaTesting:    false,
		Metalness:       0.0,
		Roughness:       0.5,
		AlbedoColor:     sprec.NewVec4(1.0, 1.0, 1.0, 1.0),
		AlbedoTexture:   e.gfxImage,
	})

	quadCount := 5
	vertexSize := 3*4 + 3*4 + 2*4
	vertexData := make([]byte, 4*vertexSize*quadCount)
	vertexPlotter := buffer.NewPlotter(vertexData, binary.LittleEndian)

	renderQuad := func(vertexPlotter *buffer.Plotter, offset sprec.Vec3, texOffset sprec.Vec2) {
		vertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(-0.5, 0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(0.0, 1.0), texOffset),
		}.Serialize(vertexPlotter)
		vertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(-0.5, -0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(0.0, 0.0), texOffset),
		}.Serialize(vertexPlotter)
		vertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(0.5, -0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(1.0, 0.0), texOffset),
		}.Serialize(vertexPlotter)
		vertex{
			Coord:    sprec.Vec3Sum(sprec.NewVec3(0.5, 0.5, 0.0), offset),
			TexCoord: sprec.Vec2Sum(sprec.NewVec2(1.0, 1.0), texOffset),
		}.Serialize(vertexPlotter)
	}

	renderQuad(vertexPlotter, sprec.NewVec3(0.0, 0.0, 0.0), sprec.NewVec2(0.0, 0.0))
	renderQuad(vertexPlotter, sprec.NewVec3(0.0, 1.01, 0.0), sprec.NewVec2(0.0, 1.0))
	renderQuad(vertexPlotter, sprec.NewVec3(0.0, -1.01, 0.0), sprec.NewVec2(0.0, -1.0))
	renderQuad(vertexPlotter, sprec.NewVec3(-1.01, 0.0, 0.0), sprec.NewVec2(-1.0, 0.0))
	renderQuad(vertexPlotter, sprec.NewVec3(1.01, 0.0, 0.0), sprec.NewVec2(1.0, 0.0))

	indexData := make([]byte, 6*2*quadCount)
	indexPlotter := buffer.NewPlotter(indexData, binary.LittleEndian)
	for i := uint16(0); i < uint16(quadCount); i++ {
		indexPlotter.PlotUint16(0 + i*4)
		indexPlotter.PlotUint16(1 + i*4)
		indexPlotter.PlotUint16(2 + i*4)

		indexPlotter.PlotUint16(0 + i*4)
		indexPlotter.PlotUint16(2 + i*4)
		indexPlotter.PlotUint16(3 + i*4)
	}

	e.gfxMeshTemplate = e.gfxEngine.CreateMeshTemplate(graphics.MeshTemplateDefinition{
		VertexData: vertexData,
		VertexFormat: graphics.VertexFormat{
			HasCoord:            true,
			CoordOffsetBytes:    0,
			CoordStrideBytes:    vertexSize,
			HasNormal:           true,
			NormalOffsetBytes:   3 * 4,
			NormalStrideBytes:   vertexSize,
			HasTexCoord:         true,
			TexCoordOffsetBytes: 3*4 + 3*4,
			TexCoordStrideBytes: vertexSize,
		},
		IndexData:   indexData,
		IndexFormat: graphics.IndexFormatU16,
		SubMeshes: []graphics.SubMeshTemplateDefinition{
			{
				Primitive:   graphics.PrimitiveTriangles,
				IndexOffset: 0,
				IndexCount:  6 * quadCount,
				Material:    e.gfxMaterial,
			},
		},
	})

	e.gfxMesh = e.gfxScene.CreateMesh(e.gfxMeshTemplate)

	if oldMesh != nil {
		oldMesh.Delete()
	}
	if oldMeshTemplate != nil {
		oldMeshTemplate.Delete()
	}
	if oldMaterial != nil {
		oldMaterial.Delete()
	}
	if oldImage != nil {
		oldImage.Delete()
	}
}

type vertex struct {
	Coord    sprec.Vec3
	TexCoord sprec.Vec2
}

func (v vertex) Serialize(plotter *buffer.Plotter) {
	plotter.PlotFloat32(v.Coord.X)
	plotter.PlotFloat32(v.Coord.Y)
	plotter.PlotFloat32(v.Coord.Z)
	plotter.PlotFloat32(0.0)
	plotter.PlotFloat32(0.0)
	plotter.PlotFloat32(1.0)
	plotter.PlotFloat32(v.TexCoord.X)
	plotter.PlotFloat32(v.TexCoord.Y)
}

func (e *TwoDTextureEditor) buildGraphicsDefinition(src asset.TwoDTexture) graphics.TwoDTextureDefinition {
	return graphics.TwoDTextureDefinition{
		Width:          int(src.Width),
		Height:         int(src.Height),
		WrapS:          e.assetToGraphicsWrap(src.WrapModeS),
		WrapT:          e.assetToGraphicsWrap(src.WrapModeT),
		MinFilter:      e.assetToGraphicsFilter(src.MinFilter),
		MagFilter:      e.assetToGraphicsFilter(src.MagFilter),
		UseAnisotropy:  false,
		InternalFormat: e.assetFormatToInternalFormat(src.Format),
		DataFormat:     e.assetFormatToDataFormat(src.Format),
		Data:           src.Data,
	}
}

func (e *TwoDTextureEditor) assetToGraphicsWrap(wrap asset.WrapMode) graphics.Wrap {
	switch wrap {
	case asset.WrapModeClampToEdge:
		return graphics.WrapClampToEdge
	case asset.WrapModeRepeat:
		return graphics.WrapRepeat
	case asset.WrapModeMirroredRepeat:
		return graphics.WrapMirroredRepat
	default:
		panic(fmt.Errorf("unsupported wrap: %v", wrap))
	}
}

func (e *TwoDTextureEditor) assetToGraphicsFilter(filter asset.FilterMode) graphics.Filter {
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

func (e *TwoDTextureEditor) assetFormatToInternalFormat(format asset.TexelFormat) graphics.InternalFormat {
	switch format {
	case asset.TexelFormatRGBA8:
		return graphics.InternalFormatRGBA8
	case asset.TexelFormatRGBA32F:
		return graphics.InternalFormatRGBA32F
	default:
		panic(fmt.Errorf("unsupported format: %v", format))
	}
}

func (e *TwoDTextureEditor) assetFormatToDataFormat(format asset.TexelFormat) graphics.DataFormat {
	switch format {
	case asset.TexelFormatRGBA8:
		return graphics.DataFormatRGBA8
	case asset.TexelFormatRGBA32F:
		return graphics.DataFormatRGBA32F
	default:
		panic(fmt.Errorf("unsupported format: %v", format))
	}
}
