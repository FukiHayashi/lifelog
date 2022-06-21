package database

import (
	"fmt"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"

	"os"

	"log"
)

// pathオプション構造体
type PathOptions struct {
	path string
}

type option func(*PathOptions)

// pathオプションを設定する関数
func SetPathOption(path string) option {
	return func(po *PathOptions) {
		po.path = path
	}
}

func DataBaseConnect(opts ...option) *gorm.DB {
	// デフォルトパラメータを定義
	g := &PathOptions{
		path: ".env",
	}
	// ユーザから渡された値だけ上書き
	for _, opt := range opts {
		opt(g)
	}

	env := godotenv.Load(g.path) // Load env file
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
