package remarks

import (
	"errors"
	"fmt"
	"lifelog/database"
	"lifelog/models"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// /remarks/new
func NewHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	// 現在時刻を初期値として画面を表示
	ctx.HTML(http.StatusOK, "remarks_new.html", gin.H{
		"profile": profile,
		"date":    time.Now().Format("2006/01/02"),
	})
}

// /remarks/create
func CreateHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	// Connect to DB
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("aud = ?", user.Aud).First(&user)

	// 入力値を取得
	appointment := models.Appointment{
		Title: ctx.PostForm("title"),
		Start: "24:00",
		End:   "25:00",
		Class: "remarks",
	}

	lifelog := models.LifeLog{}

	// 入力された日付のLifelogを取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", ctx.PostForm("date")).First(&lifelog)

	// 新規と追記の処理分け
	remark := models.Appointment{}
	if err := db.Where("life_log_id = ?", lifelog.ID).Where("class = ?", "remarks").First(&remark).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 備考が新規の場合
		lifelog.Appointments = append(lifelog.Appointments, appointment)
		db.Save(&lifelog)
	} else {
		// 追記の場合
		remark.Title = remark.Title + "," + appointment.Title
		fmt.Println(remark)
		db.Save(&remark)
	}

	ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
}

// profileからユーザ情報を取得する
func getUserInfo(map_profile map[string]interface{}) models.User {
	user := models.User{}
	user.Aud = map_profile["aud"].(string)
	user.Name = map_profile["name"].(string)
	return user
}
