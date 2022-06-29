package database

import (
	"lifelog/models"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// テストスイートの構造体
type DataBaseTestSuite struct {
	suite.Suite
	db *gorm.DB
}

// テスト開始時の処理
func (suite *DataBaseTestSuite) SetupTest() {
	if err := godotenv.Load("../.testenv"); err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}
	// テスト環境のDBに接続
	db := DataBaseConnect()
	suite.db = db
	// DBのマイグレーション
	suite.db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
}

// テスト終了時の処理
func (suite *DataBaseTestSuite) TearDownTest() {
	// DBのテーブルを削除
	suite.db.Migrator().DropTable(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
	// DBから切断
	db, _ := suite.db.DB()
	db.Close()
}

// テストの実行
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(DataBaseTestSuite))
}

// Createのテスト
func (suite *DataBaseTestSuite) TestCreate() {
	suite.Run("User Create", func() {
		aud := "testaud"
		user := models.User{
			Aud:  &aud,
			Name: "testname",
		}
		err := suite.db.Create(&user).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
	suite.Run("Lifelog Create", func() {
		aud := "lifelogtestaud"
		user := models.User{
			Aud:  &aud,
			Name: "lifelogtestname",
		}
		suite.db.Create(&user)
		name := "testname"
		lifelog := models.LifeLog{
			UserId: user.ID,
			Name:   &name,
		}
		err := suite.db.Create(&lifelog).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
	suite.Run("Appointment Create", func() {
		aud := "apptestaud"
		user := models.User{
			Aud:  &aud,
			Name: "apptestname",
		}
		suite.db.Create(&user)
		name := "apptestname"
		lifelog := models.LifeLog{
			UserId: user.ID,
			Name:   &name,
		}
		suite.db.Create(&lifelog)
		title := "testtitle"
		appointment := models.Appointment{
			LifeLogId: lifelog.ID,
			Title:     &title,
		}
		err := suite.db.Create(&appointment).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
	suite.Run("Remarks Create", func() {
		aud := "rmktestaud"
		user := models.User{
			Aud:  &aud,
			Name: "rmktestname",
		}
		suite.db.Create(&user)
		name := "rmktestname"
		lifelog := models.LifeLog{
			UserId: user.ID,
			Name:   &name,
		}
		suite.db.Create(&lifelog)
		title := "testtitle"
		date := time.Now().Format("2006/01/02")
		remarks := models.Remarks{
			LifeLogId: lifelog.ID,
			Title:     &title,
			Date:      &date,
		}
		err := suite.db.Create(&remarks).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
}
