package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Marugo/birdlax/internal/modules/user"
	usermodels "github.com/Marugo/birdlax/internal/modules/user/models"
	"github.com/Marugo/birdlax/internal/shared/password"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrEmployeeCodeAlreadyExist = errors.New("employee_code already exists")
)

type svc struct{ repo user.Repository }

func NewService(r user.Repository) user.Service { return &svc{repo: r} }

func (s *svc) List(ctx context.Context, page, perPage int, q string) ([]usermodels.User, int64, error) {
	if perPage <= 0 || perPage > 200 {
		perPage = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage
	return s.repo.FindAll(ctx, perPage, offset, q)
}

func (s *svc) Get(ctx context.Context, id string) (*usermodels.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *svc) Create(ctx context.Context,
	employeeCode, email, firstName, lastName, role string, phone *string, passwordPlain string) (*usermodels.User, error) {
	employeeCode = strings.TrimSpace(employeeCode)
	email = strings.TrimSpace(strings.ToLower(email))
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if employeeCode == "" || email == "" || firstName == "" || lastName == "" {
		return nil, errors.New("employee_code, email, first_name, last_name are required")
	}

	r, ok := usermodels.ParseRole(role)
	if !ok {
		return nil, errors.New("invalid role")
	}
	if strings.TrimSpace(passwordPlain) == "" {
		return nil, errors.New("password is required")
	}
	hash, err := password.Hash(passwordPlain)
	if err != nil {
		return nil, err
	}

	u := &usermodels.User{
		EmployeeCode: employeeCode,
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         r,
		Phone:        phone,
		IsActive:     true,
		PasswordHash: hash, // <-- เก็บ hash
	}
	if err := s.repo.Create(ctx, u); err != nil {
		// GORM จะ map duplicate เป็น ErrDuplicatedKey
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// ตัดสินใจไม่ได้ว่า key ไหนซ้ำ? ให้ repo โยน error เฉพาะมาก็ได้
			// ที่นี่ fallback ตีความจากข้อความ
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "employee_code") {
				return nil, ErrEmployeeCodeAlreadyExist
			}
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}
	return u, nil
}

func (s *svc) Update(ctx context.Context,
	id string,
	employeeCode, email, firstName, lastName, role *string, phone *string, isActive *bool,
	// เพิ่มรับรหัสผ่านใหม่ (optional)
	passwordPlain *string,
) (*usermodels.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if employeeCode != nil {
		ec := strings.TrimSpace(*employeeCode)
		if ec == "" {
			return nil, errors.New("employee_code cannot be empty")
		}
		u.EmployeeCode = ec
	}
	if email != nil {
		e := strings.TrimSpace(strings.ToLower(*email))
		if e == "" {
			return nil, errors.New("email cannot be empty")
		}
		u.Email = e
	}
	if firstName != nil {
		fn := strings.TrimSpace(*firstName)
		if fn == "" {
			return nil, errors.New("first_name cannot be empty")
		}
		u.FirstName = fn
	}
	if lastName != nil {
		ln := strings.TrimSpace(*lastName)
		if ln == "" {
			return nil, errors.New("last_name cannot be empty")
		}
		u.LastName = ln
	}
	if role != nil {
		if r, ok := usermodels.ParseRole(*role); ok {
			u.Role = r
		} else {
			return nil, errors.New("invalid role")
		}
	}

	if passwordPlain != nil {
		p := strings.TrimSpace(*passwordPlain)
		if p == "" {
			return nil, errors.New("password cannot be empty")
		}
		hash, err := password.Hash(p)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = hash // <-- อัปเดต hash
	}
	if phone != nil {
		u.Phone = phone
	}
	if isActive != nil {
		u.IsActive = *isActive
	}

	if err := s.repo.Update(ctx, u); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "employee_code") {
				return nil, ErrEmployeeCodeAlreadyExist
			}
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}
	return u, nil
}

func (s *svc) Delete(ctx context.Context, id string) error { return s.repo.Delete(ctx, id) }

func (s *svc) ListDepartments(ctx context.Context, q string) ([]usermodels.Department, error) {
	return s.repo.ListDepartments(ctx, q)
}

func (s *svc) GetDepartment(ctx context.Context, id string) (*usermodels.Department, error) {
	return s.repo.GetDepartmentByID(ctx, id)
}

func (s *svc) CreateDepartment(ctx context.Context, code, name string, isActive *bool) (*usermodels.Department, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	if code == "" || name == "" {
		return nil, errors.New("code and name are required")
	}
	d := &usermodels.Department{
		Code:     code,
		Name:     name,
		IsActive: true,
	}
	if isActive != nil {
		d.IsActive = *isActive
	}
	if err := s.repo.CreateDepartment(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *svc) UpdateDepartment(ctx context.Context, id string, code, name *string, isActive *bool) (*usermodels.Department, error) {
	d, err := s.repo.GetDepartmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if code != nil {
		c := strings.TrimSpace(*code)
		if c == "" {
			return nil, errors.New("code cannot be empty")
		}
		d.Code = c
	}
	if name != nil {
		n := strings.TrimSpace(*name)
		if n == "" {
			return nil, errors.New("name cannot be empty")
		}
		d.Name = n
	}
	if isActive != nil {
		d.IsActive = *isActive
	}
	if err := s.repo.UpdateDepartment(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *svc) DeleteDepartment(ctx context.Context, id string) error {
	return s.repo.DeleteDepartment(ctx, id)
}

// ===== User Department Roles =====

func (s *svc) ListUserDepartmentRoles(ctx context.Context, userID string) ([]usermodels.UserDepartmentRole, error) {
	// optional: validate user exists
	if _, err := s.repo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	return s.repo.ListUserDepartmentRoles(ctx, userID)
}

func (s *svc) AddUserDepartmentRole(ctx context.Context, userID, departmentID, role string) (*usermodels.UserDepartmentRole, error) {
	// validate user / dept
	if _, err := s.repo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetDepartmentByID(ctx, departmentID); err != nil {
		return nil, err
	}
	r, ok := usermodels.ParseDepartmentRole(role)
	if !ok {
		return nil, errors.New("invalid department role")
	}
	rel := &usermodels.UserDepartmentRole{
		UserID:       userID,
		DepartmentID: departmentID,
		Role:         r,
	}
	if err := s.repo.AddUserDepartmentRole(ctx, rel); err != nil {
		return nil, err
	}
	return rel, nil
}

func (s *svc) DeleteUserDepartmentRole(ctx context.Context, relID string) error {
	return s.repo.DeleteUserDepartmentRole(ctx, relID)
}

func (s *svc) ListDepartmentManagers(ctx context.Context, departmentID string) ([]usermodels.UserDepartmentRole, error) {
	// optional: validate department exists
	if _, err := s.repo.GetDepartmentByID(ctx, departmentID); err != nil {
		return nil, err
	}
	return s.repo.ListDepartmentManagers(ctx, departmentID)
}
