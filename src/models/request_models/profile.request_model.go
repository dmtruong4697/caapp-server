package models

import "time"

type UpdateProfileRequest struct {
	PhoneNumber        string    `json:"phone_number"`
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
	JobName            string    `json:"job_name"`
	TimeZone           string    `json:"time_zone"`
}

type UpdatedPaswordRequest struct {
	Password string `json:"password"`
}

type GetProfileInfoRequest struct {
	ID uint `json:"id"`
}

type FirstUpdateProfileInfoRequest struct {
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	HashtagName string    `json:"hashtag_name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Country     string    `json:"country"`
	Language    string    `json:"language"`
}

type CheckDuplicateHashtagNameRequest struct {
	HashtagName string `json:"hashtag_name"`
}
