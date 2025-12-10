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

// ===== Departments =====

type DepartmentResponse struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func FromDepartmentModel(d *usermodels.Department) *DepartmentResponse {
	if d == nil {
		return nil
	}
	return &DepartmentResponse{
		ID:       d.ID,
		Code:     d.Code,
		Name:     d.Name,
		IsActive: d.IsActive,
	}
}

// ===== User Department Roles =====

type UserDepartmentRoleResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	DepartmentID string `json:"department_id"`
	Role         string `json:"role"`
	// optional: แถมชื่อแผนก
	DepartmentName string `json:"department_name,omitempty"`
}

func FromUserDepartmentRoleModel(m *usermodels.UserDepartmentRole) *UserDepartmentRoleResponse {
	if m == nil {
		return nil
	}
	resp := &UserDepartmentRoleResponse{
		ID:           m.ID,
		UserID:       m.UserID,
		DepartmentID: m.DepartmentID,
		Role:         string(m.Role),
	}
	if m.Department != nil {
		resp.DepartmentName = m.Department.Name
	}
	return resp
}
