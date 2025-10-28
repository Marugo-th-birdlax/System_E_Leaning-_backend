package dto

type CourseResp struct {
	ID               string  `json:"id"`
	Code             string  `json:"code"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	IsActive         bool    `json:"is_active"`
	EstimatedMinutes *int    `json:"estimated_minutes,omitempty"`
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
