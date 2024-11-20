package models

import "github.com/dgrijalva/jwt-go"

type LoginClaims struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type ValidateEmailRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ValidateCode string `json:"validate_code"`
}

type LoginRequestBody struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DeviceToken string `json:"device_token"`
}

type LogoutRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResendValidateCodeRequest struct {
	Email string `json:"email"`
}
