package handler

import (
	"strconv"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc service.Service
}

func New(s service.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) UploadVideo(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file required")
	}
	resp, err := h.svc.UploadVideo(fileHeader)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) CreateLesson(c *fiber.Ctx) error {
	var req dto.CreateLessonReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	l, err := h.svc.CreateLesson(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(l)
}

// GET /v1/lessons/:id
func (h *Handler) GetLesson(c *fiber.Ctx) error {
	id := c.Params("id")
	l, err := h.svc.GetLesson(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lesson not found")
	}
	return c.JSON(l)
}

// GET /v1/lessons?module_id=...&page=1&per_page=20
func (h *Handler) ListLessons(c *fiber.Ctx) error {
	moduleID := c.Query("module_id", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))
	rows, total, err := h.svc.ListLessons(moduleID, page, per)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(fiber.Map{
		"data": rows,
		"meta": fiber.Map{"page": page, "per_page": per, "total": total},
	})
}

// GET /v1/assets/:id (เมตาดาต้า — ไม่ใช่ตัวไฟล์)
func (h *Handler) GetAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	a, err := h.svc.GetAsset(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "asset not found")
	}
	return c.JSON(a)
}
