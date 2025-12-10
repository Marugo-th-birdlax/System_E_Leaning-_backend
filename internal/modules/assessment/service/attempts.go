package service

import (
	"errors"
	"sort"
	"strings"

	"github.com/Marugo/birdlax/internal/modules/assessment/dto"
	"github.com/Marugo/birdlax/internal/modules/assessment/models"
	assrepo "github.com/Marugo/birdlax/internal/modules/assessment/repo"
	learnmodels "github.com/Marugo/birdlax/internal/modules/learning/models"
	learningservice "github.com/Marugo/birdlax/internal/modules/learning/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttemptService interface {
	StartAttempt(userID, assessmentID string, req dto.StartAttemptReq) (*models.Attempt, error)
	UpsertAnswer(userID, attemptID string, req dto.UpsertAnswerReq) (*models.Answer, error)
	SubmitAttempt(userID, attemptID string, req dto.SubmitAttemptReq) (*models.Attempt, int, int, error)
	GetAttempt(userID, attemptID string) (*models.Attempt, error)
}

type attemptSvc struct {
	repo Repo        // 👈 ใช้ Repo จาก ports.go
	rr   AttemptRepo // attempt repo
	er   EnrollmentRepo
	ms   learningservice.MetricsService
}

func NewAttemptService(repo Repo, rr AttemptRepo, er EnrollmentRepo, ms learningservice.MetricsService) AttemptService {
	return &attemptSvc{repo: repo, rr: rr, er: er, ms: ms}
}

func (s *attemptSvc) StartAttempt(userID, assessmentID string, _ dto.StartAttemptReq) (*models.Attempt, error) {
	a, err := s.repo.GetAssessmentByID(assessmentID)
	if err != nil {
		return nil, errors.New("assessment not found")
	}

	// MaxAttempts check
	if a.MaxAttempts != nil && *a.MaxAttempts > 0 {
		n, err := s.rr.CountUserAttempts(assessmentID, userID)
		if err != nil {
			return nil, err
		}
		if int(n) >= *a.MaxAttempts {
			return nil, errors.New("max attempts reached")
		}
	}

	now := s.rr.Now()
	at := &models.Attempt{
		ID:           uuid.NewString(),
		AssessmentID: assessmentID,
		UserID:       userID,
		Status:       "in_progress",
		StartedAt:    now,
		TimeLimitS:   a.TimeLimitS, // snapshot ไว้
	}
	if err := s.rr.CreateAttempt(at); err != nil {
		return nil, err
	}
	return at, nil
}

func (s *attemptSvc) GetAttempt(userID, attemptID string) (*models.Attempt, error) {
	return s.rr.GetAttempt(attemptID, userID)
}

