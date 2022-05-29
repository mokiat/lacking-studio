package observer

import "github.com/mokiat/lacking/util/filter"

// Target represents something that can be observed.
type Target interface {

	// Subscribe registers the specified callback function to this target.
	// It will only be called if the given change passes the specified filters.
	// The returned Subscription can be used to unsubscribe from this Target.
	Subscribe(callback Callback, filters ...Filter) Subscription

	// SignalChange sends the specified change to all Subscribers.
	SignalChange(change Change)

	// AccumulateChanges calls the specified closure function and any calls
	// to SignalChange during that invocation will not be sent but instead will
	// be recorded. Once the closure is complete then all recorded changes will
	// be sent via a single MultiChange change.
	AccumulateChanges(fn func() error) error
}

// WireTargets creates a new coupling between the two Targets. If a change
// occurs in the first target and passes the specified filters, it will be
// reported as a change by the second target as well.
func WireTargets(parent, child Target, fltrs ...Filter) Subscription {
	return parent.Subscribe(func(change Change) {
		child.SignalChange(change)
	}, fltrs...)
}

// NewTarget creates a new Target instance.
func NewTarget() Target {
	return &target{}
}

type target struct {
	firstSubscription *subscription
	accumulationDepth int
	ongoingChanges    []Change
}

func (t *target) Subscribe(callback Callback, filters ...Filter) Subscription {
	var fltr Filter
	if len(filters) == 0 {
		fltr = filter.Always[Change]()
	} else {
		fltr = filter.All[Change](filters...)
	}
	subscription := &subscription{
		t:    t,
		next: t.firstSubscription,
		f:    fltr,
		cb:   callback,
	}
	t.firstSubscription = subscription
	return subscription
}

func (t *target) SignalChange(change Change) {
	if t.accumulationDepth > 0 {
		t.ongoingChanges = append(t.ongoingChanges, change)
		return
	}
	current := t.firstSubscription
	for current != nil {
		if current.f(change) {
			current.cb(change)
		}
		current = current.next
	}
}

func (t *target) AccumulateChanges(fn func() error) error {
	var err error
	defer func() {
		if t.accumulationDepth == 0 {
			if err == nil {
				t.SignalChange(MultiChange{
					Changes: t.ongoingChanges,
				})
			}
			t.ongoingChanges = t.ongoingChanges[:0]
		}
	}()
	t.accumulationDepth++
	defer func() {
		t.accumulationDepth--
	}()
	err = fn()
	return err
}

func (t *target) unsubscribe(subscription *subscription) {
	if t.firstSubscription == subscription {
		t.firstSubscription = subscription.next
		return
	}
	current := t.firstSubscription
	for current != nil {
		if current.next == subscription {
			current.next = subscription.next
			return
		}
		current = current.next
	}
}
