package auth

import (
	"context"

	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/errs"
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

	return s.createToken(userID)
}

func (s AuthService) AuthorizeUser(
	ctx context.Context, login, password string,
) (string, error) {
	userID, pwdHash, err := s.repository.ReadByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password))
	if err != nil {
		return "", errs.ErrCredentials
	}

	return s.createToken(userID)
}

func (s AuthService) createToken(userID int64) (string, error) {
	tokenString, err := jwt.Create(
		userID, s.authConfig.Key(), s.authConfig.JWTLifetime(),
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
