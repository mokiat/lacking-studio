package change

import (
	"github.com/mokiat/lacking-studio/internal/studio/history"
	"github.com/mokiat/lacking-studio/internal/studio/model"
)

var _ history.Change = (*CubeTextureChangeSource)(nil)

type CubeTextureChangeSource struct {
	Controller model.CubeTextureEditor
	FromURI    string
	ToURI      string
}

func (ch *CubeTextureChangeSource) Apply() error {
	ch.Controller.SetSourcePath(ch.ToURI)
	return ch.Controller.ReloadSource()
}

func (ch *CubeTextureChangeSource) Revert() error {
	ch.Controller.SetSourcePath(ch.FromURI)
	return ch.Controller.ReloadSource()
}
