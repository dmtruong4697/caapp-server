package oauth2

import (
	"caapp-server/src/controllers"
	"caapp-server/src/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// luong khac (server thuc hien dang nhap)
var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	RedirectURL:  os.Getenv("REDIRECT_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func VerifyGoogleToken(idToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func HandleGoogleLogin(c *gin.Context) {
	var requestData struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		fmt.Println("Failed to decode request info")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
		return
	}

	userData, err := VerifyGoogleToken(requestData.Token)
	if err != nil {
		fmt.Println("loi verify gg token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "loi verify gg token"})
		return
	}

	googleID := userData["sub"].(string)
	email := userData["email"].(string)
	// name := userData["name"].(string)

	var res responce_models.LoginResponse

	var user db_models.User
	if err := database.DB.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var newUser db_models.User
			newUser.CreateAt = time.Now()
			newUser.Email = email
			newUser.VerificationStatus = "0"
			newUser.GoogleID = googleID
			newUser.AccountStatus = "0"
			newUser.DateOfBirth = time.Now()
			newUser.LastUpdate = time.Now()
			newUser.LastActive = time.Now()

			if createErr := database.DB.Create(&newUser).Error; createErr != nil {
				fmt.Println("api_error_401_xxxxxx")
				c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_xxxxxx"})
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
			tokenString, err := token.SignedString(controllers.JwtKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000010"})
				return
			}

			res.Token = tokenString
			res.UserID = newUser.ID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_401_xxxxxx"})
		}
	} else {
		// create JWT token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &request_models.LoginClaims{
			ID:    user.ID,
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(controllers.JwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000010"})
			return
		}

		res.Token = tokenString
		res.UserID = user.ID
	}

	c.JSON(http.StatusOK, res)
}
