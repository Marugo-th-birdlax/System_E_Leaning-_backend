package repo

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Marugo/birdlax/internal/modules/auth"
)

type RefreshToken struct {
	JTI       string `gorm:"primaryKey;size:64"`
	UserID    string `gorm:"size:36;index"`
	ExpiresAt int64  `gorm:"index"`
	Revoked   bool   `gorm:"index"`
	CreatedAt time.Time
}

type gormRepo struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) auth.Repository { return &gormRepo{db: db} }

func (r *gormRepo) SaveRefresh(ctx context.Context, userID, jti string, expiresAt int64) error {
	rt := &RefreshToken{JTI: jti, UserID: userID, ExpiresAt: expiresAt, Revoked: false}
	return r.db.WithContext(ctx).Create(rt).Error
}

func (r *gormRepo) IsRevoked(ctx context.Context, jti string) (bool, error) {
	var rt RefreshToken
	if err := r.db.WithContext(ctx).First(&rt, "jti = ?", jti).Error; err != nil {
		return true, err // หาไม่เจอ=ถือว่าใช้ไม่ได้
	}
	if rt.Revoked || rt.ExpiresAt <= time.Now().Unix() {
		return true, nil
	}
	return false, nil
}

func (r *gormRepo) Revoke(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("jti = ?", jti).
		Update("revoked", true).Error
}
