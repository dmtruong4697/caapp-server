package database

import (
	models "caapp-server/src/models/db_models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:truong123456@tcp(localhost:3306)/caapp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	//using the database 'sys'
	useSysDB := "USE caapp"
	if err := DB.Exec(useSysDB).Error; err != nil {
		log.Fatal("Failed to select database 'caapp':", err)
	}

	// Automatically migrate schema
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Message{})
	DB.AutoMigrate(&models.Channel{})
	DB.AutoMigrate(&models.ChannelMember{})
	DB.AutoMigrate(&models.Media{})
	DB.AutoMigrate(&models.Friend{})
	DB.AutoMigrate(&models.FriendRequest{})

	fmt.Println("Connected to caapp database ...")
}
