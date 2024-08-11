package models

import (
	"time"
)

type FriendRequest struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	CreateAt   time.Time `json:"create_at"`
}
