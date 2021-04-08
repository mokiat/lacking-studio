package studio

import (
	"log"

	"github.com/mokiat/lacking/ui"
)

type Handler struct{}

func (h *Handler) OnCreate(view ui.View) {
	log.Println("STUDIO CREATE")
}

func (h *Handler) OnShow(view ui.View) {
	log.Println("STUDIO SHOW")
}

func (h *Handler) OnHide(view ui.View) {
	log.Println("STUDIO HIDE")
}

func (h *Handler) OnDestroy(view ui.View) {
	log.Println("STUDIO DESTROY")
}
