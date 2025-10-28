package handler

import "github.com/gofiber/fiber/v2"

func Register(r fiber.Router, h *Handler) {
	g := r.Group("")
	// POST
	g.Post("/assets/video", h.UploadVideo)
	g.Post("/lessons", h.CreateLesson)
	// GET
	g.Get("/assets/:id", h.GetAsset)
	g.Get("/lessons/:id", h.GetLesson)
	g.Get("/lessons", h.ListLessons)
}
