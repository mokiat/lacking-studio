package studio

import (
	"cmp"
	"fmt"
	"path/filepath"

	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/mokiat/lacking/game/chunked"
	"github.com/urfave/cli/v2"
)

func runPackApplication(ctx *cli.Context) error {
	projectDir := cmp.Or(ctx.Args().First(), ".")

	var modelNames []string
	if modelName := ctx.Args().Get(1); modelName != "" {
		modelNames = append(modelNames, modelName)
	}

	storage, err := chunked.NewFileStorage(filepath.Join(projectDir, "assets"))
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}
	if err := dsl.Run(storage, modelNames); err != nil {
		return fmt.Errorf("error running DSL: %w", err)
	}

	return nil
}
