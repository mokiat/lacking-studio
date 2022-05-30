package view

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	co "github.com/mokiat/lacking/ui/component"
)

// TODO: Move to ui/mvc when stable enough.

func WithBinding(target observer.Target, filter observer.Filter) {
	lifecycle := co.UseLifecycle(func(handle co.LifecycleHandle) *bindingLifecycle {
		return &bindingLifecycle{
			target: target,
			filter: filter,
			handle: handle,
		}
	})
	lifecycle.newTarget = target
}

type bindingLifecycle struct {
	co.BaseLifecycle
	target       observer.Target
	newTarget    observer.Target
	filter       observer.Filter
	handle       co.LifecycleHandle
	subscription observer.Subscription
}

func (l *bindingLifecycle) OnCreate(props co.Properties, scope co.Scope) {
	l.subscription = l.target.Subscribe(func(change observer.Change) {
		l.handle.NotifyChanged()
	}, l.filter)
}

func (l *bindingLifecycle) OnUpdate(props co.Properties, scope co.Scope) {
	if l.newTarget != l.target {
		l.OnDestroy(scope)
		l.target = l.newTarget
		l.newTarget = nil
		l.OnCreate(props, scope)
	}
}

func (l *bindingLifecycle) OnDestroy(scope co.Scope) {
	l.subscription.Delete()
}
