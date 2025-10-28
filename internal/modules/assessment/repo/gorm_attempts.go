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

func (r *AttemptRepo) ListQuestionsWithCorrectChoices(assessmentID string) ([]QuestionWithChoices, error) {
	type row struct {
		QID    string
		Type   string
		Points int
		CID    string
		OK     bool
	}
	var rows []row
	err := r.db.Table("assessment_questions q").
		Select("q.id as qid, q.type, q.points, c.id as cid, c.is_correct as ok").
		Joins("LEFT JOIN assessment_choices c ON c.question_id = q.id").
		Where("q.assessment_id = ?", assessmentID).
		Order("q.seq ASC, c.seq ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Aggregate correct choices per question
	m := map[string]QuestionWithChoices{}
	for _, r := range rows {
		x, ok := m[r.QID]
		if !ok {
			x = QuestionWithChoices{QID: r.QID, Type: r.Type, Points: r.Points}
		}
		if r.OK {
			if x.CorrectCSV == "" {
				x.CorrectCSV = r.CID
			} else {
				x.CorrectCSV += "," + r.CID
			}
		}
		m[r.QID] = x
	}
	out := make([]QuestionWithChoices, 0, len(m))
	for _, v := range m {
		out = append(out, v)
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
