package global

import co "github.com/mokiat/lacking/ui/component"

func Reducer() (co.Reducer, interface{}) {
	initialStore := Store{}

	reducer := func(store *co.Store, action interface{}) interface{} {
		var value Store
		store.Inject(&value)

		// switch action := action.(type) {
		// case ChangeViewAction:
		// 	value.MainViewIndex = action.ViewIndex
		// }
		return value
	}

	return reducer, initialStore
}

type Store struct {
}

func EditorStoreReducer() (co.Reducer, interface{}) {
	initialStore := Store{}

	reducer := func(store *co.Store, action interface{}) interface{} {
		var value Store
		store.Inject(&value)

		// switch action := action.(type) {
		// case ChangeViewAction:
		// 	value.MainViewIndex = action.ViewIndex
		// }
		return value
	}

	return reducer, initialStore
}

type EditorStore struct {
	PropertiesPaneVisible bool
}
