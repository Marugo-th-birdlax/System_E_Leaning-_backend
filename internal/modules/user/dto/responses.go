package dto

import usermodels "github.com/Marugo/birdlax/internal/modules/user/models"

type UserResponse struct {
	ID           string  `json:"id"`
	EmployeeCode string  `json:"employee_code"`
	Email        string  `json:"email"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Role         string  `json:"role"`
	Phone        *string `json:"phone"`
	IsActive     bool    `json:"is_active"`
	FullName     string  `json:"full_name"`
}

func FromModel(u *usermodels.User) *UserResponse {
	if u == nil {
		return nil
	}
	fn, ln := u.FirstName, u.LastName
	full := fn
	if ln != "" {
		if full != "" {
			full += " "
		}
		full += ln
	}
	return &UserResponse{
		ID:           u.ID,
		EmployeeCode: u.EmployeeCode,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Role:         string(u.Role),
		Phone:        u.Phone,
		IsActive:     u.IsActive,
		FullName:     full,
	}
}
