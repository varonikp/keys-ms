package services

import (
	"context"

	"github.com/varonikp/keys-ms/internal/domain"
)

type SoftwareService struct {
	repo SoftwareRepository
}

func NewSoftwareService(repo SoftwareRepository) SoftwareService {
	return SoftwareService{
		repo: repo,
	}
}

func (s SoftwareService) GetSoftware(ctx context.Context, id int) (domain.Software, error) {
	return s.repo.GetSoftware(ctx, id)
}

func (s SoftwareService) GetSoftwares(ctx context.Context) ([]domain.Software, error) {
	return s.repo.GetSoftwares(ctx)
}

func (s SoftwareService) CreateSoftware(ctx context.Context, software domain.Software) (domain.Software, error) {
	return s.repo.CreateSoftware(ctx, software)
}

func (s SoftwareService) UpdateSoftware(ctx context.Context, software domain.Software) (domain.Software, error) {
	return s.repo.UpdateSoftware(ctx, software)
}

func (s SoftwareService) DeleteSoftware(ctx context.Context, id int) error {
	return s.repo.DeleteSoftware(ctx, id)
}
