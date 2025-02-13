package rcresponsemodels

import rcdbmodels "caapp-server/src/models/db_models/rc_db_models"

type GetRCChannelChatHistoryItem struct {
	Message rcdbmodels.RCMessage `json:"message"`
	Media   []rcdbmodels.RCMedia `json:"media"`
}

type GetRCChannelChatHistoryResponce struct {
	Messages []GetRCChannelChatHistoryItem `json:"messages"`
}
