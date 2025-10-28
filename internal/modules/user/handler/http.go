package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/Marugo/birdlax/internal/modules/user"
	"github.com/Marugo/birdlax/internal/modules/user/dto"
	userservice "github.com/Marugo/birdlax/internal/modules/user/service"
	"github.com/Marugo/birdlax/internal/shared/response"
)

type HTTPHandler struct {
	svc       user.Service
	validator *validator.Validate
}

func NewHTTPHandler(s user.Service) *HTTPHandler {
	return &HTTPHandler{svc: s, validator: validator.New()}
}

func (h *HTTPHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	q := c.Query("q", "")
	items, total, err := h.svc.List(c.Context(), page, perPage, q)
	if err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	out := make([]*dto.UserResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.FromModel(&items[i]))
	}
	return response.OK(c, fiber.Map{
		"items":    out,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *HTTPHandler) Create(c *fiber.Ctx) error {
	var req dto.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	u, err := h.svc.Create(
		c.Context(),
		req.EmployeeCode, req.Email, req.FirstName, req.LastName, req.Role, req.Phone,
		req.Password, // <-- ส่งรหัสผ่านไปให้ service แฮช
	)
	if err != nil {
		if errors.Is(err, userservice.ErrEmployeeCodeAlreadyExist) {
			return response.Err(c, fiber.StatusConflict, "EMPLOYEE_CODE_ALREADY_EXISTS", "employee_code is already registered")
		}
		if errors.Is(err, userservice.ErrEmailAlreadyExists) {
			return response.Err(c, fiber.StatusConflict, "EMAIL_ALREADY_EXISTS", "email is already registered")
		}
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error())
	}
	return response.OK(c, dto.FromModel(u))
}

func (h *HTTPHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	u, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return response.Err(c, fiber.StatusNotFound, "NOT_FOUND", "user not found")
	}
	return response.OK(c, dto.FromModel(u))
}

func (h *HTTPHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	u, err := h.svc.Update(
		c.Context(), id,
		req.EmployeeCode, req.Email, req.FirstName, req.LastName, req.Role, req.Phone, req.IsActive,
		req.Password, // <-- optional: ถ้า nil หรือ "" จะไม่เปลี่ยน
	)
	if err != nil {
		if errors.Is(err, userservice.ErrEmployeeCodeAlreadyExist) {
			return response.Err(c, fiber.StatusConflict, "EMPLOYEE_CODE_ALREADY_EXISTS", "employee_code is already registered")
		}
		if errors.Is(err, userservice.ErrEmailAlreadyExists) {
			return response.Err(c, fiber.StatusConflict, "EMAIL_ALREADY_EXISTS", "email is already registered")
		}
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error())
	}
	return response.OK(c, dto.FromModel(u))
}

func (h *HTTPHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	return response.OK(c, fiber.Map{"deleted": true})
}
