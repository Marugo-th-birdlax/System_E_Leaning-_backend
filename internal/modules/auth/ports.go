package auth

import "context"

type Service interface {
	Login(ctx context.Context, employeeCode, password string) (access, refresh string, err error)
	Refresh(ctx context.Context, refreshToken string) (access, newRefresh string, err error)
	Logout(ctx context.Context, refreshToken string) error
}

type Repository interface {
	SaveRefresh(ctx context.Context, userID, jti string, expiresAt int64) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
	Revoke(ctx context.Context, jti string) error
}
