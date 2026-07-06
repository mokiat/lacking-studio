package studio

import (
	"cmp"
	"fmt"
	"path/filepath"

	"github.com/mokiat/gog/filter"
	"github.com/mokiat/lacking/core/resource"
	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/urfave/cli/v2"
)

func runPackApplication(ctx *cli.Context) error {
	projectDir := cmp.Or(ctx.Args().First(), ".")

	pathFilter := filter.True[string]()
	if modelName := ctx.Args().Get(1); modelName != "" {
		pathFilter = filter.Equal(modelName)
	}

	store, err := resource.NewFileStore(filepath.Join(projectDir, "assets"))
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}
	if err := dsl.Run(store, pathFilter); err != nil {
		return fmt.Errorf("error running DSL: %w", err)
	}

	return nil
}
