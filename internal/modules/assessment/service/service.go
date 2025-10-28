package service

import (
	"strings"

	"github.com/Marugo/birdlax/internal/modules/assessment/dto"
	"github.com/Marugo/birdlax/internal/modules/assessment/models"
	"github.com/google/uuid"
)

type svc struct{ repo Repo }

func New(r Repo) Service { return &svc{repo: r} }

func (s *svc) Create(req dto.CreateAssessmentReq) (*models.Assessment, error) {
	a := &models.Assessment{
		ID:          uuid.NewString(),
		OwnerType:   req.OwnerType,
		OwnerID:     req.OwnerID,
		Type:        req.Type,
		Title:       req.Title,
		PassScore:   req.PassScore,
		MaxAttempts: req.MaxAttempts,
		TimeLimitS:  req.TimeLimitS,
	}
	return a, s.repo.CreateAssessment(a)
}

func (s *svc) AddQuestion(assessID string, req dto.AddQuestionReq) (*models.Question, error) {
	q := &models.Question{
		ID:           uuid.NewString(),
		AssessmentID: assessID,
		Type:         req.Type,
		Stem:         req.Stem,
		Explanation:  req.Explanation,
		Points:       req.Points,
		Seq:          req.Seq,
	}
	var cs []models.Choice
	for _, c := range req.Choices {
		cs = append(cs, models.Choice{
			ID:         uuid.NewString(),
			QuestionID: q.ID,
			Label:      c.Label,
			IsCorrect:  c.IsCorrect,
			Seq:        c.Seq,
		})
	}
	return q, s.repo.AddQuestion(q, cs)
}

func (s *svc) List(filter dto.ListAssessmentsFilter, page, per int) ([]dto.AssessmentItem, int64, error) {
	rows, total, err := s.repo.ListAssessments(filter.OwnerType, filter.OwnerID, filter.Type, page, per)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.AssessmentItem, 0, len(rows))
	for _, a := range rows {
		out = append(out, dto.AssessmentItem{
			ID: a.ID, OwnerType: a.OwnerType, OwnerID: a.OwnerID, Type: a.Type,
			Title: a.Title, PassScore: a.PassScore, TimeLimitS: a.TimeLimitS, MaxAttempts: a.MaxAttempts,
		})
	}
	return out, total, nil
}

func (s *svc) GetDetail(id string) (*dto.AssessmentDetailResp, error) {
	a, qs, csByQ, err := s.repo.GetAssessmentWithItems(id)
	if err != nil {
		return nil, err
	}

	item := dto.AssessmentItem{
		ID: a.ID, OwnerType: a.OwnerType, OwnerID: a.OwnerID, Type: a.Type,
		Title: a.Title, PassScore: a.PassScore, TimeLimitS: a.TimeLimitS, MaxAttempts: a.MaxAttempts,
	}
	qres := make([]dto.QuestionResp, 0, len(qs))
	for _, q := range qs {
		qr := dto.QuestionResp{
			ID: q.ID, Type: q.Type, Stem: strings.TrimSpace(q.Stem), Points: q.Points, Seq: q.Seq,
		}
		if cs, ok := csByQ[q.ID]; ok && len(cs) > 0 {
			// dto.QuestionResp does not have a Choices field, so skip populating choices here.
			_ = cs
		}
		qres = append(qres, qr)
	}
	return &dto.AssessmentDetailResp{Assessment: item, Questions: qres}, nil
}

func (s *svc) UpdateAssessment(id string, req dto.UpdateAssessmentReq) (*models.Assessment, error) {
	a, err := s.repo.GetAssessmentByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		a.Title = *req.Title
	}
	if req.Type != nil {
		a.Type = *req.Type
	}
	if req.PassScore != nil {
		a.PassScore = *req.PassScore
	}
	if req.TimeLimitS != nil {
		a.TimeLimitS = req.TimeLimitS
	}
	if req.MaxAttempts != nil {
		a.MaxAttempts = req.MaxAttempts
	}

	if err := s.repo.UpdateAssessment(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *svc) DeleteAssessment(id string) error {
	return s.repo.DeleteAssessment(id)
}

func (s *svc) UpdateQuestion(id string, req dto.UpdateQuestionReq) (*models.Question, error) {
	q, err := s.repo.GetQuestionByID(id)
	if err != nil {
		return nil, err
	}

	if req.Stem != nil {
		q.Stem = *req.Stem
	}
	if req.Points != nil {
		q.Points = *req.Points
	}
	if req.Seq != nil {
		q.Seq = *req.Seq
	}
	// if req.Type != nil { q.Type = *req.Type }

	if err := s.repo.UpdateQuestion(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *svc) DeleteQuestion(id string) error {
	return s.repo.DeleteQuestion(id)
}
