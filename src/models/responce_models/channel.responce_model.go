package models

import (
	models "caapp-server/src/models/db_models"
)

type GetChannelListResponceItem struct {
	Channel           models.Channel
	LastMessage       models.Message
	LastMessageSender GetUserInfoResponce
}

type GetChannelListResponce struct {
	Channels   []GetChannelListResponceItem `json:"channels"`
	IsLastPage bool                         `json:"is_last_page"`
}

type CheckFriendChannelResponce struct {
	Channel models.Channel `json:"channel"`
}

type GetFriendChannelInfoResponce struct {
	Channel models.Channel      `json:"channel"`
	User    GetUserInfoResponce `json:"user"`
}

type GetGroupChannelInfoResponce struct {
	Channel models.Channel        `json:"channel"`
	Users   []GetUserInfoResponce `json:"users"`
}
