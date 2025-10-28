package repo

import (
	"errors"

	"gorm.io/gorm"

	content "github.com/Marugo/birdlax/internal/modules/content/models"
	learn "github.com/Marugo/birdlax/internal/modules/learning/models"
)

// NOTE: สมมุติว่ามีตาราง modules ที่มี CourseID (เพื่อคิด % ต่อคอร์ส)
// เสริม struct บางตัวเพื่อ JOIN ได้
type CourseModule struct {
	ID       string `gorm:"type:char(36);primaryKey"`
	CourseID string `gorm:"type:char(36);index"`
}

func (CourseModule) TableName() string { return "course_modules" }

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) *Repo { return &Repo{db: db} }

// Enrollment
func (r *Repo) UpsertEnrollment(e *learn.Enrollment) error {
	// ถ้ามีอยู่แล้ว user+course → update; ไม่มี → create
	var existing learn.Enrollment
	if err := r.db.Where("user_id=? AND course_id=?", e.UserID, e.CourseID).
		First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(e).Error
		}
		return err
	}
	e.ID = existing.ID
	return r.db.Model(&existing).Updates(map[string]any{
		"status":           e.Status,
		"started_at":       e.StartedAt,
		"completed_at":     e.CompletedAt,
		"last_accessed_at": e.LastAccessedAt,
		"progress_percent": e.ProgressPercent,
	}).Error
}

func (r *Repo) GetEnrollment(userID, courseID string) (*learn.Enrollment, error) {
	var e learn.Enrollment
	if err := r.db.Where("user_id=? AND course_id=?", userID, courseID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

// Progress
func (r *Repo) GetLessonProgress(userID, lessonID string) (*learn.UserLessonProgress, error) {
	var p learn.UserLessonProgress
	if err := r.db.Where("user_id=? AND lesson_id=?", userID, lessonID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repo) CreateLessonProgress(p *learn.UserLessonProgress) error {
	return r.db.Create(p).Error
}

func (r *Repo) UpdateLessonProgress(p *learn.UserLessonProgress) error {
	return r.db.Model(&learn.UserLessonProgress{}).
		Where("id=?", p.ID).
		Updates(map[string]any{
			"progress_percent": p.ProgressPercent,
			"current_position": p.CurrentPosition,
			"max_position":     p.MaxPosition,
			"started_at":       p.StartedAt,
			"completed_at":     p.CompletedAt,
			"is_unlocked":      p.IsUnlocked,
		}).Error
}

// Content helpers
func (r *Repo) GetLesson(lessonID string) (*content.Lesson, error) {
	var l content.Lesson
	if err := r.db.First(&l, "id=?", lessonID).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repo) GetPrevLessonSameModule(moduleID string, seq int) (*content.Lesson, error) {
	if seq <= 1 {
		return nil, gorm.ErrRecordNotFound
	}
	var l content.Lesson
	err := r.db.Where("module_id=? AND seq=?", moduleID, seq-1).First(&l).Error
	return &l, err
}

// สำหรับคำนวณ % ของคอร์ส
func (r *Repo) CountMandatoryLessonsOfCourse(courseID string) (int64, error) {
	var count int64
	err := r.db.Model(&content.Lesson{}).
		Joins("JOIN course_modules m ON lessons.module_id = m.id").
		Where("m.course_id = ? AND lessons.is_mandatory = 1", courseID).
		Count(&count).Error
	return count, err
}

func (r *Repo) CountCompletedMandatoryLessons(userID, courseID string) (int64, error) {
	var count int64
	err := r.db.Model(&learn.UserLessonProgress{}).
		Joins("JOIN lessons l ON user_lesson_progress.lesson_id = l.id").
		Joins("JOIN course_modules m ON l.module_id = m.id").
		Where("user_lesson_progress.user_id = ? AND m.course_id = ? AND l.is_mandatory = 1 AND user_lesson_progress.completed_at IS NOT NULL",
			userID, courseID).
		Count(&count).Error
	return count, err
}
