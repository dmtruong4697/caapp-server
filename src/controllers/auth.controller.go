package controllers

import (
	"caapp-server/src/database"
	"caapp-server/src/enums"
	models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	utils "caapp-server/src/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func ValidateEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var validateEmailRequestBody request_models.ValidateEmailRequestBody
	if err := json.NewDecoder(r.Body).Decode(&validateEmailRequestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", validateEmailRequestBody.Email).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if dbUser.ValidateCode == validateEmailRequestBody.ValidateCode {
		dbUser.AccountStatus = string(enums.USER_ACCOUNT_STATUS_ACTIVE)
		dbUser.ValidateCode = ""

		if err := database.DB.Save(&dbUser).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := "Email validation successful. Your account has been validated."
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	} else {
		http.Error(w, "Invalid validation code", http.StatusBadRequest)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userRequest request_models.LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", userRequest.Email).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if dbUser.Password != userRequest.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// set device token
	dbUser.DeviceToken = userRequest.DeviceToken
	if err := database.DB.Save(&dbUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &request_models.LoginClaims{
		ID:    dbUser.ID,
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(dbUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"token": tokenString,
		"user":  string(jsonUser),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var userRequest request_models.LogoutRequestBody
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", userRequest.Email).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if dbUser.Password != userRequest.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// set device token
	dbUser.DeviceToken = ""
	if err := database.DB.Save(&dbUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Logout successful."
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
