package handler

import (
	"strconv"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/service"
	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	svc service.CourseService
}

func NewCourseHandler(s service.CourseService) *CourseHandler { return &CourseHandler{svc: s} }

/********* Courses *********/
func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	var req dto.CreateCourseReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	course, err := h.svc.CreateCourse(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// 🔥 ดึง department_ids เพิ่ม
	deptIDs, err := h.svc.ListCourseDepartments(course.ID)
	if err != nil {
		// ถ้าอยากจะ ignore error ก็ได้ แต่เอาแบบตรงไปตรงมา
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resp := dto.FromCourseModel(course, deptIDs)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateCourseReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	course, err := h.svc.UpdateCourse(id, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// 🔥 หลังอัปเดตแล้ว ดึง department_ids ปัจจุบัน
	deptIDs, err := h.svc.ListCourseDepartments(course.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resp := dto.FromCourseModel(course, deptIDs)
	return c.JSON(resp)
}

func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	if err := h.svc.DeleteCourse(c.Params("id")); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *CourseHandler) GetCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	course, err := h.svc.GetCourse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "course not found")
	}

	// 🔥 ดึง department_ids เพิ่มจาก service
	deptIDs, err := h.svc.ListCourseDepartments(id)
	if err != nil {
		// จะ ignore error ก็ได้ แต่เอาแบบตรง ๆ ไว้ก่อน
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resp := dto.FromCourseModel(course, deptIDs)
	return c.JSON(resp)
}
func (h *CourseHandler) ListCourses(c *fiber.Ctx) error {
	q := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))
	rows, total, err := h.svc.ListCourses(q, page, per)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	out := make([]dto.CourseResp, 0, len(rows))
	for _, x := range rows {
		out = append(out, dto.CourseResp{
			ID: x.ID, Code: x.Code, Title: x.Title, Description: x.Description,
			IsActive: x.IsActive, EstimatedMinutes: x.EstimatedMinutes,
			CategoryID: x.CategoryID,
		})
	}
	return c.JSON(dto.PagedCourses{
		Data: out, Meta: dto.PageMeta{Page: page, PerPage: per, Total: total},
	})
}

/********* Modules *********/
func (h *CourseHandler) CreateModule(c *fiber.Ctx) error {
	courseID := c.Params("id")
	var req dto.CreateModuleReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	m, err := h.svc.CreateModule(courseID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(m)
}
func (h *CourseHandler) UpdateModule(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateModuleReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	m, err := h.svc.UpdateModule(id, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(m)
}
func (h *CourseHandler) DeleteModule(c *fiber.Ctx) error {
	if err := h.svc.DeleteModule(c.Params("id")); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *CourseHandler) ListModules(c *fiber.Ctx) error {
	rows, err := h.svc.ListModules(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	out := make([]dto.ModuleResp, 0, len(rows))
	for _, x := range rows {
		out = append(out, dto.ModuleResp{
			ID: x.ID, CourseID: x.CourseID, Title: x.Title, Description: x.Description,
			Seq: x.Seq, IsMandatory: x.IsMandatory,
		})
	}
	return c.JSON(out)
}
func (h *CourseHandler) ListLessonsOfModule(c *fiber.Ctx) error {
	rows, err := h.svc.ListLessons(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(rows) // ใช้ struct Lesson ตรงๆ (มี id, module_id, title, content_type, seq, ...)
}
