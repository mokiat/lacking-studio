package common

import (
	"github.com/mokiat/lacking-studio/internal/widget"
	co "github.com/mokiat/lacking/ui/component"
)

func OpenWarning(scope co.Scope, message string) {
	co.OpenOverlay(scope, co.New(widget.NotificationModal, func() {
		co.WithData(widget.NotificationModalData{
			Icon: co.OpenImage(scope, "icons/warning.png"),
			Text: message,
		})
	}))
}

func OpenError(scope co.Scope, message string) {
	co.OpenOverlay(scope, co.New(widget.NotificationModal, func() {
		co.WithData(widget.NotificationModalData{
			Icon: co.OpenImage(scope, "icons/error.png"),
			Text: message,
		})
	}))
}

func OpenConfirmation(scope co.Scope, message string, cb func()) {
	co.OpenOverlay(scope, co.New(widget.ConfirmationModal, func() {
		co.WithData(widget.ConfirmationModalData{
			Icon: co.OpenImage(scope, "icons/warning.png"),
			Text: message,
		})
		co.WithCallbackData(widget.ConfirmationModalCallbackData{
			OnApply: cb,
		})
	}))
}

func OpenLoading(scope co.Scope) co.Overlay {
	return co.OpenOverlay(scope, co.New(widget.LoadingModal, nil))
}
