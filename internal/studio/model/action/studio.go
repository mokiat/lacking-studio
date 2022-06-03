package action

import "github.com/mokiat/lacking-studio/internal/studio/model"

type ChangeSelectedEditor struct {
	Editor *model.Editor
}

type CloseEditor struct {
	Editor *model.Editor
}

type Undo struct{}

type Redo struct{}

type Save struct{}

type ToggleProperties struct{}
