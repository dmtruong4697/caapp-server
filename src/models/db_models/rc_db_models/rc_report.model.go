package rcdbmodels

import "time"

type RCReport struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SenderID  uint      `json:"sender_id"`
	ChannelID uint      `json:"channel_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	CreateAt  time.Time `json:"create_at"`
}
