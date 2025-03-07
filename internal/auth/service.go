package auth

import (
	"context"

	"github.com/niksmo/gophermart/internal/config"
	"github.com/niksmo/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authConfig config.AuthConfig
	repository repository.UsersRepository
}

func NewService(authConfig config.AuthConfig, repository repository.UsersRepository) AuthService {
	return AuthService{authConfig: authConfig, repository: repository}
}

func (s AuthService) RegisterUser(ctx context.Context, login, password string) error {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), s.authConfig.Cost())
	if err != nil {
		return err
	}
	return s.repository.Create(ctx, login, string(pwdHash))
}
