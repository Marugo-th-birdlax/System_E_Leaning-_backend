package handler

import (
	"github.com/Marugo/birdlax/internal/modules/learning/dto"
	"github.com/Marugo/birdlax/internal/modules/learning/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{ svc service.Service }

func New(s service.Service) *Handler { return &Handler{svc: s} }

func userID(c *fiber.Ctx) string {
	if v := c.Locals("user_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// POST /v1/courses/:courseID/enroll
func (h *Handler) EnrollCourse(c *fiber.Ctx) error {
	uid := userID(c)
	courseID := c.Params("courseID")
	e, err := h.svc.EnrollCourse(uid, courseID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": e.ID, "user_id": e.UserID, "course_id": e.CourseID, "status": e.Status, "progress_percent": e.ProgressPercent,
	})
}

// GET /v1/enrollments/:courseID
func (h *Handler) GetEnrollment(c *fiber.Ctx) error {
	uid := userID(c)
	courseID := c.Params("courseID")
	e, err := h.svc.GetEnrollment(uid, courseID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "enrollment not found")
	}
	return c.JSON(fiber.Map{
		"id": e.ID, "user_id": e.UserID, "course_id": e.CourseID, "status": e.Status, "progress_percent": e.ProgressPercent,
	})
}

// POST /v1/lessons/:lessonID/start
func (h *Handler) StartLesson(c *fiber.Ctx) error {
	uid := userID(c)
	lessonID := c.Params("lessonID")
	var req dto.StartLessonReq
	_ = c.BodyParser(&req)
	p, err := h.svc.StartLesson(uid, lessonID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

// POST /v1/lessons/:lessonID/track
func (h *Handler) TrackLesson(c *fiber.Ctx) error {
	uid := userID(c)
	lessonID := c.Params("lessonID")
	var req dto.TrackLessonReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	p, err := h.svc.TrackLesson(uid, lessonID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(p)
}

// POST /v1/courses/:courseID/lessons/:lessonID/complete
func (h *Handler) CompleteLesson(c *fiber.Ctx) error {
	uid := userID(c)
	lessonID := c.Params("lessonID")
	courseID := c.Params("courseID")

	var req dto.CompleteLessonReq
	_ = c.BodyParser(&req)

	p, err := h.svc.CompleteLesson(uid, lessonID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// อัปเดตเปอร์เซ็นต์ของคอร์สทันที
	if err := h.svc.UpdateEnrollmentPercent(uid, courseID); err != nil {
		// ไม่ให้ fail ทั้งคำขอ: ส่งผลลัพธ์ progress บทเรียนไปก่อน แล้วแนบ warning
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"lesson_progress": p,
			"warning":         "course percent not updated: " + err.Error(),
		})
	}

	// ดึง enrollment เพื่อตอบ % ปัจจุบันกลับไป
	e, _ := h.svc.GetEnrollment(uid, courseID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"lesson_progress": p,
		"enrollment": fiber.Map{
			"id":               e.ID,
			"user_id":          e.UserID,
			"course_id":        e.CourseID,
			"status":           e.Status,
			"progress_percent": e.ProgressPercent,
		},
	})
}
