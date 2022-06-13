package router

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"lifelog/controllers/callback"
	"lifelog/controllers/lifelog"
	"lifelog/controllers/login"
	"lifelog/controllers/logout"
	"lifelog/platform/authenticator"
	"lifelog/platform/render"
)

// New はroutesを登録し、ルーターを返す
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	// セッション初期設定
	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// HTMLテンプレートパス設定
	//router.LoadHTMLGlob("views/*")
	router.HTMLRender = render.CreateRenderTemplates()
	//	router.SetHTMLTemplate(render.RenderTemplate("home.html"))

	router.Static("./assets", "./assets")

	// 無名ハンドラ
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	// ログイン処理のハンドラ
	router.GET("/login", login.Handler(auth))
	// コールバックハンドラ
	router.GET("/callback", callback.Handler(auth))

	auth_user_group := router.Group("/auth")
	auth_user_group.Use(authRequired())
	{
		// ライフログ画面のハンドラ
		router.GET("/lifelog", lifelog.Handler)
		router.GET("/lifelog/new", lifelog.NewHandler)
		router.POST("/lifelog/create", lifelog.CreateHandler)
		router.GET("/logout", logout.Handler)
	}
	return router
}

// ここはうまく動きません
// 未ログイン状態で/lifelogを見ようとした時にエラーが出る（トップページにリダイレクトしたい）
// 調べても分からないので後回し
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		profile := session.Get("profile")
		fmt.Println(profile)
		if profile == nil {
			c.Redirect(http.StatusMovedPermanently, "/")
			c.Abort()
		}
		c.Next()
	}
}
