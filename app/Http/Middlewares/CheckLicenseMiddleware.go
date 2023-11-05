package Middlewares

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 跨域中间件
func Handle(c *gin.Context) {
	token := c.GetHeader("token")

	s := md5.New()
	s.Write([]byte(token))
	token = hex.EncodeToString(s.Sum(nil))
	fmt.Print(token)
}
func CheckLicense() {

}
