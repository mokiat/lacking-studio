package studio

import (
	"path/filepath"

	"github.com/mokiat/lacking/game/asset"
)

func createRegistry(projectDir string) (*asset.Registry, error) {
	storage, err := asset.NewFSStorage(filepath.Join(projectDir, "assets"))
	if err != nil {
		return nil, err
	}
	formatter := asset.NewBlobFormatter()
	return asset.NewRegistry(storage, formatter)
}
