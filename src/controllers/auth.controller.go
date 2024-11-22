package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"
	utils "caapp-server/src/utils"
)

var JwtKey = []byte("20204697")

func Register(c *gin.Context) {
	var req request_models.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi lay request"})
		return
	}

	existingUser := db_models.User{}
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_000002"})
		return
	}

	var emailValidateCode db_models.EmailValidateCode

	if err := database.DB.Where("email = ?", req.Email).First(&emailValidateCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			validateCode := utils.GenerateRandomCode(6)
			newEmailValidateCode := db_models.EmailValidateCode{
				Email:        req.Email,
				ValidateCode: validateCode,
				CreateAt:     time.Now(),
				ExpireAt:     time.Now().Add(5 * time.Minute),
			}
			if createErr := database.DB.Create(&newEmailValidateCode).Error; createErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_xxxxxx"})
				return
			}

			//send email with validate code
			header := "Validate Your Email"
			body := "Validate code:" + validateCode
			utils.SendEmail(req.Email, header, body)
		} else {
			// Xử lý các lỗi khác ngoài ErrRecordNotFound
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_xxxxxx"})
		}
	} else {
		emailValidateCode.ValidateCode = utils.GenerateRandomCode(6)
		emailValidateCode.CreateAt = time.Now()
		emailValidateCode.ExpireAt = time.Now().Add(5 * time.Minute)

		if updateErr := database.DB.Save(&emailValidateCode).Error; updateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_xxxxxx"})
			return
		}

		//send email with validate code
		header := "Validate Your Email"
		body := "Validate code:" + emailValidateCode.ValidateCode
		utils.SendEmail(req.Email, header, body)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ResendValidateCode(c *gin.Context) {
	var req request_models.ResendValidateCodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi lay request"})
		return
	}

	var emailValidateCode db_models.EmailValidateCode

	if err := database.DB.Where("email = ?", req.Email).First(&emailValidateCode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi tim kiem ban ghi"})
		return
	}

	emailValidateCode.ValidateCode = utils.GenerateRandomCode(6)
	emailValidateCode.CreateAt = time.Now()
	emailValidateCode.ExpireAt = time.Now().Add(5 * time.Minute)

	if updateErr := database.DB.Save(&emailValidateCode).Error; updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi cap nhat ban ghi email validate code"})
		return
	}

	//send email with validate code
	header := "Validate Your Email"
	body := "Validate code:" + emailValidateCode.ValidateCode
	utils.SendEmail(req.Email, header, body)

	c.JSON(http.StatusOK, gin.H{})
}

func ValidateEmail(c *gin.Context) {
	var req request_models.ValidateEmailRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_xxxxxx"})
		return
	}

	emailValidateCode := db_models.EmailValidateCode{}
	if err := database.DB.Where("email = ?", req.Email).First(&emailValidateCode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "khong tim thay email tuong ung"})
		return
	}

	// kiem tra validate con hieu luc
	if emailValidateCode.ExpireAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "validate code het thoi han"})
		return
	}

	if emailValidateCode.ValidateCode != req.ValidateCode {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "sai ma xac thuc"})
		return
	}

	var newUser db_models.User
	newUser.Email = req.Email
	newUser.Password = req.Password
	newUser.AccountStatus = "0"
	newUser.DateOfBirth = time.Now()
	newUser.CreateAt = time.Now()
	newUser.LastActive = time.Now()
	newUser.LastUpdate = time.Now()

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "loi tao moi nguoi dung vao db"})
		return
	}

	// Xóa bản ghi emailValidateCode
	if err := database.DB.Delete(&emailValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "loi_xoa_email_validate_code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	var userRequest request_models.LoginRequestBody
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000008"})
		return
	}

	var dbUser db_models.User
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

	// jsonUser, err := json.Marshal(dbUser)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000011"})
	// 	return
	// }

	// responseData := map[string]interface{}{
	// 	"token":   tokenString,
	// 	"profile": string(jsonUser),
	// }

	var res responce_models.LoginResponse
	res.UserID = dbUser.ID
	res.Token = tokenString

	c.JSON(http.StatusOK, res)
}

func Logout(c *gin.Context) {
	var userRequest request_models.LogoutRequestBody
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000012"})
		return
	}

	var dbUser db_models.User
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
