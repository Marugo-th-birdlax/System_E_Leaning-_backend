package dto

type UserCreateRequest struct {
	EmployeeCode string  `json:"employee_code" validate:"required"`
	Email        string  `json:"email"         validate:"required,email"`
	FirstName    string  `json:"first_name"    validate:"required"`
	LastName     string  `json:"last_name"     validate:"required"`
	Role         string  `json:"role"          validate:"omitempty,oneof=admin employee hr"`
	Phone        *string `json:"phone"`
	Password     string  `json:"password"      validate:"required,min=6"`
}

type UserUpdateRequest struct {
	EmployeeCode *string `json:"employee_code" validate:"omitempty"`
	Email        *string `json:"email"         validate:"omitempty,email"`
	FirstName    *string `json:"first_name"    validate:"omitempty"`
	LastName     *string `json:"last_name"     validate:"omitempty"`
	Role         *string `json:"role"          validate:"omitempty,oneof=admin employee hr"`
	Phone        *string `json:"phone"`
	IsActive     *bool   `json:"is_active"`
	Password     *string `json:"password"      validate:"omitempty,min=6"`
}

// ===== Departments =====

type DepartmentCreateRequest struct {
	Code     string `json:"code"     validate:"required"`
	Name     string `json:"name"     validate:"required"`
	IsActive *bool  `json:"is_active"`
}

type DepartmentUpdateRequest struct {
	Code     *string `json:"code"`
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

// ===== User Department Roles =====

type AssignDepartmentRoleRequest struct {
	DepartmentID string `json:"department_id" validate:"required"`
	Role         string `json:"role"          validate:"required,oneof=manager leader member"`
}
