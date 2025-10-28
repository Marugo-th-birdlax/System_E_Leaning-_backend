package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Emp    string `json:"emp_code"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"uid"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func SignAccess(uid, role, emp string, ttl time.Duration) (string, error) {
	secret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	now := time.Now()
	claims := AccessClaims{
		UserID: uid, Role: role, Emp: emp,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func SignRefresh(uid, jti string, ttl time.Duration) (string, error) {
	secret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	now := time.Now()
	claims := RefreshClaims{
		UserID: uid, JTI: jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func ParseAccess(tokenStr string) (*AccessClaims, error) {
	secret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	tok, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := tok.Claims.(*AccessClaims); ok && tok.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func ParseRefresh(tokenStr string) (*RefreshClaims, error) {
	secret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	tok, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := tok.Claims.(*RefreshClaims); ok && tok.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
