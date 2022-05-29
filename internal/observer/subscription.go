package observer

// Subscription represents the association of a Callback function to a Target.
type Subscription interface {

	// Delete removes this Subscription from the associated Target and the
	// Callback associated with this Subscription will no longer be called.
	Delete()
}

type subscription struct {
	t    *target
	next *subscription
	f    Filter
	cb   Callback
}

func (s *subscription) Delete() {
	s.t.unsubscribe(s)
}
