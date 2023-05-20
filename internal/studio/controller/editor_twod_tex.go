package controller

import (
	"fmt"
	"image"
	"os"

	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	"github.com/mokiat/lacking-studio/internal/studio/model/change"
	"github.com/mokiat/lacking-studio/internal/studio/view"
	"github.com/mokiat/lacking-studio/internal/studio/visualization"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewTwoDTextureEditor(globalCtx global.Context, studio *Studio, editorModel *model.Editor, texModel *model.TwoDTexture) *TwoDTextureEditor {
	return &TwoDTextureEditor{
		studio:    studio,
		histModel: editorModel.History(),
		texModel:  texModel,
		editorModel: model.NewTwoDTextureEditor(
			editorModel,
		),
		viz: visualization.NewTwoDTexture(
			globalCtx.API,
			globalCtx.GraphicsEngine,
			texModel,
		),
	}
}

var _ Editor = (*TwoDTextureEditor)(nil)

type TwoDTextureEditor struct {
	studio      *Studio
	histModel   *model.History
	texModel    *model.TwoDTexture
	editorModel *model.TwoDTextureEditor
	viz         *visualization.TwoDTexture
}

func (e *TwoDTextureEditor) Save() error {
	previewImg := e.viz.TakeSnapshot(ui.Size{
		Width:  model.PreviewSize,
		Height: model.PreviewSize,
	})
	e.texModel.Resource().SetPreviewImage(previewImg)

	if err := e.texModel.Save(); err != nil {
		return fmt.Errorf("error saving texture model %w", err)
	}
	e.histModel.Save()
	return nil
}

func (e *TwoDTextureEditor) Render(scope co.Scope, layoutData any) co.Instance {
	return co.New(view.TwoDTextureEditor, func() {
		co.WithData(view.TwoDTextureEditorData{
			ResourceModel:    e.texModel.Resource(),
			TextureModel:     e.texModel,
			EditorModel:      e.editorModel,
			Visualization:    e.viz,
			StudioController: e.studio,
			EditorController: e,
		})
		co.WithLayoutData(layoutData)
		co.WithScope(mvc.UseReducer(scope, e))
	})
}

func (e *TwoDTextureEditor) Destroy() {
	e.viz.Destroy()
}

func (e *TwoDTextureEditor) Reduce(act mvc.Action) bool {
	switch act := act.(type) {
	case action.ChangeTwoDTextureWrapping:
		e.changeWrapping(act.Wrapping)
		return true
	case action.ChangeTwoDTextureFiltering:
		e.changeFiltering(act.Filtering)
		return true
	case action.ChangeTwoDTextureFormat:
		e.changeFormat(act.Format)
		return true
	case action.ChangeTwoDTextureContentFromPath:
		e.changeContentFromPath(act.Path)
		return true
	default:
		return false
	}
}

func (e *TwoDTextureEditor) OnRenameResource(name string) {
	e.histModel.Add(change.Name(e.texModel.Resource(),
		change.NameState{
			Value: e.texModel.Resource().Name(),
		},
		change.NameState{
			Value: name,
		},
	))
}

func (e *TwoDTextureEditor) changeWrapping(wrapping asset.WrapMode) {
	e.histModel.Add(change.Wrapping(e.texModel,
		change.WrappingState{
			Value: e.texModel.Wrapping(),
		},
		change.WrappingState{
			Value: wrapping,
		},
	))
}

func (e *TwoDTextureEditor) changeFiltering(filter asset.FilterMode) {
	e.histModel.Add(change.Filtering(
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
		panic(fmt.Errorf("failed to open source image: %w", err))
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
	if err := e.histModel.Add(ch); err != nil {
		panic(fmt.Errorf("failed to apply change: %w", err))
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
