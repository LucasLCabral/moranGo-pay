package domain

type ConfirmUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
