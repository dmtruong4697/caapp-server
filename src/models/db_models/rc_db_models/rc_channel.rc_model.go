package rcdbmodels

import (
	"time"
)

type RCChannel struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CreateAt      time.Time `json:"create_at"`
	LastMessageID uint      `json:"last_message_id"`
	Type          string    `json:"type"`
}
