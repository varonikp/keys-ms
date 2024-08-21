package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/varonikp/keys-ms/internal/domain"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s UserService) IsUserExists(ctx context.Context, login string) (bool, error) {
	user, err := s.repo.GetUser(ctx, login)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if err == nil && user.ID() != 0 {
		return true, nil
	}

	return false, nil
}

func (s UserService) GetUser(ctx context.Context, login string) (domain.User, error) {
	return s.repo.GetUser(ctx, login)
}

func (s UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s UserService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s UserService) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.UpdateUser(ctx, user)
}

func (s UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}
