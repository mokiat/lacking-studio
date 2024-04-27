package studio

import (
	"cmp"
	"fmt"

	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/urfave/cli/v2"
)

func runPackApplication(ctx *cli.Context) error {
	projectDir := cmp.Or(ctx.Args().First(), ".")

	var modelNames []string
	if modelName := ctx.Args().Get(1); modelName != "" {
		modelNames = append(modelNames, modelName)
	}

	registry, err := createRegistry(projectDir)
	if err != nil {
		return fmt.Errorf("error creating registry: %w", err)
	}
	if err := dsl.Run(registry, modelNames); err != nil {
		return fmt.Errorf("error running DSL: %w", err)
	}

	return nil
}
