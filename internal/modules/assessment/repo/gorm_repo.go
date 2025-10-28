package repo

import (
	"strings"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/assessment/models"
)

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) CreateAssessment(a *models.Assessment) error {
	return r.db.Create(a).Error
}

func (r *Repo) AddQuestion(q *models.Question, choices []models.Choice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(q).Error; err != nil {
			return err
		}
		if len(choices) > 0 {
			if err := tx.Create(&choices).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) GetAssessmentByID(id string) (*models.Assessment, error) {
	var a models.Assessment
	if err := r.db.First(&a, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repo) ListAssessments(ownerType, ownerID, atype string, page, per int) ([]models.Assessment, int64, error) {
	if page < 1 {
		page = 1
	}
	if per <= 0 || per > 100 {
		per = 20
	}

	tx := r.db.Model(&models.Assessment{})
	if s := strings.TrimSpace(ownerType); s != "" {
		tx = tx.Where("owner_type = ?", s)
	}
	if s := strings.TrimSpace(ownerID); s != "" {
		tx = tx.Where("owner_id = ?", s)
	}
	if s := strings.TrimSpace(atype); s != "" {
		tx = tx.Where("type = ?", s)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.Assessment
	if err := tx.Order("created_at DESC").
		Limit(per).Offset((page - 1) * per).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

/******** DETAIL (assessment + questions + choices) ********/
func (r *Repo) GetAssessmentWithItems(id string) (*models.Assessment, []models.Question, map[string][]models.Choice, error) {
	var a models.Assessment
	if err := r.db.First(&a, "id=?", id).Error; err != nil {
		return nil, nil, nil, err
	}

	var qs []models.Question
	if err := r.db.Where("assessment_id = ?", id).Order("seq ASC").Find(&qs).Error; err != nil {
		return &a, nil, nil, err
	}

	csByQ := make(map[string][]models.Choice, len(qs))
	if len(qs) > 0 {
		qIDs := make([]string, 0, len(qs))
		for _, q := range qs {
			qIDs = append(qIDs, q.ID)
		}

		var cs []models.Choice
		if err := r.db.Where("question_id IN ?", qIDs).Order("seq ASC").Find(&cs).Error; err == nil {
			for _, c := range cs {
				csByQ[c.QuestionID] = append(csByQ[c.QuestionID], c)
			}
		}
	}
	return &a, qs, csByQ, nil
}

func (r *Repo) UpdateAssessment(a *models.Assessment) error {
	return r.db.Model(&models.Assessment{}).Where("id = ?", a.ID).Updates(a).Error
}
func (r *Repo) DeleteAssessment(id string) error {
	return r.db.Delete(&models.Assessment{}, "id = ?", id).Error
}

// --- QUESTION: Get/Update/Delete ---
func (r *Repo) GetQuestionByID(id string) (*models.Question, error) {
	var q models.Question
	if err := r.db.First(&q, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &q, nil
}
func (r *Repo) UpdateQuestion(q *models.Question) error {
	return r.db.Model(&models.Question{}).Where("id = ?", q.ID).Updates(q).Error
}
func (r *Repo) DeleteQuestion(id string) error {
	// ลบ choice ที่ผูกด้วย เพื่อความสะอาด (on delete cascade ก็ได้ถ้าตั้ง FK)
	if err := r.db.Where("question_id = ?", id).Delete(&models.Choice{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Question{}, "id = ?", id).Error
}