func (s *attemptSvc) UpsertAnswer(userID, attemptID string, req dto.UpsertAnswerReq) (*models.Answer, error) {
	at, err := s.rr.GetAttempt(attemptID, userID)
	if err != nil {
		return nil, errors.New("attempt not found")
	}
	if at.Status != "in_progress" {
		return nil, errors.New("attempt not editable")
	}
	// Time limit
	if at.TimeLimitS != nil && *at.TimeLimitS > 0 {
		elapsed := int(s.rr.Now().Sub(at.StartedAt).Seconds())
		if elapsed > *at.TimeLimitS {
			at.Status = "expired"
			_ = s.rr.UpdateAttempt(at)
			return nil, errors.New("attempt expired")
		}
	}

	ans := &models.Answer{
		ID:         uuid.NewString(),
		AttemptID:  attemptID,
		QuestionID: req.QuestionID,
	}
	if req.TextAnswer != nil {
		ans.TextAnswer = req.TextAnswer
	}
	if len(req.SelectedChoiceIDs) > 0 {
		csv := assrepo.NormalizeCSV(req.SelectedChoiceIDs)
		ans.SelectedChoiceIDs = &csv
	}
	// is_correct จะคำนวณตอน submit
	if err := s.rr.UpsertAnswer(ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func (s *attemptSvc) SubmitAttempt(userID, attemptID string, _ dto.SubmitAttemptReq) (*models.Attempt, int, int, error) {
	at, err := s.rr.GetAttempt(attemptID, userID)
	if err != nil {
		return nil, 0, 0, errors.New("attempt not found")
	}
	if at.Status != "in_progress" {
		return nil, 0, 0, errors.New("attempt not editable")
	}
	// Time limit
	if at.TimeLimitS != nil && *at.TimeLimitS > 0 {
		elapsed := int(s.rr.Now().Sub(at.StartedAt).Seconds())
		if elapsed > *at.TimeLimitS {
			at.Status = "expired"
			_ = s.rr.UpdateAttempt(at)
			return nil, 0, 0, errors.New("attempt expired")
		}
	}

	// Load questions & answers
	qrows, err := s.rr.ListQuestionsWithCorrectChoices(at.AssessmentID)
	if err != nil {
		return nil, 0, 0, err
	}
	answers, err := s.rr.ListAnswers(at.ID)
	if err != nil {
		return nil, 0, 0, err
	}

	// Map answers by QID
	ansByQ := map[string]models.Answer{}
	for _, a := range answers {
		ansByQ[a.QuestionID] = a
	}

	totalQ := len(qrows)
	raw := 0
	correctCnt := 0

	for _, q := range qrows {
		ans, ok := ansByQ[q.QID]
		if !ok {
			continue // no answer
		}
		var isCorrect bool

		switch q.Type {
		case "single_choice", "true_false":
			if ans.SelectedChoiceIDs == nil {
				isCorrect = false
			} else {
				// ใช้ compare แบบ set แทน ==
				aSet := assrepo.SplitCSV(*ans.SelectedChoiceIDs)
				cSet := assrepo.SplitCSV(q.CorrectCSV)
				sort.Strings(aSet)
				sort.Strings(cSet)
				isCorrect = strings.Join(aSet, ",") == strings.Join(cSet, ",")
			}
		case "multiple_choice":
			if ans.SelectedChoiceIDs == nil {
				isCorrect = false
			} else {
				aSet := assrepo.SplitCSV(*ans.SelectedChoiceIDs)
				cSet := assrepo.SplitCSV(q.CorrectCSV)
				sort.Strings(aSet)
				sort.Strings(cSet)
				isCorrect = strings.Join(aSet, ",") == strings.Join(cSet, ",")
			}

		case "short_text":
			// ยังไม่ auto-grade
			isCorrect = false
		default:
			isCorrect = false
		}

		// update answer.is_correct
		v := isCorrect
		ans.IsCorrect = &v
		_ = s.rr.UpsertAnswer(&ans)

		if isCorrect {
			raw += q.Points
			correctCnt++
		}
	}

	// percent & pass/fail
	ass, err := s.repo.GetAssessmentByID(at.AssessmentID)
	if err != nil {
		return nil, 0, 0, err
	}

	var maxPoints int
	for _, q := range qrows {
		maxPoints += q.Points
	}
	percent := 0.0
	if maxPoints > 0 {
		percent = float64(raw) / float64(maxPoints) * 100.0
	}
	pass := percent >= float64(ass.PassScore)

	now := s.rr.Now()
	at.Status = "submitted"
	at.SubmittedAt = &now
	at.ScoreRaw = &raw
	at.ScorePercent = &percent
	at.IsPassed = &pass

	if err := s.rr.UpdateAttempt(at); err != nil {
		return nil, 0, 0, err
	}

	// 🔗 ผูกผล post-test → enrollment
	_ = s.updateEnrollmentFromAttempt(userID, ass, at)

	var elapsedSec int64 = 0
	if !at.StartedAt.IsZero() && at.SubmittedAt != nil {
		elapsedSec = int64(at.SubmittedAt.Sub(at.StartedAt).Seconds())
	}

	// Use the assessment we already loaded above (ass)
	if ass != nil && ass.OwnerType == "course" {
		courseID := ass.OwnerID
		// call metrics svc async-safe (we ignore error but could log)
		if s.ms != nil {
			_ = s.ms.OnSubmitAttempt(at.UserID, courseID, percent, pass, elapsedSec)
		}
	}
	return at, totalQ, correctCnt, nil

}

func (s *attemptSvc) updateEnrollmentFromAttempt(userID string, ass *models.Assessment, at *models.Attempt) error {
	// สนใจเฉพาะ post-test ที่ผูกกับ "course"
	if ass == nil || ass.OwnerType != "course" || ass.Type != "post" {
		return nil
	}
	if s.er == nil || at == nil || at.IsPassed == nil {
		return nil
	}

	// หา enrollment ถ้าไม่มีก็ไม่สร้างให้ (ตาม requirement)
	e, err := s.er.GetEnrollment(userID, ass.OwnerID)
	if err != nil {
		// ถ้าเป็น not found -> ไม่ต้องทำอะไร
		// แต่ถ้าเป็น error จริง ควรคืน error เพื่อให้ caller เห็น
		if err == gorm.ErrRecordNotFound { // ถ้าใช้ gorm
			return nil
		}
		return err
	}

	now := s.rr.Now()

	// update basic timestamps
	e.LastAccessedAt = &now
	if e.StartedAt == nil {
		e.StartedAt = &now
	}

	// ถ้าผ่าน post-test -> ให้ถือว่า complete / passed
	if *at.IsPassed {
		e.Status = learnmodels.StatusPassed
		if e.CompletedAt == nil {
			e.CompletedAt = &now
		}
		// set progress 100% เป็น conservative update
		if e.ProgressPercent < 100 {
			e.ProgressPercent = 100
		}
	} else {
		// ไม่ผ่าน -> mark failed (แต่ไม่ลด progress)
		e.Status = learnmodels.StatusFailed
	}

	// persist and return error if any
	if err := s.er.UpsertEnrollment(e); err != nil {
		return err
	}

	// ถ้ามี metrics service ให้เรียกเพื่ออัปเดต (ไม่ทำให้ flow fail ถ้า metrics ล้ม)
	if s.ms != nil && *at.IsPassed {
		// ระยะเวลา: ถ้า available ใช้ submitted-started, ถ้าไม่มีก็ส่ง 0
		var elapsed int64 = 0
		if !at.StartedAt.IsZero() && at.SubmittedAt != nil {
			elapsed = int64(at.SubmittedAt.Sub(at.StartedAt).Seconds())
		}
		_ = s.ms.OnCourseCompleted(userID, ass.OwnerID, elapsed) // ignore error, log ถ้าต้องการ
	}

	return nil
}
