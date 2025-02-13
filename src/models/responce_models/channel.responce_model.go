package models

import (
	models "caapp-server/src/models/db_models"
)

type GetChannelListResponceItem struct {
	Channel           models.Channel        `json:"channel"`
	Users             []GetUserInfoResponce `json:"users"`
	LastMessage       models.Message        `json:"last_message"`
	LastMessageSender GetUserInfoResponce   `json:"last_message_sender"`
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

type GetChannelChatHistoryItem struct {
	Sender  GetUserInfoResponce `json:"sender"`
	Message models.Message      `json:"message"`
	Media   []models.Media      `json:"media"`
}

type GetChannelChatHistoryResponce struct {
	Messages []GetChannelChatHistoryItem `json:"messages"`
}
