package service

import (
	"errors"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/models"
	"github.com/google/uuid"
)

type CategoryRepo interface {
	Exists(id string) (bool, error)
	Create(*models.Category) error
	Update(*models.Category) error
	Delete(id string) error
	GetByID(id string) (*models.Category, error)
	List(q string, page, per int) ([]models.Category, int64, error)
}

type CategoryService interface {
	CreateCategory(req dto.CreateCategoryReq) (*models.Category, error)
	UpdateCategory(id string, req dto.UpdateCategoryReq) (*models.Category, error)
	DeleteCategory(id string) error

	GetCategory(id string) (*models.Category, error)
	ListCategories(q string, page, per int) ([]models.Category, int64, error)

	// สำหรับ /categories/:id/courses
	ListCoursesOfCategory(id, q string, page, per int) ([]models.Course, int64, error)
}

type categorySvc struct {
	catRepo    CategoryRepo
	courseRepo CourseRepo // ใช้ method ListByCategory(...)
}

func NewCategoryService(cr CategoryRepo, courseRepo CourseRepo) CategoryService {
	return &categorySvc{catRepo: cr, courseRepo: courseRepo}
}

func (s *categorySvc) CreateCategory(req dto.CreateCategoryReq) (*models.Category, error) {
	if req.Code == "" || req.Title == "" {
		return nil, errors.New("code and title required")
	}
	c := &models.Category{
		ID:          uuid.NewString(),
		Code:        req.Code,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    true,
	}
	if req.IsActive != nil {
		c.IsActive = *req.IsActive
	}
	return c, s.catRepo.Create(c)
}

func (s *categorySvc) UpdateCategory(id string, req dto.UpdateCategoryReq) (*models.Category, error) {
	c, err := s.catRepo.GetByID(id)
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
	return c, s.catRepo.Update(c)
}

func (s *categorySvc) DeleteCategory(id string) error { return s.catRepo.Delete(id) }

func (s *categorySvc) GetCategory(id string) (*models.Category, error) { return s.catRepo.GetByID(id) }

func (s *categorySvc) ListCategories(q string, page, per int) ([]models.Category, int64, error) {
	return s.catRepo.List(q, page, per)
}

func (s *categorySvc) ListCoursesOfCategory(id, q string, page, per int) ([]models.Course, int64, error) {
	// optional: validate ว่าหมวดมีจริงก่อน
	ok, err := s.catRepo.Exists(id)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, errors.New("category not found")
	}
	return s.courseRepo.ListByCategory(id, q, page, per)
}
