package handler

import (
	"github.com/gofiber/fiber/v2"
)

func Register(r fiber.Router, h *Handler) {
	g := r.Group("")

	g.Post("/courses/:courseID/enroll", h.EnrollCourse)
	g.Get("/enrollments/:courseID", h.GetEnrollment)

	g.Post("/lessons/:lessonID/start", h.StartLesson)
	g.Post("/lessons/:lessonID/track", h.TrackLesson)
	g.Post("/courses/:courseID/lessons/:lessonID/complete", h.CompleteLesson)
}

// เรียกเพิ่มใน app.Register หลังเรียก Register(...)
func RegisterAdminRoutes(r fiber.Router, analyticsHandler *AnalyticsHandler) {
	admin := r.Group("/analytics")
	admin.Get("/users/:userID/metrics", analyticsHandler.GetUserMetric)
	admin.Get("/courses/:courseID/outcome", analyticsHandler.GetCourseOutcome)
}
