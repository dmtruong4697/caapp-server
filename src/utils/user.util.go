package utils

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
)

func GetUserInfo(current_user_id uint, user_id uint) responce_models.GetUserInfoResponce {

	// get user from database
	var dbUser db_models.User
	if err := database.DB.Where("id = ?", user_id).First(&dbUser).Error; err != nil {

	}

	var friend db_models.Friend
	if err := database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", current_user_id, user_id, user_id, current_user_id).First(&friend).Error; err != nil {

	}

	var responce responce_models.GetUserInfoResponce
	responce.ID = int(dbUser.ID)
	responce.AccountStatus = dbUser.AccountStatus
	responce.AvatarImage = dbUser.AvatarImage
	responce.Country = dbUser.Country
	responce.CoverImage = dbUser.CoverImage
	responce.CreateAt = dbUser.CreateAt
	responce.DateOfBirth = dbUser.DateOfBirth
	responce.Email = dbUser.Email
	responce.FirstName = dbUser.FirstName
	responce.MiddleName = dbUser.MiddleName
	responce.LastName = dbUser.LastName
	responce.Gender = dbUser.Gender
	responce.HashtagName = dbUser.HashtagName
	responce.JobName = dbUser.JobName
	responce.Language = dbUser.Language
	responce.VerificationStatus = dbUser.VerificationStatus
	responce.TimeZone = dbUser.TimeZone
	responce.ProfileDescription = dbUser.ProfileDescription
	responce.PhoneNumber = dbUser.PhoneNumber
	responce.LastActive = dbUser.LastActive
	responce.LastUpdate = dbUser.LastUpdate
	responce.Friend = friend

	return responce
}
