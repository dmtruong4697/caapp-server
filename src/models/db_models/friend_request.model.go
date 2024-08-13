package models

import (
	"time"
)

type FriendRequest struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	CreateAt   time.Time `json:"create_at"`
}
