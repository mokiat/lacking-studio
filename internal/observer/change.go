package observer

type Change any

type MultiChange struct {
	Changes []Change
}
