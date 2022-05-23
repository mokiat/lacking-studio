package controller

import (
	"fmt"
	"image"
	"os"

	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/change"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/visualization"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewTwoDTextureEditor(studio *Studio, texModel *model.TwoDTexture) (*TwoDTextureEditor, error) {
	viz := visualization.NewTwoDTexture(studio.api /* FIXME */, studio.GraphicsEngine(), texModel)

	return &TwoDTextureEditor{
		BaseEditor: NewBaseEditor(),

		studio:   studio,
		texModel: texModel,

		propsAssetExpanded:  false,
		propsConfigExpanded: true,

		viz: viz,
	}, nil
}

var _ model.TwoDTextureEditor = (*TwoDTextureEditor)(nil)

type TwoDTextureEditor struct {
	BaseEditor

	studio   *Studio
	texModel *model.TwoDTexture

	propsAssetExpanded  bool
	propsConfigExpanded bool

	viz *visualization.TwoDTexture
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
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) IsConfigAccordionExpanded() bool {
	return e.propsConfigExpanded
}

func (e *TwoDTextureEditor) SetConfigAccordionExpanded(expanded bool) {
	e.propsConfigExpanded = expanded
	e.NotifyChanged()
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
	e.changes.Push(change.TwoDTextureName(e.texModel,
		change.TwoDTextureNameState{
			Value: e.texModel.Name(),
		},
		change.TwoDTextureNameState{
			Value: newName,
		},
	))

	// TODO: Remove. This should come as a notification from the model.
	e.studio.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeContent(path string) {
	// err := e.Alter(func() error {
	// 	if !filepath.IsAbs(path) {
	// 		path = filepath.Join(e.studio.ProjectDir(), path)
	// 	}

	// 	img, err := e.openImage(path)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to open source image: %w", err)
	// 	}

	// 	twodImg := pack.BuildImageResource(img)

	// 	ch := &change.TwoDTextureData{
	// 		Controller: e,
	// 		FromAsset:  e.assetImage,
	// 		ToAsset: asset.TwoDTexture{
	// 			Width:  uint16(twodImg.Width),
	// 			Height: uint16(twodImg.Height),
	// 			Format: asset.TexelFormatRGBA8,
	// 			Data:   twodImg.RGBA8Data(),
	// 		},
	// 	}
	// 	if err := e.changes.Push(ch); err != nil {
	// 		return fmt.Errorf("failed to apply change: %w", err)
	// 	}
	// 	e.studio.NotifyChanged()
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err) // TODO
	// }
}

func (e *TwoDTextureEditor) ChangeWrapping(wrap asset.WrapMode) {
	e.changes.Push(change.TwoDTextureWrapping(e.texModel,
		change.TwoDTextureWrappingState{
			Value: e.texModel.Wrapping(),
		},
		change.TwoDTextureWrappingState{
			Value: wrap,
		},
	))

	// TODO: Remove. This should come as a notification from the model.
	e.NotifyChanged()
}

func (e *TwoDTextureEditor) ChangeFiltering(filter asset.FilterMode) {
	e.changes.Push(change.TwoDTextureFiltering(
		e.texModel,
		change.TwoDTextureFilteringState{
			Value: e.texModel.Filtering(),
		},
		change.TwoDTextureFilteringState{
			Value: filter,
		},
	))

	// TODO: Remove. This should come as a notification from the model.
	e.NotifyChanged()
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
