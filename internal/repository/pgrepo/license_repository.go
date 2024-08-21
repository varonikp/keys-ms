package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/varonikp/keys-ms/internal/domain"
	"github.com/varonikp/keys-ms/internal/repository/models"
)

type LicenseRepository struct {
	db *sqlx.DB
}

func NewLicenseRepository(db *sqlx.DB) *LicenseRepository {
	return &LicenseRepository{
		db: db,
	}
}

func (r LicenseRepository) CreateLicense(ctx context.Context, license domain.License) (domain.License, error) {
	op := "pgrepo.LicenseRepository.CreateLicense"

	modelsLicense := domainToLicense(license)
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO licenses (software_id, user_id, created_at, expire_at) VALUES (:software_id, :user_id, :created_at, :expire_at) ", modelsLicense)
	if err != nil {
		return domain.License{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.License{}, fmt.Errorf("%s(id): %w", op, err)
	}

	modelsLicense.ID = int(id)

	return licenseToDomain(modelsLicense), nil
}

func (r LicenseRepository) GetLicense(ctx context.Context, id int) (domain.License, error) {
	op := "pgrepo.LicenseRepository.GetLicense"

	var license models.License
	err := r.db.GetContext(ctx, &license, "SELECT * FROM licenses WHERE id=?", id)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.License{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.License{}, fmt.Errorf("%s: %w", op, err)
	}

	return licenseToDomain(license), nil
}

func (r LicenseRepository) GetLicensesByUserID(ctx context.Context, userId int) ([]domain.License, error) {
	op := "pgrepo.LicenseRepository.GetLicensesByUserID"

	var licenses []models.License
	err := r.db.SelectContext(ctx, &licenses, "SELECT * FROM licenses WHERE user_id=?", userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// convert all to domain
	out := make([]domain.License, 0, len(licenses))
	for _, license := range licenses {
		domainLicense := licenseToDomain(license)

		out = append(out, domainLicense)
	}

	return out, nil
}

func (r LicenseRepository) UpdateLicense(ctx context.Context, license domain.License) (domain.License, error) {
	op := "pgrepo.LicenseRepository.UpdateLicense"

	modelsLicense := domainToLicense(license)
	_, err := r.db.NamedExecContext(ctx, "UPDATE licenses SET software_id=:software_id, user_id=:user_id, expire_at=:expire_at WHERE id=:id", modelsLicense)
	if err != nil {
		return domain.License{}, fmt.Errorf("%s: %w", op, err)
	}

	return licenseToDomain(modelsLicense), nil
}

func (r LicenseRepository) DeleteLicense(ctx context.Context, id int) error {
	op := "pgrepo.LicenseRepository.DeleteLicense"

	_, err := r.db.ExecContext(ctx, "DELETE FROM licenses WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
