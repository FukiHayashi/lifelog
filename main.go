package main

import (
	"flag"
	"log"
	"os"

	"lifelog/database"
	"lifelog/models"
	"lifelog/platform/authenticator"
	"lifelog/platform/router"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
)

func main() {
	f := flag.String("env", ".env", "Set .env file path")
	flag.Parse()
	if *f != "" {
		if err := godotenv.Load(*f); err != nil {
			log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
		}
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("認証器の生成に失敗しました: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Connect to DB
	db := database.DataBaseConnect()
	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
	dbc, _ := db.DB()
	defer dbc.Close()

	rtr := router.New(auth)

	log.Print("Server listening on " + os.Getenv("PORT"))

	rtr.Run(":" + port)
}
