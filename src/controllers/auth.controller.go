package controllers

import (
	"caapp-server/src/database"
	"caapp-server/src/enums"
	models "caapp-server/src/models/db_models"
	utils "caapp-server/src/utils"
	"encoding/json"
	"net/http"
)

var JwtKey = []byte("20204697")

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser := models.User{}
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email already registered", http.StatusBadRequest)
		return
	}

	validateCode := utils.GenerateRandomCode(6)
	user.ValidateCode = validateCode
	user.AccountStatus = string(enums.USER_ACCOUNT_STATUS_NOT_ACTIVE)

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send email with validate code
	header := "Validate Your Email"
	body := "Validate code:" + validateCode
	utils.SendEmail(user.Email, header, body)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
