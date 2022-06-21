package database

import (
	"lifelog/models"
	"testing"

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
	// テスト環境のDBに接続
	db := DataBaseConnect(SetPathOption("../.testenv"))
	//	dbRepository := DataBaseReopsitory{DB: db}
	suite.db = db
	// DBのマイグレーション
	suite.db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{})
}

// テスト終了時の処理
func (suite *DataBaseTestSuite) TearDownTest() {
	// DBのテーブルを削除
	suite.db.Migrator().DropTable(&models.User{}, &models.LifeLog{}, &models.Appointment{})
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
		user := models.User{
			Aud:  "testaud",
			Name: "testname",
		}
		err := suite.db.Create(&user).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
	suite.Run("Lifelog Create", func() {
		user := models.User{
			Aud:  "lifelogtestaud",
			Name: "lifelogtestname",
		}
		suite.db.Create(&user)
		lifelog := models.LifeLog{
			UserId: user.ID,
			Name:   "testname",
		}
		err := suite.db.Create(&lifelog).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
	suite.Run("Appointment Create", func() {
		user := models.User{
			Aud:  "apptestaud",
			Name: "apptestname",
		}
		suite.db.Create(&user)
		lifelog := models.LifeLog{
			UserId: user.ID,
			Name:   "apptestname",
		}
		suite.db.Create(&lifelog)
		appointment := models.Appointment{
			LifeLogId: lifelog.ID,
		}
		err := suite.db.Create(&appointment).Error
		if err != nil {
			suite.Fail(err.Error())
		}
	})
}
