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

// helper
func normalizeRole(s string) usermodels.Role {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "admin":
		return usermodels.RoleAdmin
	case "planning":
		return usermodels.RolePlanning
	case "production_control":
		return usermodels.RolePordcutionControl
	case "production":
		return usermodels.RoleProduction
	case "outsourcing":
		return usermodels.RoleOutsourcing
	case "production_adh":
		return usermodels.RoleProductionADH
	case "production_weld":
		return usermodels.RoleProductionWeld
	case "production_assy":
		return usermodels.RoleProductionAssy
	default:
		return usermodels.RoleUser
	}
}
func isValidRole(r usermodels.Role) bool {
	for _, v := range usermodels.AllRoles() {
		if r == v {
			return true
		}
	}
	return false
}

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

	r := normalizeRole(role)
	if !isValidRole(r) {
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
		r := normalizeRole(*role)
		if !isValidRole(r) {
			return nil, errors.New("invalid role")
		}
		u.Role = r
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
