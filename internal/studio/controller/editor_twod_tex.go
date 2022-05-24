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
	TwoDTextureEditorChange                        = observer.StringChange("twod_texture_editor")
	TwoDTextureEditorAssetAccordionExpandedChange  = observer.ExtendChange(TwoDTextureEditorChange, observer.StringChange("asset_accordion_expanded"))
	TwoDTextureEditorConfigAccordionExpandedChange = observer.ExtendChange(TwoDTextureEditorChange, observer.StringChange("config_accordion_expanded"))
)

func NewTwoDTextureEditor(studio *Studio, texModel *model.TwoDTexture) *TwoDTextureEditor {
	target := observer.NewTarget()
	studioSubscription := observer.WireTargets(studio.Target(), target)
	texModelSubscription := observer.WireTargets(texModel.Target(), target)

	return &TwoDTextureEditor{
		BaseEditor: NewBaseEditor(),

		target: target,

		studio:               studio,
		studioSubscription:   studioSubscription,
		texModel:             texModel,
		texModelSubscription: texModelSubscription,

		propsAssetExpanded:  false,
		propsConfigExpanded: true,

		viz: visualization.NewTwoDTexture(studio.api /* FIXME */, studio.GraphicsEngine(), texModel),
	}
}

var _ model.TwoDTextureEditor = (*TwoDTextureEditor)(nil)

type TwoDTextureEditor struct {
	BaseEditor

	target *observer.Target

	studio               *Studio
	studioSubscription   *observer.Subscription
	texModel             *model.TwoDTexture
	texModelSubscription *observer.Subscription

	propsAssetExpanded  bool
	propsConfigExpanded bool

	viz *visualization.TwoDTexture
}

func (e *TwoDTextureEditor) Target() *observer.Target {
	return e.target
}

func (e *TwoDTextureEditor) ID() string {
	return e.texModel.ID()
}

func (e *TwoDTextureEditor) Name() string {
	return e.texModel.Name()
}

func (e *TwoDTextureEditor) Icon() *ui.Image {
	return co.OpenImage("resources/icons/texture.png")
}

func (e *TwoDTextureEditor) Save() error {
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

func (e *TwoDTextureEditor) Render(layoutData mat.LayoutData) co.Instance {
	return co.New(view.TwoDTexture, func() {
		co.WithData(e)
		co.WithLayoutData(layoutData)
	})
}

func (e *TwoDTextureEditor) Destroy() {
	e.viz.Destroy()
	e.texModelSubscription.Delete()
	e.studioSubscription.Delete()
}

func (e *TwoDTextureEditor) IsPropertiesVisible() bool {
	// TODO: Figure out how to untie this. Either create subscription hell
	// or maybe allow editors to create their own toolbar buttons in the studio.
	return e.studio.IsPropertiesVisible()
}

func (e *TwoDTextureEditor) IsAssetAccordionExpanded() bool {
	return e.propsAssetExpanded
}

func (e *TwoDTextureEditor) SetAssetAccordionExpanded(expanded bool) {
	e.propsAssetExpanded = expanded
	e.target.SignalChange(TwoDTextureEditorAssetAccordionExpandedChange)
}

func (e *TwoDTextureEditor) IsConfigAccordionExpanded() bool {
	return e.propsConfigExpanded
}

func (e *TwoDTextureEditor) SetConfigAccordionExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.target.SignalChange(TwoDTextureEditorConfigAccordionExpandedChange)
}

func (e *TwoDTextureEditor) Wrapping() asset.WrapMode {
	return e.texModel.Wrapping()
}

func (e *TwoDTextureEditor) Filtering() asset.FilterMode {
	return e.texModel.Filtering()
}

func (e *TwoDTextureEditor) DataFormat() asset.TexelFormat {
	return e.texModel.Format()
}

func (e *TwoDTextureEditor) ChangeName(newName string) {
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

func (e *TwoDTextureEditor) ChangeContent(path string) {
	img, err := e.openImage(path)
	if err != nil {
		e.studio.HandleError(fmt.Errorf("failed to open source image: %w", err))
		return
	}
	twodImg := pack.BuildImageResource(img)

	ch := change.TwoDTextureContent(e.texModel,
		change.TwoDTextureContentState{
			Width:  e.texModel.Width(),
			Height: e.texModel.Height(),
			Format: e.texModel.Format(),
			Data:   e.texModel.Data(),
		},
		change.TwoDTextureContentState{
			Width:  twodImg.Width,
			Height: twodImg.Height,
			Format: asset.TexelFormatRGBA8,
			Data:   twodImg.RGBA8Data(),
		},
	)
	if err := e.changes.Push(ch); err != nil {
		e.studio.HandleError(fmt.Errorf("failed to apply change: %w", err))
		return
	}
}

func (e *TwoDTextureEditor) ChangeWrapping(wrap asset.WrapMode) {
	e.changes.Push(change.Wrapping(e.texModel,
		change.WrappingState{
			Value: e.texModel.Wrapping(),
		},
		change.WrappingState{
			Value: wrap,
		},
	))
}

func (e *TwoDTextureEditor) ChangeFiltering(filter asset.FilterMode) {
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

func (e *TwoDTextureEditor) ChangeDataFormat(format asset.TexelFormat) {
	// TODO
}

func (e *TwoDTextureEditor) Visualization() model.Visualization {
	return e.viz
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
