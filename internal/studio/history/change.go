package history

type Change interface {
	Apply() error
	Revert() error
}
