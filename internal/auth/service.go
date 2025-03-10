package auth

import (
	"context"

	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/bonuses"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/internal/users"
	"github.com/niksmo/gophermart/pkg/jwt"
	"github.com/niksmo/gophermart/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authConfig        config.AuthConfig
	usersRepository   users.UsersRepository
	bonusesRepository bonuses.BonusesRepository
}

func NewService(
	authConfig config.AuthConfig,
	usersRepository users.UsersRepository,
	bonusesRepository bonuses.BonusesRepository,
) AuthService {
	return AuthService{
		authConfig:        authConfig,
		usersRepository:   usersRepository,
		bonusesRepository: bonusesRepository,
	}
}

func (s AuthService) RegisterUser(
	ctx context.Context, login, password string,
) (string, error) {
	pwdHash, err := bcrypt.GenerateFromPassword(
		[]byte(password), s.authConfig.Cost(),
	)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return "", err
	}
	userID, err := s.usersRepository.Create(ctx, login, string(pwdHash))
	if err != nil {
		return "", err
	}

	err = s.bonusesRepository.CreateAccount(ctx, int32(userID))
	if err != nil {
		return "", err
	}

	return s.createToken(userID)
}

func (s AuthService) AuthorizeUser(
	ctx context.Context, login, password string,
) (string, error) {
	userID, pwdHash, err := s.usersRepository.ReadByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password))
	if err != nil {
		return "", errs.ErrUserCredentials
	}

	return s.createToken(userID)
}

func (s AuthService) createToken(userID int32) (string, error) {
	tokenString, err := jwt.Create(
		userID, s.authConfig.Key(), s.authConfig.JWTLifetime(),
	)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return "", err
	}

	return tokenString, nil
}
