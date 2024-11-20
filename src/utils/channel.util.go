package utils

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
)

func GetFriendChannelInfo(current_user_id uint, channelID uint) responce_models.GetFriendChannelInfoResponce {

	var channel db_models.Channel
	if err := database.DB.Where("id = ?", channelID).First(&channel).Error; err != nil {
		return responce_models.GetFriendChannelInfoResponce{}
	}

	query := `
		SELECT *
		FROM channel_members
		WHERE channel_id = ?
		AND user_id != ?;
	`

	var channelMember db_models.ChannelMember
	result := database.DB.Raw(query, channelID, current_user_id).First(&channelMember)
	if result.Error != nil {
		return responce_models.GetFriendChannelInfoResponce{}
	}

	var res responce_models.GetFriendChannelInfoResponce
	res.Channel = channel
	res.User = GetUserInfo(current_user_id, channelMember.UserID)

	return res
}

func GetGroupChannelInfo(current_user_id uint, channelID uint) responce_models.GetGroupChannelInfoResponce {

	var channel db_models.Channel
	if err := database.DB.Where("id = ?", channelID).First(&channel).Error; err != nil {
		return responce_models.GetGroupChannelInfoResponce{}
	}

	query := `
		SELECT *
		FROM channel_members
		WHERE channel_id = ?
		AND user_id != ?;
	`

	var channelMembers []db_models.ChannelMember
	err := database.DB.Select(&channelMembers, query, channelID, current_user_id)
	if err != nil {
		return responce_models.GetGroupChannelInfoResponce{}
	}

	var res responce_models.GetGroupChannelInfoResponce
	res.Channel = channel
	for i := range channelMembers {
		res.Users[i] = GetUserInfo(current_user_id, channelMembers[i].UserID)
	}

	return res
}

func GetChannelInfo(channelID uint) db_models.Channel {
	var channel db_models.Channel
	if err := database.DB.Where("id = ?", channelID).First(&channel).Error; err != nil {
		return db_models.Channel{}
	}

	return channel
}
