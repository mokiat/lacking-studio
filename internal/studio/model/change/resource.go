package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
	"github.com/mokiat/lacking/game/asset"
)

type ResourceNameState struct {
	Value string
}

func ResourceName(target model.Resource, from, to ResourceNameState) history.Change {
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

type WrappingState struct {
	Value asset.WrapMode
}

func Wrapping(target model.Wrappable, from, to WrappingState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetWrapping(to.Value)
			return nil
		},
		func() error {
			target.SetWrapping(from.Value)
			return nil
		},
	)
}

type FilteringState struct {
	Value asset.FilterMode
}

func Filtering(target model.Filterable, from, to FilteringState) history.Change {
	return history.FuncChange(
		func() error {
			target.SetFiltering(to.Value)
			return nil
		},
		func() error {
			target.SetFiltering(from.Value)
			return nil
		},
	)
}
