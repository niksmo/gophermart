package auth

import (
	"context"

	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/repository"
	"github.com/niksmo/gophermart/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authConfig config.AuthConfig
	repository repository.UsersRepository
}

func NewService(authConfig config.AuthConfig, repository repository.UsersRepository) AuthService {
	return AuthService{authConfig: authConfig, repository: repository}
}

func (s AuthService) RegisterUser(
	ctx context.Context, login, password string,
) (string, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), s.authConfig.Cost())
	if err != nil {
		return "", err
	}
	userID, err := s.repository.Create(ctx, login, string(pwdHash))
	if err != nil {
		return "", err
	}

	tokenString, err := jwt.Create(
		userID, s.authConfig.Key(), s.authConfig.JWTLifetime(),
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
