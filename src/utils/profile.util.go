package utils

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
)

func GetProfileInfo(current_user_id uint) responce_models.GetUserInfoResponce {

	// get user from database
	var dbUser db_models.User
	if err := database.DB.Where("id = ?", current_user_id).First(&dbUser).Error; err != nil {
		return responce_models.GetUserInfoResponce{}
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

	return responce
}
