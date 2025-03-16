package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int32 `json:"userID"`
}

func newClaims(userID int32, lifetime time.Duration) Claims {
	expiresAt := jwt.NewNumericDate(time.Now().Add(lifetime))
	registeredClaims := jwt.RegisteredClaims{ExpiresAt: expiresAt}
	return Claims{registeredClaims, userID}
}

func Create(
	ID int32, key []byte, lifetime time.Duration,
) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, newClaims(ID, lifetime),
	)
	tokenString, err = token.SignedString(key)
	return
}

func Parse(tokenString string, key []byte) (int32, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc(key))
	return claims.UserID, err
}

func keyFunc(key []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		if !isValidMethod(t.Method) {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	}
}

func isValidMethod(method jwt.SigningMethod) bool {
	return method.Alg() == jwt.SigningMethodHS256.Alg()
}
