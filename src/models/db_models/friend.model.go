package models

import (
	"time"
)

type Friend struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FirstUserID  uint      `json:"first_user_id"`
	SecondUserID uint      `json:"second_user_id"`
	CreateAt     time.Time `json:"create_at"`
}
