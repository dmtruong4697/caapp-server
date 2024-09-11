package middlewares

import (
	"caapp-server/src/controllers"
	request_models "caapp-server/src/models/request_models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		// fmt.Print(tokenStr)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token not provided"})
			return
		}
		claims := &request_models.LoginClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return controllers.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized ErrSignatureInvalid"})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Invalid"})
			return
		}

		c.Set("id", claims.ID)
		c.Next()
	}
}
