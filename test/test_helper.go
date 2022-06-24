package test

import (
	"lifelog/database"
	"lifelog/models"

	"gorm.io/gorm"
)

type TestHelper struct {
	DB *gorm.DB
}

// テスト開始時の処理
func (th *TestHelper) SetupTest() {
	// テスト環境のDBに接続
	db := database.DataBaseConnect()
	th.DB = db
	// DBのマイグレーション
	th.DB.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{})
}

// テスト終了時の処理
func (th *TestHelper) TearDownTest() {
	// DBのテーブルを削除
	th.DB.Migrator().DropTable(&models.User{}, &models.LifeLog{}, &models.Appointment{})
	// DBから切断
	db, _ := th.DB.DB()
	db.Close()
}
