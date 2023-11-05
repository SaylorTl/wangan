package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"wangxin2.0/app/Http/Controllers"
	"wangxin2.0/app/Http/Middlewares"
	"wangxin2.0/app/Http/Session"
	"wangxin2.0/app/Utils/Websocket"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(Middlewares.Cors())
	// 使用 session(cookie-based)
	router.Use(sessions.Sessions("wangxin-session", Sessions.Store))
	v1 := router.Group("api/v2")

	AuthMiddleware := Middlewares.Auth()
	v1.Use(AuthMiddleware.MiddlewareFunc())
	{ //外部接口使用路由
		v1.GET("/world/read", Controllers.WorldController{}.Show)
		v1.POST("/assetrich/find", Controllers.ClueController{}.Add)
		v1.GET("/recommen/getAll", Controllers.ClueController{}.Lists)
		v1.GET("/recommen/getassetall", Controllers.ClueController{}.Show)
		v1.POST("/clue/del", Controllers.ClueController{}.Delete)
		v1.GET("/process-info", Controllers.ProcessController{}.UpDownProcessInfo)
		v1.GET("/initws", Controllers.WebsockController{}.InitWs)
		v1.GET("/sendtoclient", Controllers.WebsockController{}.SendToClient)
		v1.POST("/group/add", Controllers.RecommGroupController{}.Add)
		v1.GET("/group/show", Controllers.RecommGroupController{}.Show)
		v1.GET("/group/getall", Controllers.RecommGroupController{}.GetAll)
		v1.POST("/group/del", Controllers.RecommGroupController{}.Delete)
		v1.GET("/group/lists", Controllers.RecommGroupController{}.Lists)
		v1.GET("/clue/assetlists", Controllers.ClueAssetsController{}.Lists)
		v1.GET("/clue/assetshow", Controllers.ClueAssetsController{}.Show)
		v1.POST("/clue/addasset", Controllers.ClueAssetsController{}.AddAssetByClue)
		v1.POST("/clueasset/del", Controllers.ClueAssetsController{}.Delete)
		v1.GET("/assethistory/lists", Controllers.AssetsHistoryController{}.Lists)
		v1.GET("/search/all", Controllers.wangansearchController{}.Lists)
		v1.GET("/search/stats", Controllers.wangansearchController{}.Stats)
		v1.GET("/loopholes_status/check", Controllers.LoopholeController{}.Check)
		v1.POST("/goscanner", Controllers.ScanController{}.GoScanner)

		v1.POST("/subdomain-clue/add", Controllers.ClueSubdomainController{}.Add)
		v1.GET("/subdomain-clue/lists", Controllers.ClueSubdomainController{}.Lists)
		v1.GET("/subdomain-clue/show", Controllers.ClueSubdomainController{}.Show)
		v1.DELETE("/subdomain-clue/del", Controllers.ClueSubdomainController{}.Delete)

		v1.POST("/wangan-asset/import", Controllers.wangansearchController{}.Import)
		v1.POST("/wangan-asset/single-import", Controllers.wangansearchController{}.SingleImport)
		v1.GET("/wangan-asset/running", Controllers.wangansearchController{}.Running)

		v1.GET("/wangan-asset/syncfiledata", Controllers.wangansearchController{}.SyncFileData)
		v1.GET("/wangan-asset/aggregation", Controllers.wangansearchController{}.AggregationAsset)
		v1.GET("/wangan-asset/clearsubdomain", Controllers.wangansearchController{}.ClearSubdomain)
		v1.GET("/wangan-asset/gather-asset-from-scan", Controllers.wangansearchController{}.GatherAssetFromScan)
		v1.GET("/wangan-asset/taggedasset", Controllers.wangansearchController{}.TaggedAsset)

	}
	router.GET("/ws", Websocket.InitWebsocket)
	return router
}
