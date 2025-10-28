package handler

import "github.com/gofiber/fiber/v2"

func Register(r fiber.Router, h *Handler) {
	g := r.Group("")

	g.Post("/courses/:courseID/enroll", h.EnrollCourse)
	g.Get("/enrollments/:courseID", h.GetEnrollment)

	g.Post("/lessons/:lessonID/start", h.StartLesson)
	g.Post("/lessons/:lessonID/track", h.TrackLesson)
	g.Post("/courses/:courseID/lessons/:lessonID/complete", h.CompleteLesson)
}
