package repo

import (
	"strings"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/content/models"
)

type CategoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db: db} }

// Exists ใช้ validate ตอนผูกคอร์สเข้าหมวด
func (r *CategoryRepo) Exists(id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	var n int64
	if err := r.db.Table("categories").Where("id = ?", id).Limit(1).Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *CategoryRepo) Create(c *models.Category) error { return r.db.Create(c).Error }

func (r *CategoryRepo) Update(c *models.Category) error { return r.db.Save(c).Error }

func (r *CategoryRepo) Delete(id string) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

func (r *CategoryRepo) GetByID(id string) (*models.Category, error) {
	var c models.Category
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepo) List(q string, page, per int) ([]models.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if per <= 0 || per > 100 {
		per = 20
	}
	tx := r.db.Model(&models.Category{})
	if s := strings.TrimSpace(q); s != "" {
		tx = tx.Where("code LIKE ? OR title LIKE ?", "%"+s+"%", "%"+s+"%")
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.Category
	if err := tx.Order("created_at DESC").Limit(per).Offset((page - 1) * per).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
