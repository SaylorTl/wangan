package Websocket

import (
	"github.com/gin-gonic/gin"
)

func InitWebsocket(c *gin.Context) {
	go h.run()
	myws(c)
}
