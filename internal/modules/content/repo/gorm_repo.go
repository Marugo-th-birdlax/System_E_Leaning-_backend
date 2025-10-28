package repo

import (
	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/content/models"
)

type AssetRepo struct{ db *gorm.DB }
type LessonRepo struct{ db *gorm.DB }

func NewAssetRepo(db *gorm.DB) *AssetRepo   { return &AssetRepo{db: db} }
func NewLessonRepo(db *gorm.DB) *LessonRepo { return &LessonRepo{db: db} }

func (r *AssetRepo) CreateAsset(a *models.Asset) error    { return r.db.Create(a).Error }
func (r *LessonRepo) CreateLesson(l *models.Lesson) error { return r.db.Create(l).Error }

func (r *LessonRepo) GetLessonsByModule(moduleID string) ([]models.Lesson, error) {
	var rows []models.Lesson
	err := r.db.Where("module_id = ?", moduleID).
		Order("seq ASC").
		Find(&rows).Error
	return rows, err
}

// --- Lesson ---
func (r *LessonRepo) GetByID(id string) (*models.Lesson, error) {
	var l models.Lesson
	if err := r.db.First(&l, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LessonRepo) List(moduleID string, page, per int) ([]models.Lesson, int64, error) {
	if page < 1 {
		page = 1
	}
	if per <= 0 || per > 100 {
		per = 20
	}
	tx := r.db.Model(&models.Lesson{})
	if moduleID != "" {
		tx = tx.Where("module_id = ?", moduleID)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Lesson
	if err := tx.Order("module_id ASC, seq ASC").
		Limit(per).Offset((page - 1) * per).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// --- Asset ---
func (r *AssetRepo) GetByID(id string) (*models.Asset, error) {
	var a models.Asset
	if err := r.db.First(&a, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
