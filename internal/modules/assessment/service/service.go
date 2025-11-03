package service

import (
	"errors"
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

	ad := dto.AssessmentItem{
		ID: a.ID, OwnerType: a.OwnerType, OwnerID: a.OwnerID, Type: a.Type,
		Title: a.Title, PassScore: a.PassScore, TimeLimitS: a.TimeLimitS, MaxAttempts: a.MaxAttempts,
	}

	outQ := make([]dto.QuestionResp, 0, len(qs))
	for _, q := range qs {
		qRes := dto.QuestionResp{
			ID: q.ID, Type: q.Type, Stem: strings.TrimSpace(q.Stem), Points: q.Points, Seq: q.Seq,
			// ถ้าต้องการโชว์ assessment_id ด้วย (ตาม output เดิม)
			AssessmentID: q.AssessmentID,
		}
		if cs := csByQ[q.ID]; len(cs) > 0 {
			choices := make([]dto.ChoiceResp, 0, len(cs))
			for _, c := range cs {
				choices = append(choices, dto.ChoiceResp{
					ID: c.ID, Label: strings.TrimSpace(c.Label), IsCorrect: c.IsCorrect, Seq: c.Seq,
				})
			}
			qRes.Choices = choices
		}
		outQ = append(outQ, qRes)
	}
	return &dto.AssessmentDetailResp{Assessment: ad, Questions: outQ}, nil
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

func (s *svc) ReplaceChoices(questionID string, req dto.ReplaceChoicesReq) ([]dto.ChoiceResp, error) {
	// ensure question exists และ belong-check
	if _, err := s.repo.GetQuestionByID(questionID); err != nil {
		return nil, errors.New("question not found")
	}
	// map dto -> model
	items := make([]models.Choice, 0, len(req.Choices))
	for _, c := range req.Choices {
		id := ""
		if c.ID != nil {
			id = *c.ID
		}
		items = append(items, models.Choice{
			ID:         id,
			QuestionID: questionID,
			Label:      c.Label,
			IsCorrect:  c.IsCorrect,
			Seq:        c.Seq,
		})
	}
	if err := s.repo.ReplaceChoices(questionID, items); err != nil {
		return nil, err
	}
	// return รายการล่าสุด
	rows, err := s.repo.ListChoicesByQuestion(questionID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ChoiceResp, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.ChoiceResp{ID: r.ID, Label: r.Label, IsCorrect: r.IsCorrect, Seq: r.Seq})
	}
	return out, nil
}

func (s *svc) AddChoice(questionID string, in dto.ChoiceUpsert) (*dto.ChoiceResp, error) {
	if _, err := s.repo.GetQuestionByID(questionID); err != nil {
		return nil, errors.New("question not found")
	}
	m := models.Choice{
		ID:         uuid.NewString(),
		QuestionID: questionID,
		Label:      in.Label,
		IsCorrect:  in.IsCorrect,
		Seq:        in.Seq,
	}
	if err := s.repo.CreateChoice(&m); err != nil {
		return nil, err
	}
	return &dto.ChoiceResp{ID: m.ID, Label: m.Label, IsCorrect: m.IsCorrect, Seq: m.Seq}, nil
}

func (s *svc) UpdateChoice(choiceID string, in dto.ChoiceUpsert) (*dto.ChoiceResp, error) {
	// โหลดของเดิมก่อนเพื่อ ensure มีอยู่
	// สามารถเพิ่ม repo.GetChoiceByID ถ้าต้องการเข้มงวด
	m := models.Choice{
		ID:        choiceID,
		Label:     in.Label,
		IsCorrect: in.IsCorrect,
		Seq:       in.Seq,
	}
	if err := s.repo.UpdateChoice(&m); err != nil {
		return nil, err
	}
	return &dto.ChoiceResp{ID: m.ID, Label: m.Label, IsCorrect: m.IsCorrect, Seq: m.Seq}, nil
}

func (s *svc) DeleteChoice(choiceID string) error {
	return s.repo.DeleteChoice(choiceID)
}
