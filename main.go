package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"lifelog/platform/authenticator"
	"lifelog/platform/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Error Authentication: %v", err)
	}

	rtr := router.New(auth)

	log.Print("Server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("HTTPサーバーの起動時にエラーが発生しました: %v", err)
	}
}
