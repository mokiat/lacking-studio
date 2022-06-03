package view

import (
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking-studio/internal/studio/model/action"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/optional"
)

type StudioToolbarData struct {
	StudioModel *model.Studio
}

var StudioToolbar = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		globalCtx   = co.GetContext[global.Context]()
		data        = co.GetData[StudioToolbarData](props)
		studioModel = data.StudioModel
	)

	mvc.UseBinding(studioModel, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeStudioEditorSelection)
	})

	history := studioModel.SelectedHistory()
	mvc.UseBinding(history, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeHistory)
	})

	assetsOverlay := co.UseState(func() co.Overlay {
		return nil
	})

	onAssetsClicked := func() {
		assetsOverlay.Set(co.OpenOverlay(co.New(AssetDialog, func() {
			co.WithData(AssetDialogData{
				Registry: globalCtx.Registry,
			})
			co.WithCallbackData(AssetDialogCallbackData{
				OnAssetSelected: func(id string) {
					mvc.Dispatch(scope, action.OpenResource{
						ID: id,
					})
				},
				OnClose: func() {
					overlay := assetsOverlay.Get()
					overlay.Close()
					assetsOverlay.Set(nil)
				},
			})
		})))
	}

	return co.New(mat.Toolbar, func() {
		co.WithLayoutData(props.LayoutData())

		co.WithChild("assets", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon: co.OpenImage(scope, "icons/assets.png"),
				Text: "Assets",
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: onAssetsClicked,
			})
		}))

		co.WithChild("separator1", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("save", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/save.png"),
				Enabled: optional.Value(history.CanSave()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					mvc.Dispatch(scope, action.Save{})
				},
			})
		}))

		co.WithChild("separator2", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/undo.png"),
				Enabled: optional.Value(history.CanUndo()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					mvc.Dispatch(scope, action.Undo{})
				},
			})
		}))

		co.WithChild("redo", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/redo.png"),
				Enabled: optional.Value(history.CanRedo()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					mvc.Dispatch(scope, action.Redo{})
				},
			})
		}))

		co.WithChild("separator3", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("properties", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon: co.OpenImage(scope, "icons/properties.png"),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					mvc.Dispatch(scope, action.ToggleProperties{})
				},
			})
		}))
	})
})
