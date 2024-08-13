package models

type Media struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	SenderID  uint   `json:"sender_id"`
	MessageID uint   `json:"message_id"`
	ChannelID uint   `json:"channel_id"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	CreateAt  int    `json:"create_at"`
}
