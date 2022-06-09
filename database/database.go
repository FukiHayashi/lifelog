package database

import (
	"fmt"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"

	"os"

	"log"
)

func DataBaseConnect() *gorm.DB {
	env := godotenv.Load() // Load env file
	if env != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Database
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	return db
}
