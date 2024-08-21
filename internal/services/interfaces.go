package services

import (
	"context"

	"github.com/varonikp/keys-ms/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, login string) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type SoftwareRepository interface {
	CreateSoftware(ctx context.Context, software domain.Software) (domain.Software, error)
	GetSoftwares(ctx context.Context) ([]domain.Software, error)
	GetSoftware(ctx context.Context, id int) (domain.Software, error)
	UpdateSoftware(ctx context.Context, software domain.Software) (domain.Software, error)
	DeleteSoftware(ctx context.Context, id int) error
}

type LicenseRepository interface {
	CreateLicense(ctx context.Context, license domain.License) (domain.License, error)
	GetLicense(ctx context.Context, id int) (domain.License, error)
	GetLicensesByUserID(ctx context.Context, id int) ([]domain.License, error)
	UpdateLicense(ctx context.Context, license domain.License) (domain.License, error)
	DeleteLicense(ctx context.Context, id int) error
}
