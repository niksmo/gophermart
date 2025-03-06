package auth

import (
	"context"
	"errors"

	"github.com/niksmo/gophermart/internal/repository"
)

type AuthService struct {
	repository repository.UsersRepository
}

func NewService(repository repository.UsersRepository) AuthService {
	return AuthService{repository: repository}
}

func (s AuthService) RegisterUser(ctx context.Context, login, password string) error {
	err := s.repository.Create(ctx, login, password)
	if errors.Is(err, repository.ErrExists) {
		return ErrLoginExists
	}
	return err
}
