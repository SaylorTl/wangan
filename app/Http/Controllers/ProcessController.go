package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

type IdentifierS struct {
	Identifier string `form:"identifier" binding:"required"`
}
type ProcessController struct {
}

func (p ProcessController) UpDownProcessInfo(c *gin.Context) {
	var identifierS IdentifierS
	if err := c.ShouldBind(&identifierS); err != nil {
		Constant.Error(c, err.Error())
		return
	}
	identifier := c.Query("identifier")
	value, _ := databases.Redis.Get(Constant.FIND_ASSET_PROCESS + ":" + identifier).Result()
	if value == "" {
		Constant.Success(c, "进度查询成功!", "0%")
		return
	}
	defaultValue, _ := strconv.ParseFloat(value, 64)
	defaultMapProcess := Utils.F.If(value == "", 0, defaultValue*100)
	defaultProcess := fmt.Sprintf("%.0f", defaultMapProcess)
	Constant.Success(c, "进度查询成功!", defaultProcess+"%")
	return
}
