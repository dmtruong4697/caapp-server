package rcdbmodels

import "time"

type RCMessage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id"`
	Content    string    `json:"content"`
	CreateAt   time.Time `json:"create_at"`
	LastUpdate time.Time `json:"last_update"`
	ChannelID  uint      `json:"channel_id"`
	IsEdited   bool      `json:"is_edited"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
}

// type:
// 0: message
// 1: user join channel
// 2: user leave channel
// 3: owner delete user from channel
// 4: channel name change
// 5: channel image change
