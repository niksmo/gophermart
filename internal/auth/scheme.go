package auth

import (
	"regexp"
	"strings"

	"github.com/niksmo/gophermart/internal/errs"
)

const (
	minLoginLen = 1
	maxLoginLen = 60

	minPasswordLen = 1
	maxPasswordLen = 60
)

var (
	validLogin = regexp.MustCompile(`^[\d\w-]+$`)

	validPassword = regexp.MustCompile(
		`^[\d\w\-!@#$%^&*()_+|\\\[\]{}'";:\/?>.<,=]+$`,
	)
)

type BaseRequestScheme struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (scheme BaseRequestScheme) Validate() (result any, ok bool) {
	return validateScheme(scheme.Login, scheme.Password)
}

type SignupRequestScheme struct {
	BaseRequestScheme
}

type SigninRequestScheme struct {
	BaseRequestScheme
}

type InvalidValidationData struct {
	Login    []string `json:"login,omitempty"`
	Password []string `json:"password,omitempty"`
}

func validateScheme(
	login, password string,
) (result InvalidValidationData, ok bool) {
	ok = true
	result.Login = validateLogin(login)
	result.Password = validatePassword(password)
	if result.Login != nil || result.Password != nil {
		ok = false
		return
	}
	return
}

func validateLogin(login string) []string {
	var resultSlice []string
	login = strings.TrimSpace(login)
	length := len(login)
	if length < minLoginLen || length > maxLoginLen {
		resultSlice = append(resultSlice, errs.ErrUserLoginLength.Error())
	}
	if !validLogin.MatchString(login) {
		resultSlice = append(
			resultSlice, errs.ErrUserLoginInvalid.Error(),
		)
	}
	if len(resultSlice) != 0 {
		return resultSlice
	}
	return nil
}

func validatePassword(password string) []string {
	var resultSlice []string
	password = strings.TrimSpace(password)
	length := len(password)
	if length < minPasswordLen || length > maxPasswordLen {
		resultSlice = append(resultSlice, errs.ErrUserPasswordLength.Error())
	}
	if !validPassword.MatchString(password) {
		resultSlice = append(resultSlice, errs.ErrUserPasswordInvalid.Error())
	}
	if len(resultSlice) != 0 {
		return resultSlice
	}
	return nil
}
