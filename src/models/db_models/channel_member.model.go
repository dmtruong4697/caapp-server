package models

import (
	"time"
)

type ChannelMember struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id"`
	ChannelID     uint      `json:"channel_id"`
	JoinAt        time.Time `json:"join_at"`
	Role          string    `json:"role"`
	InChannelName string    `json:"inchannel_name"`
}
