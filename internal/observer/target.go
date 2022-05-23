package observer

// NewTarget creates a new Target instance.
func NewTarget() *Target {
	return &Target{}
}

// Target represents something that can be observed
type Target struct {
	firstSubscription *Subscription
	accumulationDepth int
	ongoingChanges    []Change
}

// Subscribe registers the specified callback to this Target and any changes
// happening to it would be reported through the callback.
// The returned Subscription can be used to unsubscribe from the Target.
func (t *Target) Subscribe(callback Callback) *Subscription {
	subscription := &Subscription{
		target:   t,
		next:     t.firstSubscription,
		callback: callback,
	}
	t.firstSubscription = subscription
	return subscription
}

// Unsubscribe unregisters the specified Subscription from this Target.
// An alternative is to call the Delete method on the Subscription object.
func (t *Target) Unsubscribe(subscription *Subscription) {
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

// SignalChange sends the specified change to all Subscribers.
func (t *Target) SignalChange(change Change) {
	if t.accumulationDepth > 0 {
		t.ongoingChanges = append(t.ongoingChanges, change)
		return
	}
	current := t.firstSubscription
	for current != nil {
		current.callback(change)
		current = current.next
	}
}

// AccumulateChanges calls the specified closure function and any calls
// to SignalChange during that invocation will not be sent but instead will
// be recorded. Once the closure is complete then all recorded changes will
// be sent via a single MultiChange change.
func (t *Target) AccumulateChanges(fn func() error) error {
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

func WireTargets(parent, child *Target) *Subscription {
	return parent.Subscribe(func(change Change) {
		child.SignalChange(change)
	})
}
