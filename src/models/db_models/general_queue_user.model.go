package models

import "time"

type GeneralQueueUser struct {
	UserID       uint      `json:"user_id"`
	JoinAt       time.Time `json:"join_at"`
	Gender       string    `json:"gender"`
	TargetGender string    `json:"target_gender"`
}
