package router

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"lifelog/controllers/callback"
	"lifelog/controllers/lifelog"
	"lifelog/controllers/login"
	"lifelog/platform/authenticator"
)

// New はroutesを登録し、ルーターを返す
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	// セッション初期設定
	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// HTMLテンプレートパス設定
	router.LoadHTMLGlob("views/template/*")

	// 無名ハンドラ
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	// ログイン処理のハンドラ
	router.GET("/login", login.Handler(auth))
	// コールバックハンドラ
	router.GET("/callback", callback.Handler(auth))
	// ライフログ画面のハンドラ
	router.GET("/lifelog", lifelog.Handler)

	return router
}
