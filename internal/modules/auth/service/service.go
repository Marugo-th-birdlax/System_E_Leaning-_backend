package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	auth "github.com/Marugo/birdlax/internal/modules/auth"
	user "github.com/Marugo/birdlax/internal/modules/user"
	"github.com/Marugo/birdlax/internal/shared/password"
	"github.com/Marugo/birdlax/internal/shared/security"
)

type svc struct {
	users  user.Repository // ต้องมี FindByEmployeeCode
	tokens auth.Repository // ต้องมี SaveRefresh/IsRevoked/Revoke
}

func New(users user.Repository, tokens auth.Repository) auth.Service { return &svc{users, tokens} }

func (s *svc) Login(ctx context.Context, employeeCode, pw string) (string, string, error) {
	u, err := s.users.FindByEmployeeCode(ctx, employeeCode) // ✅ เมธอดนี้ต้องมีแล้วใน user repo
	if err != nil || u == nil || !u.IsActive {
		return "", "", errors.New("invalid credentials")
	}
	if !password.Verify(u.PasswordHash, pw) {
		return "", "", errors.New("invalid credentials")
	}
	accessTTL := parseDuration(os.Getenv("JWT_ACCESS_TTL"), 15*time.Minute)
	refreshTTL := parseDuration(os.Getenv("JWT_REFRESH_TTL"), 7*24*time.Hour)

	access, err := security.SignAccess(u.ID, string(u.Role), u.EmployeeCode, accessTTL)
	if err != nil {
		return "", "", err
	}

	jti := newJTI()
	refresh, err := security.SignRefresh(u.ID, jti, refreshTTL)
	if err != nil {
		return "", "", err
	}

	if err := s.tokens.SaveRefresh(ctx, u.ID, jti, time.Now().Add(refreshTTL).Unix()); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *svc) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	c, err := security.ParseRefresh(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh")
	}

	revoked, err := s.tokens.IsRevoked(ctx, c.JTI)
	if err != nil || revoked {
		return "", "", errors.New("refresh revoked")
	}

	u, err := s.users.FindByID(ctx, c.UserID)
	if err != nil || u == nil || !u.IsActive {
		return "", "", errors.New("invalid user")
	}

	_ = s.tokens.Revoke(ctx, c.JTI) // rotate

	accessTTL := parseDuration(os.Getenv("JWT_ACCESS_TTL"), 15*time.Minute)
	refreshTTL := parseDuration(os.Getenv("JWT_REFRESH_TTL"), 7*24*time.Hour)

	access, err := security.SignAccess(u.ID, string(u.Role), u.EmployeeCode, accessTTL)
	if err != nil {
		return "", "", err
	}

	newJti := newJTI()
	newRefresh, err := security.SignRefresh(u.ID, newJti, refreshTTL)
	if err != nil {
		return "", "", err
	}
	if err := s.tokens.SaveRefresh(ctx, u.ID, newJti, time.Now().Add(refreshTTL).Unix()); err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

func (s *svc) Logout(ctx context.Context, refreshToken string) error {
	c, err := security.ParseRefresh(refreshToken)
	if err != nil {
		return errors.New("invalid refresh")
	}
	return s.tokens.Revoke(ctx, c.JTI)
}

func newJTI() string { b := make([]byte, 16); _, _ = rand.Read(b); return hex.EncodeToString(b) }
func parseDuration(s string, def time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return def
	}
	return d
}
