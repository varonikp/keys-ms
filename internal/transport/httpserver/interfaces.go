package httpserver

import (
	"context"

	"github.com/varonikp/keys-ms/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, login string) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	IsUserExists(ctx context.Context, login string) (bool, error)
	GetUserByID(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type TokenService interface {
	GenerateToken(user domain.User) (string, error)
	GetUser(token string) (domain.User, error)
}

type SoftwareService interface {
	CreateSoftware(ctx context.Context, software domain.Software) (domain.Software, error)
	GetSoftwares(ctx context.Context) ([]domain.Software, error)
	GetSoftware(ctx context.Context, id int) (domain.Software, error)
	UpdateSoftware(ctx context.Context, software domain.Software) (domain.Software, error)
	DeleteSoftware(ctx context.Context, id int) error
}

type LicenseService interface {
	CreateLicense(ctx context.Context, license domain.License) (domain.License, error)
	GetLicense(ctx context.Context, id int) (domain.License, error)
	GetLicensesByUserID(ctx context.Context, id int) ([]domain.License, error)
	UpdateLicense(ctx context.Context, license domain.License) (domain.License, error)
	DeleteLicense(ctx context.Context, id int) error
}
