package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/user"
	usermodels "github.com/Marugo/birdlax/internal/modules/user/models"
)

type gormRepo struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) user.Repository { return &gormRepo{db: db} }

func (r *gormRepo) FindAll(ctx context.Context, limit, offset int, q string) ([]usermodels.User, int64, error) {
	var rows []usermodels.User
	tx := r.db.WithContext(ctx).Model(&usermodels.User{})
	if q != "" {
		tx = tx.Where(`
		employee_code LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR
		CONCAT(first_name,' ',last_name) LIKE ?`,
			"%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%",
		)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *gormRepo) FindByID(ctx context.Context, id string) (*usermodels.User, error) {
	var u usermodels.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormRepo) FindByEmployeeCode(ctx context.Context, code string) (*usermodels.User, error) {
	var u usermodels.User
	if err := r.db.WithContext(ctx).First(&u, "employee_code = ?", code).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormRepo) Create(ctx context.Context, u *usermodels.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}
func (r *gormRepo) Update(ctx context.Context, u *usermodels.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}
func (r *gormRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&usermodels.User{}, "id = ?", id).Error
}
