package handler

import (
	"strconv"

	"github.com/Marugo/birdlax/internal/modules/assessment/dto"
	"github.com/Marugo/birdlax/internal/modules/assessment/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{ svc service.Service }

func New(s service.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) CreateAssessment(c *fiber.Ctx) error {
	var req dto.CreateAssessmentReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	a, err := h.svc.Create(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(a)
}

func (h *Handler) AddQuestion(c *fiber.Ctx) error {
	assessID := c.Params("id")
	var req dto.AddQuestionReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	q, err := h.svc.AddQuestion(assessID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(q)
}

type AttemptHandler struct{ svc service.AttemptService }

func NewAttemptHandler(s service.AttemptService) *AttemptHandler { return &AttemptHandler{svc: s} }

func uid(c *fiber.Ctx) string {
	if v := c.Locals("user_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// --- NEW: GET /assessments ---
func (h *Handler) ListAssessments(c *fiber.Ctx) error {
	filter := dto.ListAssessmentsFilter{
		OwnerType: c.Query("owner_type"),
		OwnerID:   c.Query("owner_id"),
		Type:      c.Query("type"),
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	per, _ := strconv.Atoi(c.Query("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if per <= 0 || per > 100 {
		per = 20
	}

	items, total, err := h.svc.List(filter, page, per)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(dto.PagedAssessments{
		Data: items,
		Meta: dto.PageMeta{Page: page, PerPage: per, Total: total},
	})
}

// --- NEW: GET /assessments/:id (with questions & choices) ---

func (h *Handler) GetAssessment(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.svc.GetDetail(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "assessment not found")
	}
	return c.JSON(res)
}

// POST /v1/assessments/:id/attempts
func (h *AttemptHandler) Start(c *fiber.Ctx) error {
	userID := uid(c)
	assessID := c.Params("id")
	var req dto.StartAttemptReq
	_ = c.BodyParser(&req)
	at, err := h.svc.StartAttempt(userID, assessID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(at)
}

// GET /v1/attempts/:id
func (h *AttemptHandler) Get(c *fiber.Ctx) error {
	userID := uid(c)
	at, err := h.svc.GetAttempt(userID, c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "attempt not found")
	}
	return c.JSON(at)
}

// POST /v1/attempts/:id/answers
func (h *AttemptHandler) UpsertAnswer(c *fiber.Ctx) error {
	userID := uid(c)
	var req dto.UpsertAnswerReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	ans, err := h.svc.UpsertAnswer(userID, c.Params("id"), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(ans)
}

// POST /v1/attempts/:id/submit
func (h *AttemptHandler) Submit(c *fiber.Ctx) error {
	userID := uid(c)
	var req dto.SubmitAttemptReq
	_ = c.BodyParser(&req)
	at, total, correct, err := h.svc.SubmitAttempt(userID, c.Params("id"), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(fiber.Map{
		"attempt": at,
		"summary": fiber.Map{"total_questions": total, "correct": correct},
	})
}

// PUT /v1/assessments/:id
func (h *Handler) UpdateAssessment(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateAssessmentReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	a, err := h.svc.UpdateAssessment(id, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(a)
}

// DELETE /v1/assessments/:id
func (h *Handler) DeleteAssessment(c *fiber.Ctx) error {
	if err := h.svc.DeleteAssessment(c.Params("id")); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// PUT /v1/assessments/:id/questions/:qid
func (h *Handler) UpdateQuestion(c *fiber.Ctx) error {
	qid := c.Params("qid")
	var req dto.UpdateQuestionReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	q, err := h.svc.UpdateQuestion(qid, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(q)
}

// DELETE /v1/assessments/:id/questions/:qid
func (h *Handler) DeleteQuestion(c *fiber.Ctx) error {
	qid := c.Params("qid")
	if err := h.svc.DeleteQuestion(qid); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
