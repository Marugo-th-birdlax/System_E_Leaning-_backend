package handler

import (
	"github.com/Marugo/birdlax/internal/modules/user/models"
	"github.com/Marugo/birdlax/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCategoryRoutes(r fiber.Router, h *CategoryHandler) {
	// READ
	r.Get("/categories", h.ListCategories)
	r.Get("/categories/:id", h.GetCategory)
	r.Get("/categories/:id/courses", h.ListCoursesOfCategory)

	// CREATE/UPDATE/DELETE
	r.Post("/categories", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.CreateCategory)
	r.Patch("/categories/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.UpdateCategory)
	// (ถ้าชอบ PUT ก็ทำเพิ่มได้)
	r.Delete("/categories/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.DeleteCategory)
}
