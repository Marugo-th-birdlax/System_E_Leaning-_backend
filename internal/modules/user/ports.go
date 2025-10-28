package user

import (
	"context"

	usermodels "github.com/Marugo/birdlax/internal/modules/user/models"
)

type Repository interface {
	FindAll(ctx context.Context, limit, offset int, q string) ([]usermodels.User, int64, error)
	FindByID(ctx context.Context, id string) (*usermodels.User, error)
	Create(ctx context.Context, u *usermodels.User) error

	FindByEmployeeCode(ctx context.Context, employeeCode string) (*usermodels.User, error)

	Update(ctx context.Context, u *usermodels.User) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	List(ctx context.Context, page, perPage int, q string) ([]usermodels.User, int64, error)
	Get(ctx context.Context, id string) (*usermodels.User, error)
	Create(ctx context.Context, employeeCode, email, firstName, lastName, role string, phone *string, password string) (*usermodels.User, error)
	Update(ctx context.Context, id string, employeeCode, email, firstName, lastName, role *string, phone *string, isActive *bool, password *string) (*usermodels.User, error)

	Delete(ctx context.Context, id string) error
}
