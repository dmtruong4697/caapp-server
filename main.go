package main

import (
	"caapp-server/src/database"
	"caapp-server/src/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	r := routes.SetupRouter()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the server!")
	})

	port := ":" + os.Getenv("PORT")

	err = r.Run(port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
