package models

import (
	models "caapp-server/src/models/db_models"
)

type GetFriendRequestResponce struct {
	FriendRequest models.FriendRequest
}

type GetListFriendRequestReceivedResponceItem struct {
	User          GetUserInfoResponce `json:"user"`
	FriendRequest models.FriendRequest
}

type GetListFriendRequestReceivedResponce struct {
	Requests []GetListFriendRequestReceivedResponceItem `json:"requests"`
}

type GetListFriendRequestSentResponceItem struct {
	User          GetUserInfoResponce `json:"user"`
	FriendRequest models.FriendRequest
}

type GetListFriendRequestSentResponce struct {
	Requests []GetListFriendRequestReceivedResponceItem `json:"requests"`
}

type GetSuggestUserResponce struct {
	Users []GetUserInfoResponce `json:"users"`
}
