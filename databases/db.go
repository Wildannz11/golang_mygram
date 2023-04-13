package database

import (
	"fmt"
	"os"
	"project4/helpers"
	"project4/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// dbHost = os.Getenv("DB_HOST")
	// dbPort = os.Getenv("DB_PORT")
	// dbName = os.Getenv("DB_NAME")
	// dbUser = os.Getenv("DB_USER")
	// dbPassword = os.Getenv("DB_PASSWORD")
	// debugMode = os.Getenv("DEBUG_MODE")
	db  *gorm.DB
	err error
)

// func loadDebugModeEnv() string {
// 	load := godotenv.Load()
//     if load != nil {
//         fmt.Println("Error loading .env file")
//         os.Exit(1)
//     }

// 	debugMode := os.Getenv("DEBUG_MODE")
// 	return debugMode
// }

func StartDB() {
	helpers.LoadEnv()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	debugMode := os.Getenv("DEBUG_MODE")

	conn := fmt.Sprintf("host=%s  user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	fmt.Println("Ini database ", conn)
	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully Connected to Database: ", dbName)
	if debugMode == "true" {
		db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
	}
	
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	helpers.LoadEnv()
	debugMode := os.Getenv("DEBUG_MODE")
	if debugMode == "true" {
		return db.Debug()
	}
	fmt.Println("ini load debug mode ", debugMode)
	return db
}
