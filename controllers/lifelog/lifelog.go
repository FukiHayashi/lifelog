package lifelog

import (
	"encoding/json"
	"errors"
	"lifelog/database"
	. "lifelog/helpers"
	"lifelog/models"
	"net/http"
	"strings"
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
	if err := db.Where("sub = ?", map_profile["sub"].(string)).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user = getUserInfo(map_profile)
		db.Create(&user)
	}
	// 月のデータが無い場合、その月のカレンダーを作成
	now := time.Now()
	if err := db.Where("name = ?", now.Format("2006/01/02")).Where("user_id", user.ID).First(&models.LifeLog{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		lifelogs := []models.LifeLog{}
		name_date := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		for name_date.Month() == now.Month() {
			lifelog_name := name_date.Format("2006/01/02")
			lifelogs = append(lifelogs, models.LifeLog{
				UserId:   user.ID,
				LoggedAt: name_date,
				Name:     &lifelog_name,
			})
			name_date = name_date.AddDate(0, 0, 1)
		}
		db.Create(&lifelogs)
	}

	// 表示するデータを取得
	lifelogs := []models.LifeLog{}
	db.Preload("Appointments").Preload("Remarks").Where(&models.LifeLog{UserId: user.ID}).Where("name like ?", now.Format("2006/01")+"/%").Order("name").Find(&lifelogs)
	schedulerjs_list, _ := json.Marshal(lifelogs)

	// DBに存在する月データを取得する
	var monthes []string
	db.Model(&models.LifeLog{}).Where("name like ?", "%/01").Order("name").Pluck("name", &monthes)
	for i, m := range monthes {
		t, _ := time.Parse("2006/01/02", m)
		monthes[i] = t.Format("2006-01")
	}

	ctx.HTML(http.StatusOK, "lifelog_index.html", gin.H{
		"profile":          profile,
		"months":           monthes,
		"this_month":       now.Format("2006-01"),
		"schedulerjs_list": string(schedulerjs_list),
	})

	defer database.DataBaseClose(db)
}

func MonthlyHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	// Connect to DB
	db := database.DataBaseConnect()
	map_profile := profile.(map[string]interface{})
	user := models.User{}
	// ユーザ情報を取得
	db.Where("sub = ?", map_profile["sub"].(string)).First(&user)

	// 表示するデータを取得
	lifelogs := []models.LifeLog{}
	db.Preload("Appointments").Preload("Remarks").Where(&models.LifeLog{UserId: user.ID}).Where("name like ?", strings.Replace(ctx.Param("month"), "-", "/", -1)+"/%").Order("name").Find(&lifelogs)
	schedulerjs_list, _ := json.Marshal(lifelogs)

	// DBに存在する月データを取得する
	var monthes []string
	db.Model(&models.LifeLog{}).Where("name like ?", "%/01").Where("user_id = ?", user.ID).Order("name").Pluck("name", &monthes)
	for i, m := range monthes {
		t, _ := time.Parse("2006/01/02", m)
		monthes[i] = t.Format("2006-01")
	}

	ctx.HTML(http.StatusOK, "lifelog_index.html", gin.H{
		"profile":          profile,
		"months":           monthes,
		"this_month":       ctx.Param("month"),
		"schedulerjs_list": string(schedulerjs_list),
	})

	defer database.DataBaseClose(db)
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
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	now := time.Now()

	if createLifelog(ctx) != nil {
		ctx.HTML(http.StatusBadRequest, "lifelog_new.html", gin.H{
			"msg":     "未入力の項目があります",
			"status":  "error",
			"profile": profile,
			"start":   now.Format("2006/01/02 15:04"),
			"end":     now.Add(time.Minute * 30).Format("2006/01/02 15:04"),
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
	}
}

func EditHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("sub = ?", user.Sub).First(&user)
	appointment := models.Appointment{}
	lifelog := models.LifeLog{}
	db.Where("id = ?", ctx.Param("appointmentId")).First(&appointment)
	if err := db.Where("user_id = ?", user.ID).Where("id = ?", appointment.LifeLogId).First(&lifelog).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.HTML(http.StatusOK, "lifelog_edit.html", gin.H{
			"profile":             profile,
			"lifelog_update_path": "/lifelog/update/" + ctx.Param("appointmentId"),
			"lifelog_delete_path": "/lifelog/delete/" + ctx.Param("appointmentId"),
			"title":               appointment.Title,
			"start":               *lifelog.Name + " " + *appointment.Start,
			"end":                 *lifelog.Name + " " + *appointment.End,
			"class":               appointment.Class,
		})
	}

	defer database.DataBaseClose(db)
}

func UpdateHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("sub = ?", user.Sub).First(&user)
	appointment := models.Appointment{}
	lifelog := models.LifeLog{}
	db.Where("id = ?", ctx.Param("appointmentId")).First(&appointment)
	if err := db.Where("user_id = ?", user.ID).Where("id = ?", appointment.LifeLogId).First(&lifelog).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusNotFound)
	} else {
		if createLifelog(ctx) != nil {
			ctx.HTML(http.StatusBadRequest, "lifelog_edit.html", gin.H{
				"msg":                 "未入力の項目があります",
				"status":              "error",
				"profile":             profile,
				"lifelog_update_path": "/lifelog/update/" + ctx.Param("appointmentId"),
				"lifelog_delete_path": "/lifelog/delete/" + ctx.Param("appointmentId"),
				"title":               appointment.Title,
				"start":               *lifelog.Name + " " + *appointment.Start,
				"end":                 *lifelog.Name + " " + *appointment.End,
				"class":               appointment.Class,
			})
		} else {
			deleteLifelog(ctx)
			ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
		}
	}

	defer database.DataBaseClose(db)
}

func DeleteHandler(ctx *gin.Context) {
	deleteLifelog(ctx)
}

func deleteLifelog(ctx *gin.Context) {
	db := database.DataBaseConnect()
	appointment := models.Appointment{}
	db.Where("id = ?", ctx.Param("appointmentId")).First(&appointment)
	db.Delete(&appointment)

	defer database.DataBaseClose(db)
}

func createLifelog(ctx *gin.Context) error {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	// Connect to DB
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("sub = ?", user.Sub).First(&user)

	// フォーム情報を取得
	appointment := models.Appointment{
		Title: GetStringPointer(ctx.PostForm("title")),
		Start: GetStringPointer(ctx.PostForm("start")),
		End:   GetStringPointer(ctx.PostForm("end")),
		Class: ctx.PostForm("class"),
	}

	// 未入力エラー対策
	if appointment.Start == nil || appointment.End == nil {
		return errors.New("FORM NO INPUT ERROR")
	}

	lifelogs := []models.LifeLog{}
	start_time, _ := time.Parse("2006/01/02 15:04", *appointment.Start)
	end_time, _ := time.Parse("2006/01/02 15:04", *appointment.End)

	// 月のデータが無い場合、その月のカレンダーを作成
	for _, t := range []time.Time{start_time, end_time} {
		if err := db.Where("name = ?", t.Format("2006/01/02")).Where("user_id = ?", user.ID).First(&models.LifeLog{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			lifelogs := []models.LifeLog{}
			name_date := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
			for name_date.Month() == t.Month() {
				lifelog_name := name_date.Format("2006/01/02")
				lifelogs = append(lifelogs, models.LifeLog{
					UserId:   user.ID,
					LoggedAt: name_date,
					Name:     &lifelog_name,
				})
				name_date = name_date.AddDate(0, 0, 1)
			}
			db.Create(&lifelogs)
		}
	}

	// Lifelogの範囲の日付を取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("logged_at BETWEEN ? AND ?",
		time.Date(start_time.Year(), start_time.Month(), start_time.Day(), 0, 0, 0, 0, time.Local),
		time.Date(end_time.Year(), end_time.Month(), end_time.Day(), 23, 59, 59, 999, time.Local)).Order("name").Find(&lifelogs)

	// 日を跨いだ場合、appointmentを分割する
	for i := 0; i < len(lifelogs); i++ {
		if end_time.Day() != start_time.Day() {
			// 開始時刻から日の終わりまで
			start := start_time.Format("15:04")
			end := "23:59"
			lifelogs[i].Appointments = append(lifelogs[i].Appointments, models.Appointment{
				Title: appointment.Title,
				Start: &start,
				End:   &end,
				Class: appointment.Class,
			})
			start_time = time.Date(start_time.Year(), start_time.Month(), start_time.Day()+1, 0, 0, 0, 0, time.Local)
		} else {
			// 日の始まりから終了時刻まで
			lifelogs[i].Appointments = append(lifelogs[i].Appointments, models.Appointment{
				Title: appointment.Title,
				Start: GetStringPointer(start_time.Format("15:04")),
				End:   GetStringPointer(end_time.Format("15:04")),
				Class: appointment.Class,
			})
			break
		}
	}

	defer database.DataBaseClose(db)

	return db.Save(&lifelogs).Error
}

func getUserInfo(map_profile map[string]interface{}) models.User {
	sub := map_profile["sub"].(string)
	user := models.User{
		Sub:  &sub,
		Name: map_profile["name"].(string),
	}
	return user
}
