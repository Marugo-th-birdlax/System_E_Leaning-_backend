package dto

type CreateCourseReq struct {
	Code             string  `json:"code" validate:"required"`
	Title            string  `json:"title" validate:"required"`
	Description      *string `json:"description"`
	IsActive         *bool   `json:"is_active"`
	EstimatedMinutes *int    `json:"estimated_minutes"`
}

type UpdateCourseReq struct {
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	IsActive         *bool   `json:"is_active"`
	EstimatedMinutes *int    `json:"estimated_minutes"`
}

type CreateModuleReq struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Seq         int     `json:"seq" validate:"required,min=1"`
	IsMandatory *bool   `json:"is_mandatory"`
}

type UpdateModuleReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Seq         *int    `json:"seq"`
	IsMandatory *bool   `json:"is_mandatory"`
}
