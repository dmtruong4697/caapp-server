package database

import (
	"fmt"
	"log"
	"os"

	models "caapp-server/src/models/db_models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	var dbErr error
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		log.Fatal("Failed to connect to database:", dbErr)
	}

	// Using the database 'caapp'
	useDB := fmt.Sprintf("USE %s", dbName)
	if dbErr := DB.Exec(useDB).Error; dbErr != nil {
		log.Fatal("Failed to select database:", dbErr)
	}

	// Automatically migrate schema
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Message{})
	DB.AutoMigrate(&models.Channel{})
	DB.AutoMigrate(&models.ChannelMember{})
	DB.AutoMigrate(&models.Media{})
	DB.AutoMigrate(&models.Friend{})
	DB.AutoMigrate(&models.FriendRequest{})

	fmt.Println("Connected to", dbName, "database...")
}
