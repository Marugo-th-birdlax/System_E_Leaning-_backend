package repo

import (
	"strings"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/content/models"
)

type CourseRepo struct{ db *gorm.DB }
type ModuleRepo struct{ db *gorm.DB }

func NewCourseRepo(db *gorm.DB) *CourseRepo { return &CourseRepo{db: db} }
func NewModuleRepo(db *gorm.DB) *ModuleRepo { return &ModuleRepo{db: db} }

/********* Courses *********/
func (r *CourseRepo) Create(c *models.Course) error {
	return r.db.Create(c).Error
}
func (r *CourseRepo) Update(c *models.Course) error {
	return r.db.Model(&models.Course{}).Where("id=?", c.ID).Updates(c).Error
}
func (r *CourseRepo) Delete(id string) error {
	return r.db.Delete(&models.Course{}, "id=?", id).Error
}
func (r *CourseRepo) GetByID(id string) (*models.Course, error) {
	var c models.Course
	if err := r.db.First(&c, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *CourseRepo) List(q string, page, per int) ([]models.Course, int64, error) {
	if page < 1 {
		page = 1
	}
	if per <= 0 || per > 100 {
		per = 20
	}
	var rows []models.Course
	var total int64

	tx := r.db.Model(&models.Course{})
	if s := strings.TrimSpace(q); s != "" {
		tx = tx.Where("code LIKE ? OR title LIKE ?", "%"+s+"%", "%"+s+"%")
	}
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("created_at DESC").Limit(per).Offset((page - 1) * per).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

/********* Modules *********/
func (r *ModuleRepo) Create(m *models.CourseModule) error {
	return r.db.Create(m).Error
}
func (r *ModuleRepo) Update(m *models.CourseModule) error {
	return r.db.Model(&models.CourseModule{}).Where("id=?", m.ID).Updates(m).Error
}
func (r *ModuleRepo) Delete(id string) error {
	return r.db.Delete(&models.CourseModule{}, "id=?", id).Error
}
func (r *ModuleRepo) GetByID(id string) (*models.CourseModule, error) {
	var m models.CourseModule
	if err := r.db.First(&m, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *ModuleRepo) ListByCourse(courseID string) ([]models.CourseModule, error) {
	var rows []models.CourseModule
	err := r.db.Where("course_id=?", courseID).Order("seq ASC").Find(&rows).Error
	return rows, err
}
