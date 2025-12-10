package handler

import (
	"github.com/Marugo/birdlax/internal/modules/user/models"
	"github.com/Marugo/birdlax/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(r fiber.Router, h *Handler) {
	g := r.Group("")
	// POST
	g.Post("/assets/video", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.UploadVideo)
	g.Post("/lessons", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.CreateLesson)

	// PUT / PATCH / DELETE
	g.Put("/lessons/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.UpdateLesson)
	g.Delete("/lessons/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.DeleteLesson)

	// GET
	g.Get("/assets/:id", h.GetAsset)
	g.Get("/lessons/:id", h.GetLesson)
	g.Get("/lessons", h.ListLessons)
}
