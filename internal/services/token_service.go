package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/varonikp/keys-ms/internal/domain"
)

type TokenService struct {
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewTokenService(jwtSecret []byte, tokenTTL time.Duration) TokenService {
	return TokenService{
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (s TokenService) GenerateToken(user domain.User) (string, error) {
	payload := UserClaims{
		UserId:       user.ID(),
		Login:        user.Login(),
		HasAdminRole: user.HasAdminRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return token, nil
}

func (s TokenService) GetUser(token string) (domain.User, error) {
	var claims UserClaims

	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !t.Valid {
		return domain.User{}, errors.New("token is not valid")
	}

	user := domain.NewUser(domain.NewUserData{
		ID:           claims.UserId,
		Login:        claims.Login,
		HasAdminRole: claims.HasAdminRole,
	})

	return user, nil
}

type UserClaims struct {
	UserId       int    `json:"user_id"`
	Login        string `json:"login"`
	HasAdminRole bool   `json:"has_admin_role"`
	jwt.RegisteredClaims
}
