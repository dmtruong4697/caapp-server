package models

import (
	db_models "caapp-server/src/models/db_models"
)

type MessageDetail struct {
	Sender  GetUserInfoResponce `json:"sender"`
	Message db_models.Message   `json:"message"`
	Media   []db_models.Media   `json:"media"`
}

type ChannelMessages struct {
	Messages []MessageDetail
}
