package services

import (
	"context"

	"github.com/varonikp/keys-ms/internal/domain"
)

type LicenseService struct {
	repo LicenseRepository
}

func NewLicenseService(repo LicenseRepository) LicenseService {
	return LicenseService{
		repo: repo,
	}
}

func (s LicenseService) CreateLicense(ctx context.Context, license domain.License) (domain.License, error) {
	return s.repo.CreateLicense(ctx, license)
}

func (s LicenseService) GetLicense(ctx context.Context, id int) (domain.License, error) {
	return s.repo.GetLicense(ctx, id)
}

func (s LicenseService) GetLicensesByUserID(ctx context.Context, userId int) ([]domain.License, error) {
	return s.repo.GetLicensesByUserID(ctx, userId)
}

func (s LicenseService) UpdateLicense(ctx context.Context, license domain.License) (domain.License, error) {
	return s.repo.UpdateLicense(ctx, license)
}

func (s LicenseService) DeleteLicense(ctx context.Context, id int) error {
	return s.repo.DeleteLicense(ctx, id)
}
