package lifelog

import (
	"encoding/json"
	"errors"
	"lifelog/database"
	"lifelog/models"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	// Connect to DB
	db := database.DataBaseConnect()
	map_profile := profile.(map[string]interface{})

	user := models.User{}
	// 初めてログインする場合、ユーザを作成
	if err := db.Where("aud = ?", map_profile["aud"].(string)).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user = getUserInfo(map_profile)
		db.Create(&user)
	}
	// 月のデータが無い場合、その月のカレンダーを作成
	now := time.Now()
	if err := db.Where("name = ?", now.Format("2006/01/02")).First(&models.LifeLog{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		lifelogs := []models.LifeLog{}
		name_date := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		for name_date.Month() == now.Month() {
			lifelogs = append(lifelogs, models.LifeLog{
				UserId:   user.ID,
				LoggedAt: name_date,
				Name:     name_date.Format("2006/01/02"),
			})
			name_date = name_date.AddDate(0, 0, 1)
		}
		db.Create(&lifelogs)
	}

	// 表示するデータを取得
	lifelogs := []models.LifeLog{}
	db.Preload("Appointments").Where(&models.LifeLog{UserId: user.ID}).Find(&lifelogs)
	schedulerjs_list, _ := json.Marshal(lifelogs)

	ctx.HTML(http.StatusOK, "lifelog_index.html", gin.H{
		"profile":          profile,
		"schedulerjs_list": string(schedulerjs_list),
	})
}

func NewHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	now := time.Now()
	ctx.HTML(http.StatusOK, "lifelog_new.html", gin.H{
		"profile": profile,
		"start":   now.Format("2006/01/02 15:04"),
		"end":     now.Add(time.Minute * 30).Format("2006/01/02 15:04"),
	})
}

func CreateHandler(ctx *gin.Context) {
	createLifelog(ctx)
	ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
}

func EditHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	db := database.DataBaseConnect()

	appointment := models.Appointment{}
	lifelog := models.LifeLog{}
	db.Where("id = ?", ctx.Param("appointmentId")).First(&appointment)
	db.Where("id = ?", appointment.LifeLogId).First(&lifelog)
	ctx.HTML(http.StatusOK, "lifelog_edit.html", gin.H{
		"profile":             profile,
		"lifelog_update_path": "/lifelog/update/" + ctx.Param("appointmentId"),
		"title":               appointment.Title,
		"start":               lifelog.Name + " " + appointment.Start,
		"end":                 lifelog.Name + " " + appointment.End,
		"class":               appointment.Class,
	})
}

func UpdateHandler(ctx *gin.Context) {
	createLifelog(ctx)
	deleteLifelog(ctx)
	ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
}

func DeleteHandler(ctx *gin.Context) {
	deleteLifelog(ctx)
}

func deleteLifelog(ctx *gin.Context) {
	db := database.DataBaseConnect()
	appointment := models.Appointment{}
	db.Where("id = ?", ctx.Param("appointmentId")).First(&appointment)
	db.Delete(&appointment)
}

func createLifelog(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	// Connect to DB
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("aud = ?", user.Aud).First(&user)

	// フォーム情報を取得
	appointment := models.Appointment{
		Title: ctx.PostForm("title"),
		Start: ctx.PostForm("start"),
		End:   ctx.PostForm("end"),
		Class: ctx.PostForm("class"),
	}

	lifelogs := []models.LifeLog{}
	start_time, _ := time.Parse("2006/01/02 15:04", appointment.Start)
	end_time, _ := time.Parse("2006/01/02 15:04", appointment.End)

	// Lifelogの範囲の日付を取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("logged_at BETWEEN ? AND ?",
		time.Date(start_time.Year(), start_time.Month(), start_time.Day(), 0, 0, 0, 0, time.Local),
		time.Date(end_time.Year(), end_time.Month(), end_time.Day(), 23, 59, 59, 999, time.Local)).Find(&lifelogs)

	// 日を跨いだ場合、appointmentを分割する
	for i := 0; i < len(lifelogs); i++ {
		if end_time.Day() > start_time.Day() {
			// 開始時刻から日の終わりまで
			lifelogs[i].Appointments = append(lifelogs[i].Appointments, models.Appointment{
				Title: appointment.Title,
				Start: start_time.Format("15:04"),
				End:   "23:59",
				Class: appointment.Class,
			})
			start_time = time.Date(start_time.Year(), start_time.Month(), start_time.Day()+1, 0, 0, 0, 0, time.Local)
		} else {
			// 日の始まりから終了時刻まで
			lifelogs[i].Appointments = append(lifelogs[i].Appointments, models.Appointment{
				Title: appointment.Title,
				Start: start_time.Format("15:04"),
				End:   end_time.Format("15:04"),
				Class: appointment.Class,
			})
			break
		}
	}
	db.Save(&lifelogs)
}

func getUserInfo(map_profile map[string]interface{}) models.User {
	user := models.User{}
	user.Aud = map_profile["aud"].(string)
	user.Name = map_profile["name"].(string)
	return user
}
