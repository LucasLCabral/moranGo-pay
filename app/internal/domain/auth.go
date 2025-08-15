package domain

type LoginResult struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

type ConfirmUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
