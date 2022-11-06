package view

import (
	"github.com/mokiat/gog/filter"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
)

type StudioToolbarData struct {
	StudioModel      *model.Studio
	StudioController StudioController
}

var StudioToolbar = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data        = co.GetData[StudioToolbarData](props)
		studioModel = data.StudioModel
		controller  = data.StudioController
	)

	mvc.UseBinding(studioModel, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeStudioEditorSelection)
	})

	history := studioModel.SelectedHistory()
	mvc.UseBinding(history, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeHistory)
	})

	mvc.UseBinding(studioModel.Registry(), filter.True[mvc.Change]())

	assetsOverlay := co.UseState(func() co.Overlay {
		return nil
	})

	onAssetsClicked := func() {
		assetsOverlay.Set(co.OpenOverlay(co.New(AssetDialog, func() {
			co.WithData(AssetDialogData{
				Registry:   studioModel.Registry(),
				Controller: controller,
			})
			co.WithCallbackData(AssetDialogCallbackData{
				OnOpen: func(id string) {
					controller.OnOpenResource(id)
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
				Enabled: opt.V(history.CanSave()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					controller.OnSave()
				},
			})
		}))

		co.WithChild("separator2", co.New(mat.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/undo.png"),
				Enabled: opt.V(history.CanUndo()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					controller.OnUndo()
				},
			})
		}))

		co.WithChild("redo", co.New(mat.ToolbarButton, func() {
			co.WithData(mat.ToolbarButtonData{
				Icon:    co.OpenImage(scope, "icons/redo.png"),
				Enabled: opt.V(history.CanRedo()),
			})
			co.WithCallbackData(mat.ToolbarButtonCallbackData{
				OnClick: func() {
					controller.OnRedo()
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
					controller.OnToggleProperties()
				},
			})
		}))
	})
})
