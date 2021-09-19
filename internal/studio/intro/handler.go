package intro

// import (
// 	"fmt"
// 	"log"

// 	"github.com/mokiat/lacking/ui"
// )

// type Config struct{}

// func (c Config) SetupView(view *ui.View) error {
// 	template, err := view.Context().OpenTemplate("resources/studio/intro/view.xml")
// 	if err != nil {
// 		return fmt.Errorf("failed to open template: %w", err)
// 	}
// 	rootControl, err := view.Context().InstantiateTemplate(template, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to instantiate template: %w", err)
// 	}
// 	view.SetRoot(rootControl)
// 	view.SetHandler(&Handler{})
// 	return nil
// }

// type Handler struct{}

// func (h *Handler) OnCreate(view *ui.View) {
// 	log.Println("INTRO CREATE")
// }

// func (h *Handler) OnShow(view *ui.View) {
// 	log.Println("INTRO SHOW")
// }

// func (h *Handler) OnHide(view *ui.View) {
// 	log.Println("INTRO HIDE")
// }

// func (h *Handler) OnDestroy(view *ui.View) {
// 	log.Println("INTRO DESTROY")
// }
