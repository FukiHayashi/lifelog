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
	"lifelog/controllers/logout"
	"lifelog/controllers/remarks"
	"lifelog/platform/authenticator"
	"lifelog/platform/render"
)

// New はroutesを登録し、ルーターを返す
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// セッション初期設定
	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// HTMLテンプレートパス設定
	router.HTMLRender = render.CreateRenderTemplates()

	router.Static("./assets", "./assets")

	// 無名ハンドラ
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	// ログイン処理のハンドラ
	router.GET("/login", login.Handler(auth))
	// コールバックハンドラ
	router.GET("/callback", callback.Handler(auth))

	// ライフログ画面のハンドラ
	auth_group := router.Group("/", authRequired())
	{
		lifelog_group := auth_group.Group("/lifelog")
		{
			lifelog_group.GET("/", lifelog.Handler)
			lifelog_group.GET("/new", lifelog.NewHandler)
			lifelog_group.POST("/create", lifelog.CreateHandler)
			lifelog_group.GET("/:month", lifelog.MonthlyHandler)
			lifelog_group.GET("/edit/:appointmentId", lifelog.EditHandler)
			lifelog_group.POST("/update/:appointmentId", lifelog.UpdateHandler)
			lifelog_group.DELETE("/delete/:appointmentId", lifelog.DeleteHandler)
		}
		remarks_group := auth_group.Group("/remarks")
		{
			remarks_group.GET("/new", remarks.NewHandler)
			remarks_group.POST("/create", remarks.CreateHandler)
			remarks_group.GET("/edit/:remarksId", remarks.EditHandler)
			remarks_group.POST("/update/:remarksId", remarks.UpdateHandler)
			remarks_group.DELETE("/delete/:remarksId", remarks.DeleteHandler)
		}
		auth_group.DELETE("/logout", logout.Handler)
	}

	return router
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		profile := session.Get("profile")
		if profile == nil {
			// 未ログインの場合403を返す
			c.Status(http.StatusForbidden)
			c.Abort()
		}
		c.Next()
	}
}
