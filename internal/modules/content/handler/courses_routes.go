package handler

import "github.com/gofiber/fiber/v2"

func RegisterCourseRoutes(r fiber.Router, h *CourseHandler) {
	g := r.Group("")

	// Courses (Admin CRUD + Learner list/get)
	g.Get("/courses", h.ListCourses)
	g.Get("/courses/:id", h.GetCourse)
	g.Post("/courses", h.CreateCourse)
	g.Put("/courses/:id", h.UpdateCourse)
	g.Delete("/courses/:id", h.DeleteCourse)

	// Modules
	g.Get("/courses/:id/modules", h.ListModules) // by course
	g.Post("/courses/:id/modules", h.CreateModule)

	g.Get("/modules/:id/lessons", h.ListLessonsOfModule)
	g.Put("/modules/:id", h.UpdateModule)
	g.Delete("/modules/:id", h.DeleteModule)
}
