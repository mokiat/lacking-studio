package studio

import (
	"os"

	"github.com/urfave/cli/v2"
)

func Run() error {
	app := &cli.App{
		Name:        "studio",
		Usage:       "run the studio application",
		Description: "Runs the Studio application for the Lacking game engine.",
		Commands: []*cli.Command{
			{
				Name:      "pack",
				Usage:     "Packs the assets of the project",
				Args:      true,
				ArgsUsage: "[project dir] [model name]",
				Action:    runPackApplication,
			},
			{
				Name:      "preview",
				Usage:     "Runs the studio in preview mode",
				Args:      true,
				ArgsUsage: "[project dir]",
				Action:    runPreviewApplication,
			},
		},
	}
	return app.Run(os.Args)
}
