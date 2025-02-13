package main

import (
	"caapp-server/src/controllers"
	"caapp-server/src/database"
	"caapp-server/src/routes"
	rcroutines "caapp-server/src/routines/rc_routines"
	"caapp-server/src/ws"
	rcws "caapp-server/src/ws/rc_ws"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	env := os.Getenv("APP_ENV")
	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}

	database.Connect()

	// chat ws
	go controllers.HandleMessages()
	go ws.HandleChatListMessages()

	// rc chat ws
	go rcws.HandleRCChatMessages()

	// rc queue routine + ws
	go rcroutines.PairUser()
	go rcws.HandleWaitingQueueMessages()

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
