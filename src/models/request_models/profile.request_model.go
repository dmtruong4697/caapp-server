package models

import "time"

type UpdatedProfileRequest struct {
	FirstName          string    `json:"first_name"`
	MiddleName         string    `json:"middle_name"`
	LastName           string    `json:"last_name"`
	PhoneNumber        string    `json:"phone_number"`
	AvatarImage        string    `json:"avatar_image"`
	CoverImage         string    `json:"cover_image"`
	HashtagName        string    `json:"hashtag_name"`
	Gender             string    `json:"gender"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	Language           string    `json:"language"`
	Country            string    `json:"country"`
	ProfileDescription string    `json:"profile_description"`
	JobName            string    `json:"job_name"`
	TimeZone           string    `json:"time_zone"`
}

type UpdatedPaswordRequest struct {
	Password string `json:"password"`
}
