package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
)

var _ history.Change = (*Combined)(nil)

type Combined struct {
	Changes []history.Change
}

func (ch *Combined) Apply() error {
	for i := 0; i < len(ch.Changes); i++ {
		if err := ch.Changes[i].Apply(); err != nil {
			return err
		}
	}
	return nil
}

func (ch *Combined) Revert() error {
	for i := len(ch.Changes) - 1; i >= 0; i-- {
		if err := ch.Changes[i].Apply(); err != nil {
			return err
		}
	}
	return nil
}
