package models

import (
	models "caapp-server/src/models/db_models"
)

type GetFriendRequestResponce struct {
	FriendRequest models.FriendRequest
}

type GetListFriendRequestReceivedResponceItem struct {
	User          GetUserInfoResponce  `json:"user"`
	FriendRequest models.FriendRequest `json:"friend_request"`
}

type GetListFriendRequestReceivedResponce struct {
	Requests []GetListFriendRequestReceivedResponceItem `json:"requests"`
}

type GetListFriendRequestSentResponceItem struct {
	User          GetUserInfoResponce  `json:"user"`
	FriendRequest models.FriendRequest `json:"friend_request"`
}

type GetListFriendRequestSentResponce struct {
	Requests []GetListFriendRequestReceivedResponceItem `json:"requests"`
}

type GetSuggestUserResponce struct {
	Users []GetUserInfoResponce `json:"users"`
}

type GetAllMyFriendResponce struct {
	Friends []GetUserInfoResponce `json:"friends"`
}
