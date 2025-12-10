package handler

import (
	"github.com/Marugo/birdlax/internal/modules/user/models"
	"github.com/Marugo/birdlax/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCourseRoutes(r fiber.Router, h *CourseHandler) {
	g := r.Group("")

	// Courses (Admin CRUD + Learner list/get)
	g.Get("/courses", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.ListCourses)
	g.Get("/courses/:id", h.GetCourse)
	g.Post("/courses", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.CreateCourse)
	g.Put("/courses/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.UpdateCourse)
	g.Delete("/courses/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.DeleteCourse)

	// Modules
	g.Get("/courses/:id/modules", h.ListModules) // by course
	g.Post("/courses/:id/modules", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.CreateModule)

	g.Get("/modules/:id/lessons", h.ListLessonsOfModule)
	g.Put("/modules/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.UpdateModule)
	g.Delete("/modules/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleHR), h.DeleteModule)
}
