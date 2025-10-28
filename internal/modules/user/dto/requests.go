package dto

type UserCreateRequest struct {
	EmployeeCode string  `json:"employee_code" validate:"required,max=30"`
	Email        string  `json:"email"         validate:"required,email,max=190"`
	FirstName    string  `json:"first_name"    validate:"required,max=100"`
	LastName     string  `json:"last_name"     validate:"required,max=100"`
	Role         string  `json:"role" validate:"required,oneof=user admin planning production_control production outsourcing production_adh production_weld production_assy"`
	Phone        *string `json:"phone"         validate:"omitempty,max=50"`
	Password     string  `json:"password"      validate:"required,min=8"` // <-- เพิ่ม
	// PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"` // (ถ้าต้องการยืนยัน)
}

type UserUpdateRequest struct {
	EmployeeCode *string `json:"employee_code" validate:"omitempty,max=30"`
	Email        *string `json:"email"         validate:"omitempty,email,max=190"`
	FirstName    *string `json:"first_name"    validate:"omitempty,max=100"`
	LastName     *string `json:"last_name"     validate:"omitempty,max=100"`
	Role         *string `json:"role"          validate:"omitempty,oneof=user admin planning production_control production outsourcing production_adh production_weld production_assy"`
	Phone        *string `json:"phone"         validate:"omitempty,max=50"`
	IsActive     *bool   `json:"is_active"     validate:"omitempty"`
	Password     *string `json:"password"      validate:"omitempty,min=8"`
}
