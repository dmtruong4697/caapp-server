package main

import (
	"caapp-server/src/controllers"
	"caapp-server/src/database"
	"caapp-server/src/routes"
	rcroutines "caapp-server/src/routines/rc_routines"
	"caapp-server/src/ws"
	rcws "caapp-server/src/ws/rc_ws"
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

	// chat ws
	go controllers.HandleMessages()
	go ws.HandleChatListMessages()

	// rc ws
	go rcws.HandleRCChatMessages()

	// rc routine
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
