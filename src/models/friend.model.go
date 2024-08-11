package models

import (
	"time"
)

type Friend struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	FirstUserID  int       `json:"first_user_id"`
	SecondUserID int       `json:"second_user_id"`
	CreateAt     time.Time `json:"create_at"`
}
