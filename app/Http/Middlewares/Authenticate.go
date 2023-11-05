package Middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/databases"
)

var identityKey = "id"

func Auth() *jwt.GinJWTMiddleware {
	conf := databases.Conf
	if databases.Conf == nil {
		absoulePath, _ := os.Getwd()
		parentPath := filepath.Dir(absoulePath)
		fileconf, fileErr := ini.Load(parentPath + "/config/my.ini")
		if fileErr != nil {
			log.Fatal(Constant.CONF_READ_ERROR+", err = ", fileErr)
		}
		conf = fileconf
	}
	jwt_secret := conf.Section("jwt").Key("JWT_SECRET").String()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(jwt_secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
