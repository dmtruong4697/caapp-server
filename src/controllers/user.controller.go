package controllers

import (
	"encoding/json"
	"net/http"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	current_user_id := r.Context().Value("id").(int)

	// get user id from request body
	var user_id uint
	err := json.NewDecoder(r.Body).Decode(&user_id)
	if err != nil {
		http.Error(w, "Failed to decode user info", http.StatusBadRequest)
		return
	}

	// get user from database
	var dbUser db_models.User
	if err := database.DB.Where("id = ?", user_id).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responce); err != nil {
		http.Error(w, "Failed to encode responce info", http.StatusInternalServerError)
	}
}
