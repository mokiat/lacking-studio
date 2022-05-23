package view

import (
	"fmt"

	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking-studio/internal/studio/global"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
)

type notifLifecycle struct {
	co.BaseLifecycle
	target       *observer.Target
	filter       func(change observer.Change) bool
	handle       co.LifecycleHandle
	subscription *observer.Subscription
}

func (l *notifLifecycle) OnCreate(props co.Properties) {
	l.subscription = l.target.Subscribe(func(change observer.Change) {
		if l.filter(change) {
			l.handle.NotifyChanged()
		}
	})
}

func (l *notifLifecycle) OnDestroy() {
	l.subscription.Delete()
}

func WithNotifications(target *observer.Target, filter func(change observer.Change) bool) {
	co.UseLifecycle(func(handle co.LifecycleHandle) *notifLifecycle {
		return &notifLifecycle{
			target: target,
			filter: filter,
			handle: handle,
		}
	})
}

var TwoDTexture = co.Define(func(props co.Properties) co.Instance {
	editor := props.Data().(model.TwoDTextureEditor)

	WithNotifications(editor.Target(), func(change observer.Change) bool {
		fmt.Println("CHANGE:", change.Description())
		return true // TODO
	})

	viz := editor.Visualization()

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(mat.SurfaceColor),
			Layout:          mat.NewFrameLayout(),
		})
		co.WithLayoutData(props.LayoutData())

		co.WithChild("center", co.New(mat.DropZone, func() {
			co.WithCallbackData(mat.DropZoneCallbackData{
				OnDrop: func(paths []string) bool {
					editor.ChangeContent(paths[0])
					return true
				},
			})
			co.WithLayoutData(mat.LayoutData{
				Alignment: mat.AlignmentCenter,
			})

			co.WithChild("viewport", co.New(mat.Viewport, func() {
				co.WithData(mat.ViewportData{
					API: co.GetContext[global.Context]().API,
				})
				co.WithCallbackData(mat.ViewportCallbackData{
					OnMouseEvent: viz.OnViewportMouseEvent,
					OnRender:     viz.OnViewportRender,
				})
			}))
		}))

		if editor.IsPropertiesVisible() {
			co.WithChild("right", co.New(TwoDTextureProperties, func() {
				co.WithData(editor)
				co.WithLayoutData(mat.LayoutData{
					Alignment: mat.AlignmentRight,
					Width:     optional.Value(500),
				})
			}))
		}
	})
})
