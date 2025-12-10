package repo

import (
	"errors"
	"time"

	"gorm.io/gorm"

	lmodels "github.com/Marugo/birdlax/internal/modules/learning/models"
	"github.com/google/uuid"
)

type MetricsRepo struct {
	db *gorm.DB
}

func NewMetricsRepo(db *gorm.DB) *MetricsRepo { return &MetricsRepo{db: db} }

func (r *MetricsRepo) GetLearningMetric(userID, courseID string) (*lmodels.LearningMetric, error) {
	var m lmodels.LearningMetric
	if err := r.db.First(&m, "user_id = ? AND course_id = ?", userID, courseID).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MetricsRepo) UpsertLearningMetric(m *lmodels.LearningMetric) error {
	// try update existing first
	var existing lmodels.LearningMetric
	err := r.db.Where("user_id = ? AND course_id = ?", m.UserID, m.CourseID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if m.ID == "" {
			m.ID = uuid.NewString()
		}
		return r.db.Create(m).Error
	}
	if err != nil {
		return err
	}
	m.ID = existing.ID
	m.UpdatedAt = time.Now()
	return r.db.Model(&lmodels.LearningMetric{}).
		Where("id = ?", m.ID).
		Updates(map[string]any{
			"avg_score":          m.AvgScore,
			"last_score":         m.LastScore,
			"attempts_count":     m.AttemptsCount,
			"pass_count":         m.PassCount,
			"total_time_seconds": m.TotalTimeSeconds,
			"completion_status":  m.CompletionStatus,
			"updated_at":         m.UpdatedAt,
		}).Error
}

func (r *MetricsRepo) GetCourseOutcome(courseID string) (*lmodels.CourseOutcome, error) {
	var c lmodels.CourseOutcome
	if err := r.db.First(&c, "course_id = ?", courseID).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MetricsRepo) UpsertCourseOutcome(c *lmodels.CourseOutcome) error {
	var existing lmodels.CourseOutcome
	err := r.db.First(&existing, "course_id = ?", c.CourseID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(c).Error
	}
	if err != nil {
		return err
	}
	c.UpdatedAt = time.Now()
	return r.db.Model(&lmodels.CourseOutcome{}).
		Where("course_id = ?", c.CourseID).
		Updates(map[string]any{
			"total_enrollments":   c.TotalEnrollments,
			"total_completed":     c.TotalCompleted,
			"avg_score":           c.AvgScore,
			"pass_rate":           c.PassRate,
			"median_time_seconds": c.MedianTimeSeconds,
			"updated_at":          c.UpdatedAt,
		}).Error
}
