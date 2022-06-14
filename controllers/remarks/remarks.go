package remarks

import (
	"lifelog/database"
	"lifelog/models"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	now := time.Now()
	ctx.HTML(http.StatusOK, "remarks_new.html", gin.H{
		"profile": profile,
		"date":    now.Format("2006/01/02"),
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
		Start: "24:00",
		End:   "25:00",
		Class: "remarks",
	}

	lifelog := models.LifeLog{}

	// Lifelogの備考を取得
	db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", ctx.PostForm("date")).First(&lifelog)

	lifelog.Appointments = append(lifelog.Appointments, appointment)

	db.Save(&lifelog)
	ctx.Redirect(http.StatusMovedPermanently, "/lifelog")
}

func getUserInfo(map_profile map[string]interface{}) models.User {
	user := models.User{}
	user.Aud = map_profile["aud"].(string)
	user.Name = map_profile["name"].(string)
	return user
}
