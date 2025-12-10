package handler

import (
	"strconv"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/service"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler { return &CategoryHandler{svc: s} }

// ====== CRUD ======

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	obj, err := h.svc.CreateCategory(req)
	if err != nil {
		if err.Error() == "code and title required" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": obj})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req dto.UpdateCategoryReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	obj, err := h.svc.UpdateCategory(c.Params("id"), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(fiber.Map{"data": obj})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.svc.DeleteCategory(c.Params("id")); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ====== READ ======

func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	q := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))

	rows, total, err := h.svc.ListCategories(q, page, per)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(fiber.Map{
		"data": rows,
		"meta": fiber.Map{"page": page, "per_page": per, "total": total},
	})
}

func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	x, err := h.svc.GetCategory(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "category not found")
	}
	return c.JSON(fiber.Map{"data": x})
}

func (h *CategoryHandler) ListCoursesOfCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	q := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))

	rows, total, err := h.svc.ListCoursesOfCategory(id, q, page, per)
	if err != nil {
		if err.Error() == "category not found" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(fiber.Map{
		"data": rows,
		"meta": fiber.Map{"page": page, "per_page": per, "total": total},
	})
}
