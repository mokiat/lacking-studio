package controller

import (
	"fmt"
	"image"
	"os"

	"github.com/mokiat/lacking-studio/internal/studio/data"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking-studio/internal/studio/model/change"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/visualization"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

func NewTwoDTextureEditor(studio *Studio, texModel *model.TwoDTexture) *TwoDTextureEditor {
	return &TwoDTextureEditor{
		BaseEditor:  NewBaseEditor(),
		studio:      studio,
		texModel:    texModel,
		editorModel: model.NewTwoDTextureEditor(),
		viz:         visualization.NewTwoDTexture(studio.api /* FIXME */, studio.GraphicsEngine(), texModel),
	}
}

var _ model.Editor = (*TwoDTextureEditor)(nil)

type TwoDTextureEditor struct {
	BaseEditor
	studio      *Studio
	texModel    *model.TwoDTexture
	editorModel *model.TwoDTextureEditor
	viz         *visualization.TwoDTexture
}

func (e *TwoDTextureEditor) ID() string {
	return e.texModel.Resource().ID()
}

func (e *TwoDTextureEditor) Name() string {
	return e.texModel.Resource().Name()
}

func (e *TwoDTextureEditor) Icon(scope co.Scope) *ui.Image {
	return co.OpenImage(scope, "icons/texture.png")
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
	return co.New(view.TwoDTextureEditor, func() {
		co.WithData(view.TwoDTextureEditorData{
			ResourceModel: e.texModel.Resource(),
			TextureModel:  e.texModel,
			EditorModel:   e.editorModel,
			Visualization: e.viz,
			Controller:    e,
		})
		co.WithLayoutData(layoutData)
	})
}

func (e *TwoDTextureEditor) Destroy() {
	e.viz.Destroy()
}

func (e *TwoDTextureEditor) Dispatch(act interface{}) {
	switch act := act.(type) {
	case action.ChangeResourceName:
		e.changeResourceName(act.Name)
	case action.ChangeTwoDTextureWrapping:
		e.changeWrapping(act.Wrapping)
	case action.ChangeTwoDTextureFiltering:
		e.changeFiltering(act.Filtering)
	case action.ChangeTwoDTextureFormat:
		e.changeFormat(act.Format)
	case action.ChangeTwoDTextureContentFromPath:
		e.changeContentFromPath(act.Path)
	default:
		e.studio.Dispatch(act)
	}
}

func (e *TwoDTextureEditor) changeResourceName(name string) {
	e.changes.Push(change.Name(e.texModel.Resource(),
		change.NameState{
			Value: e.texModel.Resource().Name(),
		},
		change.NameState{
			Value: name,
		},
	))

	// FIXME: Figure out how to avoid this:
	e.studio.NotifyChanged()
}

func (e *TwoDTextureEditor) changeWrapping(wrapping asset.WrapMode) {
	e.changes.Push(change.Wrapping(e.texModel,
		change.WrappingState{
			Value: e.texModel.Wrapping(),
		},
		change.WrappingState{
			Value: wrapping,
		},
	))
}

func (e *TwoDTextureEditor) changeFiltering(filter asset.FilterMode) {
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

func (e *TwoDTextureEditor) changeFormat(format asset.TexelFormat) {
	// TODO
}

func (e *TwoDTextureEditor) changeContentFromPath(path string) {
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

func (e *TwoDTextureEditor) openImage(path string) (image.Image, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image resource: %w", err)
	}
	defer in.Close()

	// TODO: Register image decoders above and ideally move this to
	// a util package.

	img, _, err := image.Decode(in)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}
