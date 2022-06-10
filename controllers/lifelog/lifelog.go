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
		"lifelog_new_path": "/lifelog/new",
	})
}

func NewHandler(ctx *gin.Context) {
	now := time.Now()
	ctx.HTML(http.StatusOK, "lifelog_new.html", gin.H{
		"start": now.Format("2006/01/02 15:04"),
		"end":   now.Add(time.Minute * 30).Format("2006/01/02 15:04"),
	})
}

func CreateHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user := getUserInfo(profile.(map[string]interface{}))
	// Connect to DB
	db := database.DataBaseConnect()
	// ユーザの取得
	db.Where("aud = ?", user.Aud).First(&user)

	appointment := models.Appointment{
		Title: ctx.PostForm("title"),
		Start: ctx.PostForm("start"),
		End:   ctx.PostForm("end"),
		Class: ctx.PostForm("class"),
	}

	lifelogs := []models.LifeLog{}
	// 日を跨いだ場合、appointmentを分割する
	start_time, _ := time.Parse("2006/01/02 15:04", appointment.Start)
	end_time, _ := time.Parse("2006/01/02 15:04", appointment.End)

	// Lifelogの範囲の日付を取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("logged_at BETWEEN ? AND ?",
		time.Date(start_time.Year(), start_time.Month(), start_time.Day(), 0, 0, 0, 0, time.Local),
		time.Date(end_time.Year(), end_time.Month(), end_time.Day(), 23, 59, 59, 999, time.Local)).Find(&lifelogs)

	for i := 0; i < len(lifelogs); i++ {
		//		lifelog := models.LifeLog{Appointments: []models.Appointment{}}
		//		lifelog.Name = start_time.Format("2006/01/02")
		if end_time.Day() > start_time.Day() {
			lifelogs[i].Appointments = append(lifelogs[i].Appointments, models.Appointment{
				Title: appointment.Title,
				Start: start_time.Format("15:04"),
				End:   "23:59",
				Class: appointment.Class,
			})
			start_time = time.Date(start_time.Year(), start_time.Month(), start_time.Day()+1, 0, 0, 0, 0, time.Local)
			//			lifelogs = append(lifelogs, lifelog)
		} else {
			lifelogs[i].Appointments = append(lifelogs[i].Appointments, models.Appointment{
				Title: appointment.Title,
				Start: start_time.Format("15:04"),
				End:   end_time.Format("15:04"),
				Class: appointment.Class,
			})
			//			lifelogs = append(lifelogs, lifelog)
			break
		}
	}
	//	user.Lifelogs = lifelogs
	//	fmt.Println(lifelogs)
	db.Save(&lifelogs)
	ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
}

func getUserInfo(map_profile map[string]interface{}) models.User {
	user := models.User{}
	user.Aud = map_profile["aud"].(string)
	user.Name = map_profile["name"].(string)
	return user
}
