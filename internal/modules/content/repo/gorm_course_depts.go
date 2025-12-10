package repo

import (
	"strings"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/content/models"
)

type gormCourseDeptRepo struct {
	db *gorm.DB
}

func NewCourseDeptRepo(db *gorm.DB) *gormCourseDeptRepo {
	return &gormCourseDeptRepo{db: db}
}

// ReplaceTargets: ลบ mapping เดิมของ course แล้วใส่ชุดใหม่
func (r *gormCourseDeptRepo) ReplaceTargets(courseID string, deptIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("course_id = ?", courseID).
			Delete(&models.CourseDepartmentTarget{}).Error; err != nil {
			return err
		}
		seen := map[string]struct{}{}
		for _, d := range deptIDs {
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			if _, ok := seen[d]; ok {
				continue
			}
			seen[d] = struct{}{}
			link := &models.CourseDepartmentTarget{
				CourseID:     courseID,
				DepartmentID: d,
				IsMandatory:  false,
			}
			if err := tx.Create(link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormCourseDeptRepo) ListDepartmentIDs(courseID string) ([]string, error) {
	var rows []models.CourseDepartmentTarget
	if err := r.db.
		Where("course_id = ?", courseID).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, v := range rows {
		out = append(out, v.DepartmentID)
	}
	return out, nil
}
