package controller

import (
	"fmt"
	"image"
	"os"

	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/change"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/visualization"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var (
	// TODO: Move these to model package so that they are usable from other packages
	CubeTextureEditorChange                        = observer.StringChange("cube_texture_editor")
	CubeTextureEditorAssetAccordionExpandedChange  = observer.ExtendChange(CubeTextureEditorChange, observer.StringChange("asset_accordion_expanded"))
	CubeTextureEditorConfigAccordionExpandedChange = observer.ExtendChange(CubeTextureEditorChange, observer.StringChange("config_accordion_expanded"))
)

func NewCubeTextureEditor(studio *Studio, texModel *model.CubeTexture) *CubeTextureEditor {
	target := observer.NewTarget()
	studioSubscription := observer.WireTargets(studio.Target(), target)
	texModelSubscription := observer.WireTargets(texModel.Target(), target)

	return &CubeTextureEditor{
		BaseEditor: NewBaseEditor(),

		target: target,

		studio:               studio,
		studioSubscription:   studioSubscription,
		texModel:             texModel,
		texModelSubscription: texModelSubscription,

		propsAssetExpanded:  false,
		propsConfigExpanded: true,

		viz: visualization.NewCubeTexture(studio.api /* FIXME */, studio.GraphicsEngine(), texModel),
	}
}

var _ model.CubeTextureEditor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	BaseEditor

	target *observer.Target

	studio               *Studio
	studioSubscription   *observer.Subscription
	texModel             *model.CubeTexture
	texModelSubscription *observer.Subscription

	propsAssetExpanded  bool
	propsConfigExpanded bool

	viz *visualization.CubeTexture
}

func (e *CubeTextureEditor) Target() *observer.Target {
	return e.target
}

func (e *CubeTextureEditor) ID() string {
	return e.texModel.ID()
}

func (e *CubeTextureEditor) Name() string {
	return e.texModel.Name()
}

func (e *CubeTextureEditor) Icon() *ui.Image {
	return co.OpenImage("icons/texture.png")
}

func (e *CubeTextureEditor) Save() error {
	previewImg := e.viz.TakeSnapshot(ui.Size{
		Width:  data.PreviewSize,
		Height: data.PreviewSize,
	})
	e.texModel.SetPreviewImage(previewImg)

	if err := e.texModel.Save(); err != nil {
		return fmt.Errorf("error saving texture model %w", err)
	}
	return e.BaseEditor.Save()
}

func (e *CubeTextureEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(view.CubeTexture, func() {
		co.WithData(e)
		co.WithLayoutData(layoutData)
	})
}

func (e *CubeTextureEditor) Destroy() {
	e.viz.Destroy()
	e.texModelSubscription.Delete()
	e.studioSubscription.Delete()
}

func (e *CubeTextureEditor) IsPropertiesVisible() bool {
	// TODO: Figure out how to untie this. Either create subscription hell
	// or maybe allow editors to create their own toolbar buttons in the studio.
	return e.studio.IsPropertiesVisible()
}

func (e *CubeTextureEditor) IsAssetAccordionExpanded() bool {
	return e.propsAssetExpanded
}

func (e *CubeTextureEditor) SetAssetAccordionExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.target.SignalChange(TwoDTextureEditorAssetAccordionExpandedChange)
}

func (e *CubeTextureEditor) IsConfigAccordionExpanded() bool {
	return e.propsConfigExpanded
}

func (e *CubeTextureEditor) SetConfigAccordionExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.target.SignalChange(TwoDTextureEditorConfigAccordionExpandedChange)
}

func (e *CubeTextureEditor) Filtering() asset.FilterMode {
	return e.texModel.Filtering()
}

func (e *CubeTextureEditor) DataFormat() asset.TexelFormat {
	return e.texModel.Format()
}

func (e *CubeTextureEditor) ChangeName(newName string) {
	e.changes.Push(change.ResourceName(e.texModel,
		change.ResourceNameState{
			Value: e.texModel.Name(),
		},
		change.ResourceNameState{
			Value: newName,
		},
	))

	// FIXME: Figure out how to avoid this:
	e.studio.NotifyChanged()
}

func (e *CubeTextureEditor) ChangeContent(path string) {
	img, err := e.openImage(path)
	if err != nil {
		e.studio.HandleError(fmt.Errorf("failed to open source image: %w", err))
		return
	}

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
		e.studio.HandleError(fmt.Errorf("failed to build cube image: %w", err))
		return
	}

	ch := change.CubeTextureContent(e.texModel,
		change.CubeTextureContentState{
			Dimension:  e.texModel.Dimension(),
			Format:     e.texModel.Format(),
			FrontData:  e.texModel.FrontData(),
			BackData:   e.texModel.BackData(),
			LeftData:   e.texModel.LeftData(),
			RightData:  e.texModel.RightData(),
			TopData:    e.texModel.TopData(),
			BottomData: e.texModel.BottomData(),
		},
		change.CubeTextureContentState{
			Dimension:  cubeImg.Dimension,
			Format:     asset.TexelFormatRGBA32F,
			FrontData:  cubeImg.RGBA32FData(pack.CubeSideFront),
			BackData:   cubeImg.RGBA32FData(pack.CubeSideRear),
			LeftData:   cubeImg.RGBA32FData(pack.CubeSideLeft),
			RightData:  cubeImg.RGBA32FData(pack.CubeSideRight),
			TopData:    cubeImg.RGBA32FData(pack.CubeSideTop),
			BottomData: cubeImg.RGBA32FData(pack.CubeSideBottom),
		},
	)
	if err := e.changes.Push(ch); err != nil {
		e.studio.HandleError(fmt.Errorf("failed to apply change: %w", err))
		return
	}
}

func (e *CubeTextureEditor) ChangeFiltering(filter asset.FilterMode) {
	e.changes.Push(change.Filtering(
		e.texModel,
		change.FilteringState{
			Value: e.texModel.Filtering(),
		},
		change.FilteringState{
			Value: filter,
		},
	))
}

func (e *CubeTextureEditor) ChangeDataFormat(format asset.TexelFormat) {
	// TODO
}

func (e *CubeTextureEditor) Visualization() model.Visualization {
	return e.viz
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
