package models

import "time"

type Message struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	SenderID   int       `json:"sender_id"`
	Content    string    `json:"content"`
	CreateAt   time.Time `json:"create_at"`
	LastUpdate time.Time `json:"last_update"`
	ChannelID  int       `json:"channel_id"`
	IsEdited   bool      `json:"is_edited"`
	Status     string    `json:"status"`
}
