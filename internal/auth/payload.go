package auth

import (
	"regexp"
	"strings"

	"github.com/niksmo/gophermart/internal/errs"
)

const (
	tokenType = "Bearer"

	minLoginLen = 1
	maxLoginLen = 60

	minPasswordLen = 1
	maxPasswordLen = 60
)

var (
	validLogin    = regexp.MustCompile(`^[\d\w-]+$`)
	validPassword = regexp.MustCompile(`^[\d\w\-!@#$%^&*()_+|\\\[\]{}'";:\/?>.<,=]+$`)
)

type SignupReqPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (signupReq SignupReqPayload) Validate() (result InvalidValidationData, ok bool) {
	return validatePayload(signupReq.Login, signupReq.Password)
}

type SignupResPayload struct {
	TokenType  string `json:"tokenType"`
	TokenValue string `json:"tokenValue"`
}

func NewSignupResPayload(tokenValue string) SignupResPayload {
	return SignupResPayload{TokenType: tokenType, TokenValue: tokenValue}
}

type SigninReqPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (signinReq SigninReqPayload) Validate() (result InvalidValidationData, ok bool) {
	return validatePayload(signinReq.Login, signinReq.Password)
}

type SigninResPayload struct {
	TokenType  string `json:"tokenType"`
	TokenValue string `json:"tokenValue"`
}

func NewSigninResPayload(tokenValue string) SigninResPayload {
	return SigninResPayload{TokenType: tokenType, TokenValue: tokenValue}
}

type InvalidValidationData struct {
	Login    []string `json:"login,omitempty"`
	Password []string `json:"password,omitempty"`
}

func validatePayload(
	login, password string,
) (result InvalidValidationData, ok bool) {
	result.Login = validateLogin(login)
	result.Password = validatePassword(password)
	if result.Login != nil || result.Password != nil {
		ok = false
		return
	}
	ok = true
	return
}

func validateLogin(login string) []string {
	var sErr []string
	login = strings.TrimSpace(login)
	length := len(login)
	if length < minLoginLen || length > maxLoginLen {
		sErr = append(sErr, errs.ErrLoginLength.Error())
	}
	if !validLogin.MatchString(login) {
		sErr = append(
			sErr, errs.ErrLoginInvalid.Error(),
		)
	}
	if len(sErr) != 0 {
		return sErr
	}
	return nil
}

func validatePassword(password string) []string {
	var sErr []string
	password = strings.TrimSpace(password)
	length := len(password)
	if length < minPasswordLen || length > maxPasswordLen {
		sErr = append(sErr, errs.ErrPasswordLength.Error())
	}
	if !validPassword.MatchString(password) {
		sErr = append(sErr, errs.ErrPasswordInvalid.Error())
	}
	if len(sErr) != 0 {
		return sErr
	}
	return nil
}
