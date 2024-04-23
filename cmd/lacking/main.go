package main

import (
	"flag"

	"github.com/mokiat/lacking-studio/studio"
)

// Deprecated: The approach is to have each project run the studio on its own.

func main() {
	flag.Parse()
	projectDir := "."
	if flag.NArg() > 0 {
		projectDir = flag.Arg(0)
	}
	studio.Run(studio.WithProjectDir(projectDir))
}
