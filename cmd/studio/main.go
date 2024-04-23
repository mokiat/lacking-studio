package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "studio",
		Usage:       "run the studio application",
		Description: "Runs the Studio application for the Lacking game engine.",
		Commands: []*cli.Command{
			{
				Name:      "preview",
				Usage:     "Runs the studio in preview mode",
				Args:      true,
				ArgsUsage: "[project dir]",
				Action:    runPreviewApplication,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
