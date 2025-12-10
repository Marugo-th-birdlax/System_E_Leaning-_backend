package handler

import (
	"errors"
	"strconv"
	"strings"

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

// ===================== Users =====================

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
	actorRole, _ := c.Locals("role").(string)

	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	if req.Role != "" && actorRole != "admin" && actorRole != "hr" {
		return response.Err(c, fiber.StatusForbidden, "FORBIDDEN", "not allowed to set role")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	// ถ้าไม่ได้ส่ง role มา → ตั้ง default เป็น employee
	roleInput := req.Role
	if strings.TrimSpace(roleInput) == "" {
		roleInput = "employee"
	}
	u, err := h.svc.Create(c.Context(),
		req.EmployeeCode, req.Email, req.FirstName, req.LastName,
		roleInput, req.Phone, req.Password,
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
	actorRole, _ := c.Locals("role").(string)
	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	// ถ้าไม่ใช่ admin/hr และมีการส่ง req.Role → ห้าม
	if req.Role != nil && actorRole != "admin" && actorRole != "hr" {
		return response.Err(c, fiber.StatusForbidden, "FORBIDDEN", "not allowed to change role")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	u, err := h.svc.Update(
		c.Context(), id,
		req.EmployeeCode, req.Email, req.FirstName, req.LastName, req.Role, req.Phone, req.IsActive,
		req.Password,
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

// ===================== Departments =====================

func (h *HTTPHandler) ListDepartments(c *fiber.Ctx) error {
	q := c.Query("q", "")
	items, err := h.svc.ListDepartments(c.Context(), q)
	if err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	out := make([]*dto.DepartmentResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.FromDepartmentModel(&items[i])) // ✅ ใช้ dto.
	}
	return response.OK(c, out)
}

func (h *HTTPHandler) CreateDepartment(c *fiber.Ctx) error {
	var req dto.DepartmentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	d, err := h.svc.CreateDepartment(c.Context(), req.Code, req.Name, req.IsActive)
	if err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error())
	}
	return response.OK(c, dto.FromDepartmentModel(d)) // ✅
}

func (h *HTTPHandler) GetDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	d, err := h.svc.GetDepartment(c.Context(), id)
	if err != nil {
		return response.Err(c, fiber.StatusNotFound, "NOT_FOUND", "department not found")
	}
	return response.OK(c, dto.FromDepartmentModel(d)) // ✅
}

func (h *HTTPHandler) UpdateDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.DepartmentUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	d, err := h.svc.UpdateDepartment(c.Context(), id, req.Code, req.Name, req.IsActive)
	if err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error())
	}
	return response.OK(c, dto.FromDepartmentModel(d)) // ✅
}

func (h *HTTPHandler) DeleteDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteDepartment(c.Context(), id); err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	return response.OK(c, fiber.Map{"deleted": true})
}

// GET /departments/:id/managers
func (h *HTTPHandler) ListDepartmentManagers(c *fiber.Ctx) error {
	id := c.Params("id")
	items, err := h.svc.ListDepartmentManagers(c.Context(), id)
	if err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	out := make([]*dto.UserDepartmentRoleResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.FromUserDepartmentRoleModel(&items[i])) // ✅
	}
	return response.OK(c, out)
}

// ===================== User Department Roles =====================

// GET /users/:id/department-roles
func (h *HTTPHandler) ListUserDepartmentRoles(c *fiber.Ctx) error {
	id := c.Params("id")
	items, err := h.svc.ListUserDepartmentRoles(c.Context(), id)
	if err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	out := make([]*dto.UserDepartmentRoleResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.FromUserDepartmentRoleModel(&items[i]))
	}
	return response.OK(c, out)
}

// POST /users/:id/department-roles
func (h *HTTPHandler) AddUserDepartmentRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.AssignDepartmentRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	if err := h.validator.Struct(&req); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	rel, err := h.svc.AddUserDepartmentRole(c.Context(), id, req.DepartmentID, req.Role)
	if err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error())
	}
	return response.OK(c, dto.FromUserDepartmentRoleModel(rel)) // ✅
}

// DELETE /users/:id/department-roles/:relID
func (h *HTTPHandler) DeleteUserDepartmentRole(c *fiber.Ctx) error {
	relID := c.Params("relID")
	if err := h.svc.DeleteUserDepartmentRole(c.Context(), relID); err != nil {
		return response.Err(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	return response.OK(c, fiber.Map{"deleted": true})
}
