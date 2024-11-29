package rcdbmodels

import (
	"time"
)

type RCChannelMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ChannelID uint      `json:"channel_id"`
	JoinAt    time.Time `json:"join_at"`
}
