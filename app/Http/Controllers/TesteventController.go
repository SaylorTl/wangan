package Controllers

import (
	"github.com/gin-gonic/gin"
	"wangxin2.0/app/Events"
)

func TestEvent(c *gin.Context) {
	// 测试事件监听
	Events.DemoQueue{}.Handle("777778989898")
}
