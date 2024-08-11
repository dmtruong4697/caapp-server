package models

import (
	"time"
)

type Channel struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	CreatorID         int       `json:"creator_id"`
	CreateAt          time.Time `json:"create_at"`
	InviteCode        string    `json:"invite_code"`
	LastMessageID     int       `json:"last_message_id"`
	ChannelImage      string    `json:"channel_image"`
	Type              string    `json:"type"`
	IsAllowInviteCode bool      `json:"is_allow_invite_code"`
}
