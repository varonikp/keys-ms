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

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	op := "pgrepo.UserRepository.CreateLicense"

	modelsUser := domainToUser(user)
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO users (login, password, has_admin_role) VALUES (:login, :password, :has_admin_role)", modelsUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, fmt.Errorf("%s(id): %w", op, err)
	}

	modelsUser.ID = int(id)

	return userToDomain(modelsUser), nil
}

func (r UserRepository) GetUser(ctx context.Context, login string) (domain.User, error) {
	op := "pgrepo.UserRepository.GetUser"

	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login=?", login)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userToDomain(user), nil
}

func (r UserRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	op := "pgrepo.UserRepository.GetUsers"

	var modelUsers []models.User
	err := r.db.SelectContext(ctx, &modelUsers, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	out := make([]domain.User, 0, len(modelUsers))
	for _, user := range modelUsers {
		domainUser := userToDomain(user)

		out = append(out, domainUser)
	}

	return out, nil
}

func (r UserRepository) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	op := "pgrepo.UserRepository.GetUserByID"

	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id=?", id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userToDomain(user), nil
}

func (r UserRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	op := "pgrepo.UserRepository.UpdateUser"

	modelsUser := domainToUser(user)
	_, err := r.db.NamedExecContext(ctx, "UPDATE users SET login=:login, password=:password, has_admin_role=:has_admin_role WHERE id=:id", modelsUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userToDomain(modelsUser), nil
}

func (r UserRepository) DeleteUser(ctx context.Context, id int) error {
	op := "pgrepo.UserRepository.UpdateUser"

	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
