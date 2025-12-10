package repo

import (
	"context"
	"time"

	contentmodels "github.com/Marugo/birdlax/internal/modules/content/models"
	"gorm.io/gorm"
)

type MyCoursesRepo struct {
	db *gorm.DB
}

func NewMyCoursesRepo(db *gorm.DB) *MyCoursesRepo {
	return &MyCoursesRepo{db: db}
}

// คอร์สตามแผนกของ user (ยังไม่สน enroll)
func (r *MyCoursesRepo) ListDepartmentCourses(
	ctx context.Context,
	userID string,
	categoryID *string,
	limit, offset int,
) ([]contentmodels.Course, int64, error) {

	var courses []contentmodels.Course

	tx := r.db.WithContext(ctx).
		Model(&contentmodels.Course{}).
		Joins(`JOIN course_department_targets t ON t.course_id = courses.id`).
		Joins(`JOIN user_department_roles udr ON udr.department_id = t.department_id`).
		Where("udr.user_id = ?", userID).
		Where("courses.is_active = 1").
		Where("courses.deleted_at IS NULL").
		Where("udr.deleted_at IS NULL").
		Group("courses.id")

	// 👇 ถ้ามี category_id ให้ filter เพิ่ม
	if categoryID != nil && *categoryID != "" {
		tx = tx.Where("courses.category_id = ?", *categoryID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := tx.
		Order("courses.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// คอร์สที่ user ลงทะเบียนแล้ว (จาก enrollments)
type MyEnrolledCourse struct {
	contentmodels.Course
	EnrollmentID    string  `json:"enrollment_id"`
	ProgressPercent float64 `json:"progress_percent"`
	Status          string  `json:"status"`
}

func (r *MyCoursesRepo) ListMyEnrolledCourses(ctx context.Context, userID string, limit, offset int) ([]MyEnrolledCourse, int64, error) {
	var rows []MyEnrolledCourse

	tx := r.db.WithContext(ctx).
		Table("enrollments AS e").
		Joins(`JOIN courses c ON c.id = e.course_id`).
		Where("e.user_id = ?", userID).
		Where("e.deleted_at IS NULL").
		Where("c.deleted_at IS NULL")

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.
		Select(`c.*, e.id AS enrollment_id, e.progress_percent, e.status`).
		Limit(limit).Offset(offset).
		Order("e.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

type LessonProgressItem struct {
	LessonID        string     `json:"lesson_id"`
	Title           string     `json:"title"`
	ModuleID        string     `json:"module_id"`
	Seq             int        `json:"seq"`
	ContentType     string     `json:"content_type"`
	ProgressPercent float64    `json:"progress_percent"`
	CurrentPosition int64      `json:"current_position"`
	MaxPosition     int64      `json:"max_position"`
	IsUnlocked      bool       `json:"is_unlocked"`
	StartedAt       *time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

// CourseProgress aggregation
type CourseProgress struct {
	Enrollment *struct {
		ID              string  `json:"id"`
		UserID          string  `json:"user_id"`
		CourseID        string  `json:"course_id"`
		Status          string  `json:"status"`
		ProgressPercent float64 `json:"progress_percent"`
	} `json:"enrollment"`
	Lessons []LessonProgressItem `json:"lessons"`
}

// GetCourseProgress returns enrollment (if any) and lesson-level progress for a course
func (r *MyCoursesRepo) GetCourseProgress(ctx context.Context, userID, courseID string) (*CourseProgress, error) {
	res := &CourseProgress{}

	// 1) try load enrollment
	var e struct {
		ID              string  `json:"id"`
		UserID          string  `json:"user_id"`
		CourseID        string  `json:"course_id"`
		Status          string  `json:"status"`
		ProgressPercent float64 `json:"progress_percent"`
	}
	if err := r.db.WithContext(ctx).
		Table("enrollments").
		Select("id, user_id, course_id, status, progress_percent").
		Where("user_id = ? AND course_id = ? AND deleted_at IS NULL", userID, courseID).
		First(&e).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// not found => leave Enrollment nil
	} else {
		res.Enrollment = &e
	}

	// 2) get lessons + left join user_lesson_progress (user-specific)
	var lessons []LessonProgressItem
	tx := r.db.WithContext(ctx).
		Table("lessons AS l").
		Select(`l.id AS lesson_id, l.title, l.module_id, l.seq, l.content_type,
			COALESCE(ulp.progress_percent, 0) AS progress_percent,
			COALESCE(ulp.current_position, 0) AS current_position,
			COALESCE(ulp.max_position, 0) AS max_position,
			COALESCE(ulp.is_unlocked, false) AS is_unlocked,
			ulp.started_at, ulp.completed_at`).
		Joins("JOIN course_modules m ON m.id = l.module_id").
		Joins("LEFT JOIN user_lesson_progress ulp ON ulp.lesson_id = l.id AND ulp.user_id = ?", userID).
		Where("m.course_id = ? AND l.deleted_at IS NULL", courseID).
		Order("l.seq ASC")

	if err := tx.Scan(&lessons).Error; err != nil {
		return nil, err
	}
	res.Lessons = lessons
	return res, nil
}
