package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"lifelog/database"
	"lifelog/models"
	"lifelog/platform/authenticator"
	"lifelog/platform/router"
)

func main() {
	f := flag.String("env", ".env", "Set .env file path")
	flag.Parse()
	if err := godotenv.Load(*f); err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("認証器の生成に失敗しました: %v", err)
	}

	// Connect to DB
	db := database.DataBaseConnect()
	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
	dbc, _ := db.DB()
	defer dbc.Close()

	rtr := router.New(auth)

	log.Print("Server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("HTTPサーバーの起動時にエラーが発生しました: %v", err)
	}
}
