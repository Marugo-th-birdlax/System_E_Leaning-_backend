package service

import (
	"context"

	contentmodels "github.com/Marugo/birdlax/internal/modules/content/models"
	learningrepo "github.com/Marugo/birdlax/internal/modules/learning/repo"
)

type MyCoursesService interface {
	ListDepartmentCourses(ctx context.Context, userID string, categoryID *string, page, per int) ([]contentmodels.Course, int64, error)
	ListMyCourses(ctx context.Context, userID string, page, per int) ([]learningrepo.MyEnrolledCourse, int64, error)
	GetCourseProgress(ctx context.Context, userID, courseID string) (*learningrepo.CourseProgress, error)
}

type myCoursesSvc struct {
	repo *learningrepo.MyCoursesRepo
}

func NewMyCoursesService(r *learningrepo.MyCoursesRepo) MyCoursesService {
	return &myCoursesSvc{repo: r}
}

func (s *myCoursesSvc) ListDepartmentCourses(
	ctx context.Context,
	userID string,
	categoryID *string,
	page, per int,
) ([]contentmodels.Course, int64, error) {
	if per <= 0 || per > 200 {
		per = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * per
	return s.repo.ListDepartmentCourses(ctx, userID, categoryID, per, offset)
}

func (s *myCoursesSvc) ListMyCourses(ctx context.Context, userID string, page, per int) ([]learningrepo.MyEnrolledCourse, int64, error) {
	if per <= 0 || per > 200 {
		per = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * per
	return s.repo.ListMyEnrolledCourses(ctx, userID, per, offset)
}

func (s *myCoursesSvc) GetCourseProgress(ctx context.Context, userID, courseID string) (*learningrepo.CourseProgress, error) {
	return s.repo.GetCourseProgress(ctx, userID, courseID)
}
