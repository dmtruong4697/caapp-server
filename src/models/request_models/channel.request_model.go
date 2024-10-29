package models

type GetFriendChannelInfoRequest struct {
	ChannelID uint `json:"channel_id"`
}

type GetGroupChannelInfoRequest struct {
	ChannelID uint `json:"channel_id"`
}

type GetChannelListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CheckFriendChannelRequest struct {
	UserID uint `json:"user_id"`
}
