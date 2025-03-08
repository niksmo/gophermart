package auth

const tokenType = "Bearer"

type SignupReqPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

type SigninResPayload struct {
	TokenType  string `json:"tokenType"`
	TokenValue string `json:"tokenValue"`
}

func NewSigninResPayload(tokenValue string) SigninResPayload {
	return SigninResPayload{TokenType: tokenType, TokenValue: tokenValue}
}
