package wsmodels

import (
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
)

type WSChannelListItemForListening struct {
	UserID    uint              `json:"user_id"`
	ChannelID uint              `json:"channel_id"`
	Message   db_models.Message `json:"message"`
}

type WSChannelListItemForBroadcast struct {
	Channel           db_models.Channel                     `json:"channel"`
	Users             []responce_models.GetUserInfoResponce `json:"users"`
	LastMessage       db_models.Message                     `json:"last_message"`
	LastMessageSender responce_models.GetUserInfoResponce   `json:"last_message_sender"`
}

type WSChannelListItemForHandleMessage struct {
	Data                     WSChannelListItemForBroadcast `json:"data"`
	UserIDForMakingChannelID string                        `json:"user_id_for_making_channel_id"`
}
