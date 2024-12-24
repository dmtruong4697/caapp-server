// package database

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	models "caapp-server/src/models/db_models"
// 	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func Connect() {
// 	// err := godotenv.Load(".env")
// 	// if err != nil {
// 	// 	log.Fatal("Error loading .env file")
// 	// }

// 	// env := os.Getenv("APP_ENV")
// 	// envFile := fmt.Sprintf(".env.%s", env)
// 	// err := godotenv.Load(envFile)
// 	// if err != nil {
// 	// 	log.Fatalf("Error loading %s file: %v", envFile, err)
// 	// }

// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

// 	var dbErr error
// 	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if dbErr != nil {
// 		log.Fatal("Failed to connect to database:", dbErr)
// 	}

// 	// Using the database 'caapp'
// 	useDB := fmt.Sprintf("USE %s", dbName)
// 	if dbErr := DB.Exec(useDB).Error; dbErr != nil {
// 		log.Fatal("Failed to select database:", dbErr)
// 	}

// 	// Automatically migrate schema
// 	DB.AutoMigrate(&models.User{})
// 	DB.AutoMigrate(&models.Message{})
// 	DB.AutoMigrate(&models.Channel{})
// 	DB.AutoMigrate(&models.ChannelMember{})
// 	DB.AutoMigrate(&models.Media{})
// 	DB.AutoMigrate(&models.Friend{})
// 	DB.AutoMigrate(&models.FriendRequest{})
// 	DB.AutoMigrate(&models.EmailValidateCode{})
// 	DB.AutoMigrate(&models.ForgotPasswordValidateCode{})

// 	// rc module
// 	DB.AutoMigrate(&rcdbmodels.MQueue{})
// 	DB.AutoMigrate(&rcdbmodels.FQueue{})
// 	DB.AutoMigrate(&rcdbmodels.RCChannel{})
// 	DB.AutoMigrate(&rcdbmodels.RCChannelMember{})
// 	DB.AutoMigrate(&rcdbmodels.RCMessage{})
// 	DB.AutoMigrate(&rcdbmodels.GeneralQueueUser{})
// 	DB.AutoMigrate(&rcdbmodels.RCMedia{})

//		fmt.Println("Connected to", dbName, "database...")
//	}
package database

import (
	"fmt"
	"log"
	"os"

	models "caapp-server/src/models/db_models"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		log.Fatal("Failed to connect to database:", dbErr)
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Channel{},
		&models.ChannelMember{},
		&models.Media{},
		&models.Friend{},
		&models.FriendRequest{},
		&models.EmailValidateCode{},
		&models.ForgotPasswordValidateCode{},
		&rcdbmodels.MQueue{},
		&rcdbmodels.FQueue{},
		&rcdbmodels.RCChannel{},
		&rcdbmodels.RCChannelMember{},
		&rcdbmodels.RCMessage{},
		&rcdbmodels.GeneralQueueUser{},
		&rcdbmodels.RCMedia{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	fmt.Println("Connected to", dbName, "database...")
}
