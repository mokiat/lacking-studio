package action

import "github.com/mokiat/lacking-studio/internal/studio/model"

type OpenResource struct {
	ID string
}

type CloneResource struct {
	Resource *model.Resource
}

type DeleteResource struct {
	Resource *model.Resource
}

type ChangeResourceName struct {
	Resource *model.Resource
	Name     string
}
