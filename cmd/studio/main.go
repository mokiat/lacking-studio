package main

import (
	"os"

	"github.com/mokiat/lacking-studio/studio"
	"github.com/mokiat/lacking/debug/log"
)

// TODO: Get rid of this. Apps are expected to package the studio
// in their own cmd folder.

func main() {
	if err := studio.Run(); err != nil {
		log.Error("Error: %v", err)
		os.Exit(1)
	}
}
