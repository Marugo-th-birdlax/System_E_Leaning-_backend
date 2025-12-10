package dto

type CreateCategoryReq struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateCategoryReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
