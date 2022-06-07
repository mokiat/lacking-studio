package change

import "github.com/mokiat/lacking-studio/internal/studio/history"

type Datable interface {
	SetData([]byte)
}

type BinaryContentState struct {
	Data []byte
}

func BinaryContent(target Datable, from, to BinaryContentState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetData(to.Data)
			return nil
		},
		func() error {
			target.SetData(from.Data)
			return nil
		},
	)
}
