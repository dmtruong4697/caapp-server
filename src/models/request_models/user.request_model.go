package models

type GetUserInfoRequestBody struct {
	ID uint `json:"id"`
}

type GetUserFriendRequest struct {
	UserID uint `json:"id"`
}
