package lifelog

import (
	"errors"
	"lifelog/database"
	"lifelog/models"
	"log"
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

	if err := db.Where("aud = ?", map_profile["aud"].(string)).First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&models.User{Aud: map_profile["aud"].(string), Name: map_profile["name"].(string)})
		log.Print("User created")
	}
	//	db.Create(&user)
	ctx.HTML(http.StatusOK, "lifelog.html", profile)
}
