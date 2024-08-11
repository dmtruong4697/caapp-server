package models

import (
	"time"
)

type ChannelMember struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	UserID        int       `json:"user_id"`
	ChannelID     int       `json:"channel_id"`
	JoinAt        time.Time `json:"join_at"`
	Role          string    `json:"role"`
	InChannelName string    `json:"inchannel_name"`
}
