package controllers

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	"encoding/json"
	"net/http"
	"time"
)

func GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("id").(uint)

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", user_id).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbUser); err != nil {
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}
}

func UpdateProfileInfo(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	var dbUser db_models.User
	if err := database.DB.Where("email = ?", email).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var updatedProfile request_models.UpdatedProfileRequest
	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		http.Error(w, "Failed to decode user info", http.StatusBadRequest)
		return
	}

	// update user profile
	dbUser.FirstName = updatedProfile.FirstName
	dbUser.MiddleName = updatedProfile.MiddleName
	dbUser.LastName = updatedProfile.LastName
	dbUser.PhoneNumber = updatedProfile.PhoneNumber
	dbUser.AvatarImage = updatedProfile.AvatarImage
	dbUser.CoverImage = updatedProfile.CoverImage
	dbUser.HashtagName = updatedProfile.HashtagName
	dbUser.Gender = updatedProfile.Gender
	dbUser.DateOfBirth = updatedProfile.DateOfBirth
	dbUser.Language = updatedProfile.Language
	dbUser.Country = updatedProfile.Country
	dbUser.ProfileDescription = updatedProfile.ProfileDescription
	dbUser.JobName = updatedProfile.JobName
	dbUser.TimeZone = updatedProfile.TimeZone
	dbUser.LastUpdate = time.Now()

	// save update
	if err := database.DB.Save(&dbUser).Error; err != nil {
		http.Error(w, "Failed to update user information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbUser); err != nil {
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("id").(uint)

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", user_id).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var updatedPassword request_models.UpdatedPaswordRequest
	err := json.NewDecoder(r.Body).Decode(&updatedPassword)
	if err != nil {
		http.Error(w, "Failed to decode user info", http.StatusBadRequest)
		return
	}

	// update account password
	dbUser.Password = updatedPassword.Password
	dbUser.LastUpdate = time.Now()

	// save update
	if err := database.DB.Save(&dbUser).Error; err != nil {
		http.Error(w, "Failed to update user information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbUser); err != nil {
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}
}
