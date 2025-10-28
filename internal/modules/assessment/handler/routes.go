package handler

import "github.com/gofiber/fiber/v2"

func Register(r fiber.Router, h *Handler) {
	g := r.Group("/assessments")
	g.Get("/", h.ListAssessments)
	g.Get("/:id", h.GetAssessment)
	g.Post("/", h.CreateAssessment)
	g.Post("/:id/questions", h.AddQuestion)
	g.Delete("/:id", h.DeleteAssessment)
	g.Put("/:id/questions/:qid", h.UpdateQuestion)
	g.Put("/:id", h.UpdateAssessment)
	g.Delete("/:id/questions/:qid", h.DeleteQuestion)

}

// NEW: register attempts
func RegisterAttemptRoutes(r fiber.Router, h *AttemptHandler) {
	g := r.Group("/assessments")
	g.Post("/:id/attempts", h.Start)

	a := r.Group("/attempts")
	a.Get("/:id", h.Get)
	a.Post("/:id/answers", h.UpsertAnswer)
	a.Post("/:id/submit", h.Submit)
}
