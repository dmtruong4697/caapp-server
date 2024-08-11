package models

import "time"

type User struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	Email              string    `json:"email"`
	PhoneNumber        string    `json:"phone_number"`
	Password           string    `json:"password"`
	FirstName          string    `json:"first_name"`
	MiddleName         string    `json:"middle_name"`
	LastName           string    `json:"last_name"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	HashtagName        string    `json:"hashtag_name"`
	Gender             string    `json:"gender"`
	Language           string    `json:"language"`
	Country            string    `json:"country"`
	ProfileDescription string    `json:"profile_description"`
	AvatarImage        string    `json:"avatar_image"`
	CoverImage         string    `json:"cover_image"`
	ValidateCode       string    `json:"validate_code"`
	AccountStatus      string    `json:"account_status"`
	VerificationStatus string    `json:"verification_status"`
	CreateAt           time.Time `json:"create_at"`
	LastUpdate         time.Time `json:"last_update"`
	LastActive         time.Time `json:"last_active"`
	DeviceToken        string    `json:"device_token"`
	JobName            string    `json:"job_name"`
	TimeZone           string    `json:"time_zone"`
}
