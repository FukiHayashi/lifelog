package lifelog

import (
	"encoding/json"
	"errors"
	"lifelog/database"
	"lifelog/models"
	"net/http"

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
		user.Aud = map_profile["aud"].(string)
		user.Name = map_profile["name"].(string)
		db.Create(&user)
	}
	// 表示するデータを取得
	lifelogs := []models.LifeLog{}
	db.Preload("Appointments").Where(&models.LifeLog{UserId: user.ID}, "user_id").Find(&lifelogs)
	schedulerjs_list, _ := json.Marshal(lifelogs)

	ctx.HTML(http.StatusOK, "lifelog.html", gin.H{
		"profile":          profile,
		"schedulerjs_list": string(schedulerjs_list),
	})
}
