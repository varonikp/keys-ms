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

type SoftwareRepo struct {
	db *sqlx.DB
}

func NewSoftwareRepo(db *sqlx.DB) *SoftwareRepo {
	return &SoftwareRepo{
		db: db,
	}
}

func (r SoftwareRepo) CreateSoftware(ctx context.Context, software domain.Software) (domain.Software, error) {
	op := "pgrepo.SoftwareRepo.CreateSoftware"

	modelsSoftware := domainToSoftware(software)
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO softwares (name) VALUES (:name)", modelsSoftware)
	if err != nil {
		return domain.Software{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Software{}, fmt.Errorf("%s(id): %w", op, err)
	}

	modelsSoftware.ID = int(id)

	return softwareToDomain(modelsSoftware), nil
}

func (r SoftwareRepo) GetSoftware(ctx context.Context, id int) (domain.Software, error) {
	op := "pgrepo.SoftwareRepo.GetSoftware"

	var software models.Software
	err := r.db.GetContext(ctx, &software, "SELECT * FROM softwares WHERE id=?", id)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.Software{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.Software{}, fmt.Errorf("%s: %w", op, err)
	}

	return softwareToDomain(software), nil
}

func (r SoftwareRepo) UpdateSoftware(ctx context.Context, software domain.Software) (domain.Software, error) {
	op := "pgrepo.SoftwareRepo.UpdateSoftware"

	modelsSoftware := domainToSoftware(software)
	_, err := r.db.NamedExecContext(ctx, "UPDATE softwares SET name=:name WHERE id=:id", modelsSoftware)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.Software{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.Software{}, fmt.Errorf("%s: %w", op, err)
	}

	modelsSoftware.Name = software.Name()

	return softwareToDomain(modelsSoftware), nil
}

func (r SoftwareRepo) DeleteSoftware(ctx context.Context, id int) error {
	op := "pgrepo.SoftwareRepo.DeleteSoftware"

	_, err := r.db.ExecContext(ctx, "DELETE FROM softwares WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r SoftwareRepo) GetSoftwares(ctx context.Context) ([]domain.Software, error) {
	op := "pgrepo.SoftwareRepo.GetSoftwares"

	var modelSoftwares []models.Software
	err := r.db.SelectContext(ctx, &modelSoftwares, "SELECT * FROM softwares")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	out := make([]domain.Software, 0, len(modelSoftwares))
	for _, software := range modelSoftwares {
		domainSoftware := softwareToDomain(software)

		out = append(out, domainSoftware)
	}

	return out, nil
}
