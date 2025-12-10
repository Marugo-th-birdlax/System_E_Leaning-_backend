package handler

import (
	"github.com/Marugo/birdlax/internal/modules/user"
	"github.com/Marugo/birdlax/internal/modules/user/models"
	"github.com/Marugo/birdlax/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(r fiber.Router, svc user.Service) {
	h := NewHTTPHandler(svc)

	// ===== Users CRUD (admin, hr เท่านั้น) =====
	g := r.Group("/users", middleware.RequireRoles(models.RoleAdmin, models.RoleHR))
	g.Get("/", h.List)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)

	// User department roles
	g.Get("/:id/department-roles", h.ListUserDepartmentRoles)
	g.Post("/:id/department-roles", h.AddUserDepartmentRole)
	g.Delete("/:id/department-roles/:relID", h.DeleteUserDepartmentRole)

	// ===== Departments CRUD =====
	dg := r.Group("/departments", middleware.RequireRoles(models.RoleAdmin, models.RoleHR))
	dg.Get("/", h.ListDepartments)
	dg.Post("/", h.CreateDepartment)
	dg.Get("/:id", h.GetDepartment)
	dg.Put("/:id", h.UpdateDepartment)
	dg.Delete("/:id", h.DeleteDepartment)
	dg.Get("/:id/managers", h.ListDepartmentManagers)
}
