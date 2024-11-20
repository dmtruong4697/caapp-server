package models

import (
	"time"
)

type EmailValidateCode struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email"`
	ValidateCode string    `json:"validate_code"`
	CreateAt     time.Time `json:"create_at"`
	ExpireAt     time.Time `json:"expire_at"`
}
