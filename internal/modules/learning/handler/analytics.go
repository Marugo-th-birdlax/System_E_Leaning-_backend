package handler

import (
	learningservice "github.com/Marugo/birdlax/internal/modules/learning/service"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	metrics learningservice.MetricsService
}

func NewAnalyticsHandler(m learningservice.MetricsService) *AnalyticsHandler {
	return &AnalyticsHandler{metrics: m}
}

func (h *AnalyticsHandler) GetUserMetric(c *fiber.Ctx) error {
	userID := c.Params("userID")
	courseID := c.Query("course_id", "")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userID required")
	}
	if courseID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "course_id required")
	}
	m, err := h.metrics.GetLearningMetric(userID, courseID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	}
	return c.JSON(m)
}

func (h *AnalyticsHandler) GetCourseOutcome(c *fiber.Ctx) error {
	courseID := c.Params("courseID")
	if courseID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "courseID required")
	}
	o, err := h.metrics.GetCourseOutcome(courseID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	}
	return c.JSON(o)
}
