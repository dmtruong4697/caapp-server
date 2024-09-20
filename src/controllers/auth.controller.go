package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	"caapp-server/src/enums"
	models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	utils "caapp-server/src/utils"
)

var JwtKey = []byte("20204697")

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser := models.User{}
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_000002"})
		return
	}

	validateCode := utils.GenerateRandomCode(6)
	user.ValidateCode = validateCode
	user.AccountStatus = string(enums.USER_ACCOUNT_STATUS_NOT_ACTIVE)

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_400_000003"})
		return
	}

	// send email with validate code
	header := "Validate Your Email"
	body := "Validate code:" + validateCode
	utils.SendEmail(user.Email, header, body)

	c.JSON(http.StatusCreated, user)
}

func ValidateEmail(c *gin.Context) {
	var validateEmailRequestBody request_models.ValidateEmailRequestBody
	if err := c.BindJSON(&validateEmailRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000004"})
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", validateEmailRequestBody.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000005"})
		return
	}

	if dbUser.ValidateCode == validateEmailRequestBody.ValidateCode {
		dbUser.AccountStatus = string(enums.USER_ACCOUNT_STATUS_ACTIVE)
		dbUser.ValidateCode = ""

		if err := database.DB.Save(&dbUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000006"})
			return
		}

		message := "Email validation successful. Your account has been validated."
		c.JSON(http.StatusOK, gin.H{"message": message})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_000007"})
	}
}

func Login(c *gin.Context) {
	var userRequest request_models.LoginRequestBody
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000008"})
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", userRequest.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000001"})
		return
	}

	if dbUser.Password != userRequest.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000001"})
		return
	}

	// set device token
	dbUser.DeviceToken = userRequest.DeviceToken
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000009"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000010"})
		return
	}

	jsonUser, err := json.Marshal(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000011"})
		return
	}

	responseData := map[string]interface{}{
		"token":   tokenString,
		"profile": string(jsonUser),
	}

	c.JSON(http.StatusOK, responseData)
}

func Logout(c *gin.Context) {
	var userRequest request_models.LogoutRequestBody
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000012"})
		return
	}

	var dbUser models.User
	if err := database.DB.Where("email = ?", userRequest.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000013"})
		return
	}

	if dbUser.Password != userRequest.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000033"})
		return
	}

	// set device token
	dbUser.DeviceToken = ""
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000014"})
		return
	}

	message := "Logout successful."
	c.JSON(http.StatusOK, gin.H{"message": message})
}
