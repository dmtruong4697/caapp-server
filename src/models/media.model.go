package models

type Media struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	SenderID  int    `json:"sender_id"`
	MessageID int    `json:"message_id"`
	ChannelID int    `json:"channel_id"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	CreateAt  int    `json:"create_at"`
}
