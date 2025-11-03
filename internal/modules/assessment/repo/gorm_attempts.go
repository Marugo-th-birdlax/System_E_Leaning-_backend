package repo

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/assessment/models"
)

type AttemptRepo struct{ db *gorm.DB }

func NewAttemptRepo(db *gorm.DB) *AttemptRepo { return &AttemptRepo{db: db} }

/******** Attempts ********/
func (r *AttemptRepo) CountUserAttempts(assessmentID, userID string) (int64, error) {
	var n int64
	err := r.db.Model(&models.Attempt{}).
		Where("assessment_id=? AND user_id=?", assessmentID, userID).
		Count(&n).Error
	return n, err
}
func (r *AttemptRepo) CreateAttempt(a *models.Attempt) error {
	return r.db.Create(a).Error
}
func (r *AttemptRepo) GetAttempt(id, userID string) (*models.Attempt, error) {
	var a models.Attempt
	if err := r.db.First(&a, "id=? AND user_id=?", id, userID).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
func (r *AttemptRepo) UpdateAttempt(a *models.Attempt) error {
	return r.db.Model(&models.Attempt{}).Where("id=?", a.ID).Updates(a).Error
}

/******** Answers ********/
func (r *AttemptRepo) UpsertAnswer(ans *models.Answer) error {
	// unique key: attempt_id + question_id
	var existing models.Answer
	err := r.db.Where("attempt_id=? AND question_id=?", ans.AttemptID, ans.QuestionID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(ans).Error
	}
	if err != nil {
		return err
	}
	// update
	payload := map[string]any{
		"selected_choice_ids": ans.SelectedChoiceIDs,
		"text_answer":         ans.TextAnswer,
		"is_correct":          ans.IsCorrect,
	}
	return r.db.Model(&existing).Updates(payload).Error
}

func (r *AttemptRepo) ListAnswers(attemptID string) ([]models.Answer, error) {
	var rows []models.Answer
	err := r.db.Where("attempt_id=?", attemptID).Find(&rows).Error
	return rows, err
}

/******** Helpers for grading ********/
type QuestionWithChoices struct {
	QID        string
	Type       string
	Points     int
	CorrectCSV string // correct choice ids joined by comma
}

// internal/modules/assessment/repo/gorm_repo.go (หรือไฟล์ repo ที่คุณวาง AttemptRepo ไว้)

func (r *AttemptRepo) ListQuestionsWithCorrectChoices(assessmentID string) ([]QuestionWithChoices, error) {
	// 1) ดึงคำถามของ assessment นี้
	var qs []models.Question
	if err := r.db.
		Where("assessment_id = ?", assessmentID).
		Order("seq ASC").
		Find(&qs).Error; err != nil {
		return nil, err
	}
	if len(qs) == 0 {
		return nil, nil
	}

	// 2) รวบรวม question_ids
	qIDs := make([]string, 0, len(qs))
	for _, q := range qs {
		qIDs = append(qIDs, q.ID)
	}

	// 3) ดึง choices ที่ถูกต้องของชุดคำถามนี้
	var cs []models.Choice
	if err := r.db.
		Where("question_id IN ? AND is_correct = ?", qIDs, true).
		Order("seq ASC").
		Find(&cs).Error; err != nil {
		return nil, err
	}

	// 4) map choices ต่อคำถาม
	byQ := make(map[string][]string, len(qs))
	for _, c := range cs {
		byQ[c.QuestionID] = append(byQ[c.QuestionID], c.ID)
	}

	// 5) ประกอบผลลัพธ์ที่ service ต้องใช้
	out := make([]QuestionWithChoices, 0, len(qs))
	for _, q := range qs {
		corr := NormalizeCSV(byQ[q.ID]) // join ด้วย comma (หรือจะเก็บ []string ก็ได้)
		out = append(out, QuestionWithChoices{
			QID:        q.ID,
			Type:       q.Type,
			Points:     q.Points,
			CorrectCSV: corr,
		})
	}
	return out, nil
}

func NormalizeCSV(arr []string) string {
	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
	}
	return strings.Join(arr, ",")
}

func SplitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (r *AttemptRepo) Now() time.Time { return time.Now() }
