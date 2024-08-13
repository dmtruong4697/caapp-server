package models

type FriendRequestRequest struct {
	ReceiverID int `json:"receiver_id"`
}

type AcceptRequestRequest struct {
	ID int `json:"id"`
}

type GetRelationshipRequest struct {
	UserID uint `json:"user_id"`
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
