package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	learningservice "github.com/Marugo/birdlax/internal/modules/learning/service"
)

type MyHandler struct {
	mySvc learningservice.MyCoursesService
}

func NewMyHandler(s learningservice.MyCoursesService) *MyHandler {
	return &MyHandler{mySvc: s}
}

// GET /api/v1/my/department-courses
func (h *MyHandler) MyDepartmentCourses(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))

	// 👇 ดึง category_id จาก query
	catID := c.Query("category_id", "")
	var catPtr *string
	if catID != "" {
		catPtr = &catID
	}

	items, total, err := h.mySvc.ListDepartmentCourses(
		c.Context(),
		userID,
		catPtr, // 👈 ส่ง pointer ลง service
		page,
		per,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"items":    items,
		"total":    total,
		"page":     page,
		"per_page": per,
	})
}

// GET /api/v1/my/courses
func (h *MyHandler) MyCourses(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))

	items, total, err := h.mySvc.ListMyCourses(c.Context(), userID, page, per)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"items":    items,
		"total":    total,
		"page":     page,
		"per_page": per,
	})
}

func (h *MyHandler) MyCourseProgress(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	courseID := c.Params("courseID")
	if courseID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "courseID required")
	}
	res, err := h.mySvc.GetCourseProgress(c.Context(), userID, courseID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(res)
}
