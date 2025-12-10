package service

import (
	"errors"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/models"
	"github.com/google/uuid"
)

type CourseRepo interface {
	Create(*models.Course) error
	Update(*models.Course) error
	Delete(id string) error
	GetByID(id string) (*models.Course, error)
	List(q string, page, per int) ([]models.Course, int64, error)
	ListByCategory(categoryID, q string, page, per int) ([]models.Course, int64, error)
}

type CourseDeptRepo interface {
	ReplaceTargets(courseID string, deptIDs []string) error
	ListDepartmentIDs(courseID string) ([]string, error)
}

type ModuleRepo interface {
	Create(*models.CourseModule) error
	Update(*models.CourseModule) error
	Delete(id string) error
	GetByID(id string) (*models.CourseModule, error)
	ListByCourse(courseID string) ([]models.CourseModule, error)
}
type LessonLister interface {
	// มีอยู่แล้วใน content repo เดิม: ดึงบทเรียนของโมดูล
	GetLessonsByModule(moduleID string) ([]models.Lesson, error)
}

type CourseService interface {
	CreateCourse(req dto.CreateCourseReq) (*models.Course, error)
	UpdateCourse(id string, req dto.UpdateCourseReq) (*models.Course, error)
	DeleteCourse(id string) error
	GetCourse(id string) (*models.Course, error)
	ListCourses(q string, page, per int) ([]models.Course, int64, error)

	ListCourseDepartments(courseID string) ([]string, error)

	CreateModule(courseID string, req dto.CreateModuleReq) (*models.CourseModule, error)
	UpdateModule(id string, req dto.UpdateModuleReq) (*models.CourseModule, error)
	DeleteModule(id string) error
	ListModules(courseID string) ([]models.CourseModule, error)
	ListLessons(moduleID string) ([]models.Lesson, error)
}

type courseSvc struct {
	courseRepo CourseRepo
	moduleRepo ModuleRepo
	lessonList LessonLister
	catRepo    CategoryRepo
	deptRepo   CourseDeptRepo
}

func NewCourseService(cr CourseRepo, mr ModuleRepo, ll LessonLister, cats CategoryRepo, dr CourseDeptRepo) CourseService {
	return &courseSvc{courseRepo: cr, moduleRepo: mr, lessonList: ll, catRepo: cats, deptRepo: dr}
}

/******** Courses ********/
/******** Courses ********/
func (s *courseSvc) CreateCourse(req dto.CreateCourseReq) (*models.Course, error) {
	if req.Code == "" || req.Title == "" {
		return nil, errors.New("code and title required")
	}

	// ✅ validate category_id ถ้ามี
	if req.CategoryID != nil && *req.CategoryID != "" {
		if s.catRepo == nil {
			return nil, errors.New("category repo not wired")
		}
		ok, err := s.catRepo.Exists(*req.CategoryID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("invalid category_id")
		}
	}

	c := &models.Course{
		ID:               uuid.NewString(),
		Code:             req.Code,
		Title:            req.Title,
		Description:      req.Description,
		IsActive:         true,
		EstimatedMinutes: req.EstimatedMinutes,
		CategoryID:       req.CategoryID,
	}
	if req.IsActive != nil {
		c.IsActive = *req.IsActive
	}

	// ⛳ 1) เซฟตัว course ก่อน
	if err := s.courseRepo.Create(c); err != nil {
		return nil, err
	}

	// ⛳ 2) ผูก department targets ถ้ามีส่งมา และมี repo
	if s.deptRepo != nil && len(req.DepartmentIDs) > 0 {
		if err := s.deptRepo.ReplaceTargets(c.ID, req.DepartmentIDs); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (s *courseSvc) UpdateCourse(id string, req dto.UpdateCourseReq) (*models.Course, error) {
	c, err := s.courseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		c.Title = *req.Title
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.IsActive != nil {
		c.IsActive = *req.IsActive
	}
	if req.EstimatedMinutes != nil {
		c.EstimatedMinutes = req.EstimatedMinutes
	}
	if req.CategoryID != nil {
		if *req.CategoryID == "" {
			c.CategoryID = nil // เคลียร์หมวด
		} else {
			if s.catRepo == nil {
				return nil, errors.New("category repo not wired")
			}
			ok, err := s.catRepo.Exists(*req.CategoryID)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, errors.New("invalid category_id")
			}
			c.CategoryID = req.CategoryID
		}
	}

	// ⛳ 1) อัปเดต course ก่อน
	if err := s.courseRepo.Update(c); err != nil {
		return nil, err
	}

	// ⛳ 2) ถ้า client ส่ง department_ids มา → แทนที่ mapping ทั้งชุด
	// ใช้ pointer (*[]string) เพื่อแยกกรณี:
	//   - nil  = ไม่แตะเรื่อง department
	//   - []{} = เคลียร์ทุก department
	if s.deptRepo != nil && req.DepartmentIDs != nil {
		if err := s.deptRepo.ReplaceTargets(c.ID, *req.DepartmentIDs); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (s *courseSvc) DeleteCourse(id string) error                { return s.courseRepo.Delete(id) }
func (s *courseSvc) GetCourse(id string) (*models.Course, error) { return s.courseRepo.GetByID(id) }
func (s *courseSvc) ListCourses(q string, page, per int) ([]models.Course, int64, error) {
	return s.courseRepo.List(q, page, per)
}

func (s *courseSvc) ListCourseDepartments(courseID string) ([]string, error) {
	if s.deptRepo == nil {
		return []string{}, nil
	}
	return s.deptRepo.ListDepartmentIDs(courseID)
}

/******** Modules ********/
func (s *courseSvc) CreateModule(courseID string, req dto.CreateModuleReq) (*models.CourseModule, error) {
	if req.Title == "" || req.Seq < 1 {
		return nil, errors.New("title and seq required")
	}
	m := &models.CourseModule{
		ID:          uuid.NewString(),
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		Seq:         req.Seq,
		IsMandatory: true,
	}
	if req.IsMandatory != nil {
		m.IsMandatory = *req.IsMandatory
	}
	return m, s.moduleRepo.Create(m)
}
func (s *courseSvc) UpdateModule(id string, req dto.UpdateModuleReq) (*models.CourseModule, error) {
	m, err := s.moduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.Seq != nil {
		m.Seq = *req.Seq
	}
	if req.IsMandatory != nil {
		m.IsMandatory = *req.IsMandatory
	}
	return m, s.moduleRepo.Update(m)
}
func (s *courseSvc) DeleteModule(id string) error { return s.moduleRepo.Delete(id) }
func (s *courseSvc) ListModules(courseID string) ([]models.CourseModule, error) {
	return s.moduleRepo.ListByCourse(courseID)
}
func (s *courseSvc) ListLessons(moduleID string) ([]models.Lesson, error) {
	return s.lessonList.GetLessonsByModule(moduleID)
}
