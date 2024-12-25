package controllers

import (
	"errors"
	"net/http"
	"os"
	"strings"
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

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func Register(c *gin.Context) {
	var req request_models.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_001001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	existingUser := db_models.User{}
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error_code": "api_error_409_001002"})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_001003"})
				return
			}

			//send email with validate code
			header := "Validate Your Email"
			body := "Validate code:" + validateCode
			utils.SendEmail(req.Email, header, body)
		} else {
			// Xử lý các lỗi khác ngoài ErrRecordNotFound
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_001004"})
		}
	} else {
		emailValidateCode.ValidateCode = utils.GenerateRandomCode(6)
		emailValidateCode.CreateAt = time.Now()
		emailValidateCode.ExpireAt = time.Now().Add(5 * time.Minute)

		if updateErr := database.DB.Save(&emailValidateCode).Error; updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_001005"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_002001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	var emailValidateCode db_models.EmailValidateCode

	if err := database.DB.Where("email = ?", req.Email).First(&emailValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_002002"})
		return
	}

	emailValidateCode.ValidateCode = utils.GenerateRandomCode(6)
	emailValidateCode.CreateAt = time.Now()
	emailValidateCode.ExpireAt = time.Now().Add(5 * time.Minute)

	if updateErr := database.DB.Save(&emailValidateCode).Error; updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_002003"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_003001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	emailValidateCode := db_models.EmailValidateCode{}
	if err := database.DB.Where("email = ?", req.Email).First(&emailValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_003002"})
		return
	}

	// kiem tra validate con hieu luc
	if emailValidateCode.ExpireAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_003003"})
		return
	}

	if emailValidateCode.ValidateCode != req.ValidateCode {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_003004"})
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
	newUser.AvatarImage = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT5bP009mhqzNrdKcPPnWSSkhtIu1FPSmeI8iSJDpTs4B4oDv--N_qM7y2xzZoF2uBX8mI&usqp=CAU"

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_003005"})
		return
	}

	// Xóa bản ghi emailValidateCode
	if err := database.DB.Delete(&emailValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_003006"})
		return
	}

	// create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &request_models.LoginClaims{
		ID:    newUser.ID,
		Email: newUser.Email,
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

	var res responce_models.LoginResponse
	res.UserID = newUser.ID
	res.Token = tokenString

	c.JSON(http.StatusOK, res)
}

func Login(c *gin.Context) {
	var userRequest request_models.LoginRequestBody
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_004001"})
		return
	}

	userRequest.Email = strings.ToLower(userRequest.Email)

	var dbUser db_models.User
	if err := database.DB.Where("email = ?", userRequest.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_004002"})
		return
	}

	if dbUser.Password != userRequest.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_004003"})
		return
	}

	// set device token
	dbUser.DeviceToken = userRequest.DeviceToken
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_004004"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_004005"})
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
	currentUserID := c.MustGet("id").(uint)
	// var userRequest request_models.LogoutRequestBody
	// if err := c.BindJSON(&userRequest); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000012"})
	// 	return
	// }

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", currentUserID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_005001"})
		return
	}

	// if dbUser.Password != userRequest.Password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000033"})
	// 	return
	// }

	// set device token
	dbUser.DeviceToken = ""
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_005002"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ForgotPassword(c *gin.Context) {
	var req request_models.ForgotPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_006001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	var dbUser db_models.User
	if err := database.DB.Where("email = ?", req.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_006002"})
		return
	}

	var forgotPasswordVaidateCode db_models.ForgotPasswordValidateCode
	if err := database.DB.Where("email = ?", req.Email).First(&forgotPasswordVaidateCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			validateCode := utils.GenerateRandomCode(6)
			newForgotPasswordValidateCode := db_models.ForgotPasswordValidateCode{
				Email:        req.Email,
				ValidateCode: validateCode,
				CreateAt:     time.Now(),
				ExpireAt:     time.Now().Add(5 * time.Minute),
			}
			if createErr := database.DB.Create(&newForgotPasswordValidateCode).Error; createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_006003"})
				return
			}

			//send email with validate code
			header := "Forgot password"
			body := "Validate code:" + validateCode
			utils.SendEmail(req.Email, header, body)
		} else {
			// Xử lý các lỗi khác ngoài ErrRecordNotFound
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_006004"})
		}
	} else {
		forgotPasswordVaidateCode.ValidateCode = utils.GenerateRandomCode(6)
		forgotPasswordVaidateCode.CreateAt = time.Now()
		forgotPasswordVaidateCode.ExpireAt = time.Now().Add(5 * time.Minute)

		if updateErr := database.DB.Save(&forgotPasswordVaidateCode).Error; updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_006005"})
			return
		}

		//send email with validate code
		header := "Forgot password"
		body := "Validate code:" + forgotPasswordVaidateCode.ValidateCode
		utils.SendEmail(req.Email, header, body)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ForgotPasswordValidate(c *gin.Context) {
	var req request_models.ForgotPasswordValidateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_007001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	forgotPasswordValidateCode := db_models.ForgotPasswordValidateCode{}
	if err := database.DB.Where("email = ?", req.Email).First(&forgotPasswordValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_007002"})
		return
	}

	// kiem tra validate con hieu luc
	if forgotPasswordValidateCode.ExpireAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_007003"})
		return
	}

	if forgotPasswordValidateCode.ValidateCode != req.ValidateCode {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_007004"})
		return
	}

	// Xóa bản ghi emailValidateCode
	if err := database.DB.Delete(&forgotPasswordValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_007005"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ForgotPasswordChangePassword(c *gin.Context) {
	var req request_models.ForgotPasswordChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_008001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	var dbUser db_models.User
	if err := database.DB.Where("email = ?", req.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_008002"})
		return
	}

	dbUser.Password = req.Password
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_008003"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ResendForgotPasswordValidateCode(c *gin.Context) {
	var req request_models.ResendForgotPasswordValidateCodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_009001"})
		return
	}

	req.Email = strings.ToLower(req.Email)

	var forgotPasswordValidateCode db_models.ForgotPasswordValidateCode

	if err := database.DB.Where("email = ?", req.Email).First(&forgotPasswordValidateCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_009002"})
		return
	}

	forgotPasswordValidateCode.ValidateCode = utils.GenerateRandomCode(6)
	forgotPasswordValidateCode.CreateAt = time.Now()
	forgotPasswordValidateCode.ExpireAt = time.Now().Add(5 * time.Minute)

	if updateErr := database.DB.Save(&forgotPasswordValidateCode).Error; updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "loi cap nhat ban ghi email validate code"})
		return
	}

	//send email with validate code
	header := "Forgot password"
	body := "Validate code:" + forgotPasswordValidateCode.ValidateCode
	utils.SendEmail(req.Email, header, body)

	c.JSON(http.StatusOK, gin.H{})
}
