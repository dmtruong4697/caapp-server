package rcdbmodels

import "time"

type GeneralQueueUser struct {
	UserID       uint      `json:"user_id" gorm:"primaryKey"`
	JoinAt       time.Time `json:"join_at"`
	Gender       string    `json:"gender"`
	TargetGender string    `json:"target_gender"`
}
