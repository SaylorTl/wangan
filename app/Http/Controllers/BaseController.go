package Controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
)

type BaseController struct {
}

func (b BaseController) HasRole(c *gin.Context, roleData []string) bool {
	queryData := make(map[string]interface{})
	jwt_payload := jwt.ExtractClaims(c)
	queryData["id"] = int(jwt_payload["sub"].(float64))
	data, _ := Models.UserModel.FindUser(queryData)
	if Utils.F.InArray(data["username"].(string), roleData) {
		return true
	}
	return false
}

func (b BaseController) GetUserId(c *gin.Context) int {
	jwt_payload := jwt.ExtractClaims(c)
	return int(jwt_payload["sub"].(float64))
}
