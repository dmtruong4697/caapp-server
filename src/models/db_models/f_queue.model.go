package models

import "time"

type FQueue struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id"`
	CreateAt time.Time `json:"create_at"`
}
