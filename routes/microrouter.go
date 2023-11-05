package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"wangxin2.0/app/Http/Controllers"
	"wangxin2.0/app/Http/Middlewares"
	"wangxin2.0/app/Http/Session"
)

func InitMicroRouter() *gin.Engine {
	router := gin.Default()
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(Middlewares.Cors())
	// 使用 session(cookie-based)
	router.Use(sessions.Sessions("wangxin-session", Sessions.Store))
	v2 := router.Group("api/v2")
	v2.Use()
	{
		v2.POST("/sendtouser", Controllers.WebsockController{}.SendToUserPhp)
		v2.POST("/sendtoall", Controllers.WebsockController{}.SendToAllPhp)
		v2.POST("/sendtoclient", Controllers.WebsockController{}.SendToClientPhp)
		v2.POST("/micro/subdomain-clue/add", Controllers.ClueSubdomainController{}.Add)
	}
	return router
}
