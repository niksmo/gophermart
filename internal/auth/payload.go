package auth

type SignupReqPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignupResPayload struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type SigninReqPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SigninResPayload struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
