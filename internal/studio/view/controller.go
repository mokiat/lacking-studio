package view

// TODO: Consider dispatching actions and having controllers handle them in
// depth. E.g. the Clone/Delete resource actions could bubble up all the way
// to the StudioController.
// Furthermore, a controller could be made up of multiple handlers that can
// try and handle the action so that not all code is in a single place.
// If this idea does not work, an alternative could be to dispatch the action
// down the chain of models but that might not work well with the idea above.

type Controller interface {
	Dispatch(action interface{})
}
