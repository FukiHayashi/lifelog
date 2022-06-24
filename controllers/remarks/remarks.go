package remarks

import (
	"errors"
	"fmt"
	"lifelog/database"
	. "lifelog/helpers"
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
	if createRemarks(ctx) != nil {
		ctx.HTML(http.StatusBadRequest, "remarks_new.html", gin.H{
			"msg":     "未入力の項目があります",
			"status":  "error",
			"profile": profile,
			"date":    time.Now().Format("2006/01/02"),
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
	}
}

func createRemarks(ctx *gin.Context) error {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	// Connect to DB
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("aud = ?", user.Aud).First(&user)

	// 入力値を取得
	appointment := models.Appointment{
		Title: GetStringPointer(ctx.PostForm("title")),
		Start: GetStringPointer("24:00"),
		End:   GetStringPointer("25:00"),
		Class: "remarks",
	}

	// titleがnilになった場合のエラー処理
	if appointment.Title == nil {
		return errors.New("FORM INPUT ERROR")
	}

	lifelog := models.LifeLog{}

	// 入力された日付のLifelogを取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", ctx.PostForm("date")).First(&lifelog)

	// 新規と追記の処理分け
	remarks := models.Appointment{}
	var error error
	if err := db.Where("life_log_id = ?", lifelog.ID).Where("class = ?", "remarks").First(&remarks).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 備考が新規の場合
		lifelog.Appointments = append(lifelog.Appointments, appointment)
		error = db.Save(&lifelog).Error
	} else {
		// 追記の場合
		remarks.Title = GetStringPointer(*remarks.Title + "," + *appointment.Title)
		fmt.Println(remarks)
		error = db.Save(&remarks).Error
	}
	return error
}

func EditHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	db := database.DataBaseConnect()

	remarks := models.Appointment{}
	lifelog := models.LifeLog{}
	db.Where("id = ?", ctx.Param("remarksId")).First(&remarks)
	db.Where("id = ?", remarks.LifeLogId).First(&lifelog)
	ctx.HTML(http.StatusOK, "remarks_edit.html", gin.H{
		"profile":             profile,
		"remarks_update_path": "/remarks/update/" + ctx.Param("remarksId"),
		"remarks_delete_path": "/remarks/delete/" + ctx.Param("remarksId"),
		"title":               remarks.Title,
		"date":                lifelog.LoggedAt.Format("2006/01/02"),
	})
}

func UpdateHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	// Connect to DB
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("aud = ?", user.Aud).First(&user)

	lifelog := models.LifeLog{}
	// 入力された日付のLifelogを取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", ctx.PostForm("date")).First(&lifelog)
	// 該当のremarksを取得
	remarks := models.Appointment{}
	db.Where(&models.Appointment{LifeLogId: lifelog.ID}).Where("id = ?", ctx.Param("remarksId")).First(&remarks)

	// 値を更新
	remarks.Title = GetStringPointer(ctx.PostForm("title"))
	db.Save(&remarks)

	ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
}

func DeleteHandler(ctx *gin.Context) {
	deleteRemarks(ctx)
}

func deleteRemarks(ctx *gin.Context) {
	db := database.DataBaseConnect()
	remarks := models.Appointment{}
	db.Where("id = ?", ctx.Param("remarksId")).First(&remarks)
	db.Delete(&remarks)
}

// profileからユーザ情報を取得する
func getUserInfo(map_profile map[string]interface{}) models.User {
	user := models.User{}
	user.Aud = GetStringPointer(map_profile["aud"].(string))
	user.Name = map_profile["name"].(string)
	return user
}
