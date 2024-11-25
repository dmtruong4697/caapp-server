package models

type LoginResponse struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type ValidateEmailResponse struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}
