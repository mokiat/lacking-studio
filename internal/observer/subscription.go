package observer

// Subscription represents the association of a Callback function to a Target.
type Subscription struct {
	target   *Target
	next     *Subscription
	callback Callback
}

// Delete removes this Subscription.
// The Callback associated with this Subscription will no longer be called.
func (s *Subscription) Delete() {
	s.target.Unsubscribe(s)
}
