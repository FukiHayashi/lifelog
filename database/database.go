package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"os"
)

func DataBaseConnect() *gorm.DB {
	// Connect to Database
	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	defer sqlDB.Close()
	return db
}

func DataBaseClose(db *gorm.DB) {
	dbc, _ := db.DB()
	dbc.Close()
}
