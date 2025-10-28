package service

import (
	content "github.com/Marugo/birdlax/internal/modules/content/models"
	"github.com/Marugo/birdlax/internal/modules/learning/dto"
	"github.com/Marugo/birdlax/internal/modules/learning/models"
)

type Repo interface {
	UpsertEnrollment(e *models.Enrollment) error
	GetEnrollment(userID, courseID string) (*models.Enrollment, error)

	GetLessonProgress(userID, lessonID string) (*models.UserLessonProgress, error)
	CreateLessonProgress(p *models.UserLessonProgress) error
	UpdateLessonProgress(p *models.UserLessonProgress) error

	GetLesson(lessonID string) (*content.Lesson, error)
	GetPrevLessonSameModule(moduleID string, seq int) (*content.Lesson, error)

	CountMandatoryLessonsOfCourse(courseID string) (int64, error)
	CountCompletedMandatoryLessons(userID, courseID string) (int64, error)
}

type Service interface {
	EnrollCourse(userID, courseID string) (*models.Enrollment, error)
	GetEnrollment(userID, courseID string) (*models.Enrollment, error)

	StartLesson(userID, lessonID string, req dto.StartLessonReq) (*models.UserLessonProgress, error)
	TrackLesson(userID, lessonID string, req dto.TrackLessonReq) (*models.UserLessonProgress, error)
	CompleteLesson(userID, lessonID string, req dto.CompleteLessonReq) (*models.UserLessonProgress, error)
	UpdateEnrollmentPercent(userID, courseID string) error
}
