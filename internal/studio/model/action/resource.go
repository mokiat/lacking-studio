package action

import "github.com/mokiat/lacking-studio/internal/studio/model"

type ChangeResourceName struct {
	Resource *model.Resource
	Name     string
}

type DeleteResource struct {
	Resource *model.Resource
}

type CloneResource struct {
	Resource *model.Resource
}
