package models

import "time"

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id"`
	Content    string    `json:"content"`
	CreateAt   time.Time `json:"create_at"`
	LastUpdate time.Time `json:"last_update"`
	ChannelID  uint      `json:"channel_id"`
	IsEdited   bool      `json:"is_edited"`
	Status     string    `json:"status"`
}
