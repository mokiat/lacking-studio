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
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewCubeTextureEditor(globalCtx global.Context, studio *Studio, editorModel *model.Editor, texModel *model.CubeTexture) *CubeTextureEditor {
	return &CubeTextureEditor{
		studio:    studio,
		histModel: editorModel.History(),
		texModel:  texModel,
		editorModel: model.NewCubeTextureEditor(
			editorModel,
		),
		viz: visualization.NewCubeTexture(
			globalCtx.API,
			globalCtx.GraphicsEngine,
			texModel,
		),
	}
}

var _ Editor = (*CubeTextureEditor)(nil)

type CubeTextureEditor struct {
	studio      *Studio
	histModel   *model.History
	texModel    *model.CubeTexture
	editorModel *model.CubeTextureEditor
	viz         *visualization.CubeTexture
}

func (e *CubeTextureEditor) Save() error {
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

func (e *CubeTextureEditor) Render(scope co.Scope, layoutData mat.LayoutData) co.Instance {
	return co.New(view.CubeTextureEditor, func() {
		co.WithData(view.CubeTextureEditorData{
			ResourceModel:    e.texModel.Resource(),
			TextureModel:     e.texModel,
			EditorModel:      e.editorModel,
			StudioController: e.studio,
			EditorController: e,
			Visualization:    e.viz,
		})
		co.WithLayoutData(layoutData)
		co.WithScope(mvc.UseReducer(scope, e))
	})
}

func (e *CubeTextureEditor) Destroy() {
	e.viz.Destroy()
}

func (e *CubeTextureEditor) Reduce(act mvc.Action) bool {
	switch act := act.(type) {
	case action.ChangeCubeTextureFiltering:
		e.changeFiltering(act.Filtering)
		return true
	case action.ChangeCubeTextureFormat:
		e.changeFormat(act.Format)
		return true
	case action.ChangeCubeTextureContentFromPath:
		e.changeContentFromPath(act.Path)
		return true
	default:
		return false
	}
}

func (e *CubeTextureEditor) OnRenameResource(name string) {
	e.histModel.Add(change.Name(e.texModel.Resource(),
		change.NameState{
			Value: e.texModel.Resource().Name(),
		},
		change.NameState{
			Value: name,
		},
	))
}

func (e *CubeTextureEditor) changeFiltering(filter asset.FilterMode) {
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

func (e *CubeTextureEditor) changeFormat(format asset.TexelFormat) {
	// TODO
}

func (e *CubeTextureEditor) changeContentFromPath(path string) {
	img, err := e.openImage(path)
	if err != nil {
		panic(fmt.Errorf("failed to open source image: %w", err))
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
		panic(fmt.Errorf("failed to build cube image: %w", err))
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
	if err := e.histModel.Add(ch); err != nil {
		panic(fmt.Errorf("failed to apply change: %w", err)) // TODO
	}
}

func (e *CubeTextureEditor) openImage(path string) (image.Image, error) {
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
