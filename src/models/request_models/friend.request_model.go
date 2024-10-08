package models

type GetFriendRequestRequest struct {
	ID uint `json:"id"`
}

type GetAllFriendRequestSentRequest struct {
}

type AcceptRequestRequest struct {
	ID uint `json:"id"`
}

type CreateFriendRequestRequest struct {
	UserID uint `json:"user_id"`
}

type AcceptFriendRequestRequest struct {
	ID uint `json:"id"`
}

type RefuseFriendRequestRequest struct {
	ID uint `json:"id"`
}

type DeleteFriendRequestRequest struct {
	ID uint `json:"id"`
}
