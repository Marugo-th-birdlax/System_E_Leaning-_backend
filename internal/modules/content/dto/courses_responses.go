package dto

import "github.com/Marugo/birdlax/internal/modules/content/models"

type CourseResp struct {
	ID               string   `json:"id"`
	Code             string   `json:"code"`
	Title            string   `json:"title"`
	Description      *string  `json:"description,omitempty"`
	IsActive         bool     `json:"is_active"`
	EstimatedMinutes *int     `json:"estimated_minutes,omitempty"`
	CategoryID       *string  `json:"category_id"`
	DepartmentIDs    []string `json:"department_ids,omitempty"`
	CreatedAt        string   `json:"CreatedAt"`
	UpdatedAt        string   `json:"UpdatedAt"`
}

type ModuleResp struct {
	ID          string  `json:"id"`
	CourseID    string  `json:"course_id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Seq         int     `json:"seq"`
	IsMandatory bool    `json:"is_mandatory"`
}

type PageMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

type PagedCourses struct {
	Data []CourseResp `json:"data"`
	Meta PageMeta     `json:"meta"`
}

func FromCourseModel(c *models.Course, deptIDs []string) *CourseResp {
	if c == nil {
		return nil
	}
	return &CourseResp{
		ID:               c.ID,
		Code:             c.Code,
		Title:            c.Title,
		Description:      c.Description,
		IsActive:         c.IsActive,
		EstimatedMinutes: c.EstimatedMinutes,
		CategoryID:       c.CategoryID,
		DepartmentIDs:    deptIDs,
		CreatedAt:        c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
