package rcdbmodels

import "time"

type MQueue struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id"`
	CreateAt time.Time `json:"create_at"`
}
