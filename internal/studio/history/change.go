package history

type Change interface {
	Apply() error
	Revert() error
}

type CombinedChange struct {
	Changes []Change
}

func (ch *CombinedChange) Apply() error {
	for i := 0; i < len(ch.Changes); i++ {
		if err := ch.Changes[i].Apply(); err != nil {
			return err
		}
	}
	return nil
}

func (ch *CombinedChange) Revert() error {
	for i := len(ch.Changes) - 1; i >= 0; i-- {
		if err := ch.Changes[i].Revert(); err != nil {
			return err
		}
	}
	return nil
}

func FuncChange(apply, revert func() error) Change {
	return &funcChange{
		apply:  apply,
		revert: revert,
	}
}

type funcChange struct {
	apply  func() error
	revert func() error
}

func (ch *funcChange) Apply() error {
	return ch.apply()
}

func (ch *funcChange) Revert() error {
	return ch.revert()
}
