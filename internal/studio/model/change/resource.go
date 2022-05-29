package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
)

type Nameable interface {
	SetName(name string)
}

type NameState struct {
	Value string
}

func Name(target Nameable, from, to NameState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetName(to.Value)
			return nil
		},
		func() error {
			target.SetName(from.Value)
			return nil
		},
	)
}
