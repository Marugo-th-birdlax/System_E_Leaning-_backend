package user

import (
	"context"

	usermodels "github.com/Marugo/birdlax/internal/modules/user/models"
)

type Repository interface {
	// Users
	FindAll(ctx context.Context, limit, offset int, q string) ([]usermodels.User, int64, error)
	FindByID(ctx context.Context, id string) (*usermodels.User, error)
	FindByEmployeeCode(ctx context.Context, employeeCode string) (*usermodels.User, error)
	Create(ctx context.Context, u *usermodels.User) error
	Update(ctx context.Context, u *usermodels.User) error
	Delete(ctx context.Context, id string) error

	// Departments
	ListDepartments(ctx context.Context, q string) ([]usermodels.Department, error)
	GetDepartmentByID(ctx context.Context, id string) (*usermodels.Department, error)
	CreateDepartment(ctx context.Context, d *usermodels.Department) error
	UpdateDepartment(ctx context.Context, d *usermodels.Department) error
	DeleteDepartment(ctx context.Context, id string) error

	// UserDepartmentRoles
	ListUserDepartmentRoles(ctx context.Context, userID string) ([]usermodels.UserDepartmentRole, error)
	AddUserDepartmentRole(ctx context.Context, r *usermodels.UserDepartmentRole) error
	DeleteUserDepartmentRole(ctx context.Context, relID string) error
	ListDepartmentManagers(ctx context.Context, departmentID string) ([]usermodels.UserDepartmentRole, error)
}

type Service interface {
	// Users
	List(ctx context.Context, page, perPage int, q string) ([]usermodels.User, int64, error)
	Get(ctx context.Context, id string) (*usermodels.User, error)
	Create(ctx context.Context, employeeCode, email, firstName, lastName, role string, phone *string, password string) (*usermodels.User, error)
	Update(ctx context.Context, id string,
		employeeCode, email, firstName, lastName, role *string,
		phone *string, isActive *bool, password *string) (*usermodels.User, error)
	Delete(ctx context.Context, id string) error

	// Departments
	ListDepartments(ctx context.Context, q string) ([]usermodels.Department, error)
	GetDepartment(ctx context.Context, id string) (*usermodels.Department, error)
	CreateDepartment(ctx context.Context, code, name string, isActive *bool) (*usermodels.Department, error)
	UpdateDepartment(ctx context.Context, id string, code, name *string, isActive *bool) (*usermodels.Department, error)
	DeleteDepartment(ctx context.Context, id string) error

	// User department roles
	ListUserDepartmentRoles(ctx context.Context, userID string) ([]usermodels.UserDepartmentRole, error)
	AddUserDepartmentRole(ctx context.Context, userID, departmentID, role string) (*usermodels.UserDepartmentRole, error)
	DeleteUserDepartmentRole(ctx context.Context, relID string) error
	ListDepartmentManagers(ctx context.Context, departmentID string) ([]usermodels.UserDepartmentRole, error)
}
