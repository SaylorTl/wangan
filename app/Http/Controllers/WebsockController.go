package Controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"strconv"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils/Websocket"
)

// Booking 包含绑定和验证的数据。
type Booking struct {
	Age  int    `form:"age" binding:"required"`
	Name string `form:"name" binding:"required"`
}

type send struct {
	Content   string `form:"content" binding:"required"`
	Client_id string `form:"client_id" binding:"required"`
}

type WebsockController struct {
}

func (w WebsockController) InitWs(c *gin.Context) {
	jwt_payload := jwt.ExtractClaims(c)
	user_id := int(jwt_payload["sub"].(float64))
	client_id := c.Query("client_id")
	res := Websocket.WebConnection.Login(user_id, client_id)
	if !res {
		Constant.Error(c, "登录失败")
		return
	}
	Constant.Success(c, "登录成功!", "")
}

func (w WebsockController) SendToAll(c *gin.Context) {
	params := make(map[string]interface{})
	content := c.Query("content")
	params["Type"] = "user"
	params["Data"] = ""
	params["Message"] = content
	Websocket.WebConnection.SendToAll(params)
	Constant.Success(c, "查询成功!", "")
}

func (w WebsockController) SendToUser(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	params := make(map[string]interface{})
	content := c.Query("content")
	params["Type"] = "user"
	params["Data"] = ""
	params["Message"] = content
	Websocket.WebConnection.SendToUser(params, user_id)
	Constant.Success(c, "查询成功!", "")
}

func (w WebsockController) SendToClient(c *gin.Context) {
	params := make(map[string]interface{})
	content := c.Query("content")
	params["Type"] = "user"
	params["Data"] = content
	params["client_id"] = c.Query("client_id")
	params["Message"] = c.Query("message")
	Websocket.WebConnection.SendToClient(params, params["client_id"].(string))
	Constant.Success(c, "查询成功!", "")
}

func (w WebsockController) SendToAllPhp(c *gin.Context) {
	params := make(map[string]interface{})
	var data Websocket.Data
	log.Print("content", c.PostForm("content"))
	err := json.Unmarshal([]byte(c.PostForm("content")), &data)
	if err != nil {
		log.Fatal(err)

	}
	params["Type"] = data.Type
	params["Data"] = data.Data
	params["Success"] = data.Success
	Websocket.WebConnection.SendToAll(params)
	Constant.Success(c, "查询成功!", params)
}
func (w WebsockController) SendToUserPhp(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.PostForm("user_id"))
	params := make(map[string]interface{})
	var data Websocket.Data
	err := json.Unmarshal([]byte(c.PostForm("content")), &data)
	if err != nil {
		log.Fatal(err)

	}
	params["Type"] = data.Type
	params["Data"] = data.Data
	params["Success"] = data.Success
	Websocket.WebConnection.SendToUser(params, user_id)
	Constant.Success(c, "查询成功!", params)
}

func (w WebsockController) SendToClientPhp(c *gin.Context) {
	params := make(map[string]interface{})
	var data Websocket.Data
	err := json.Unmarshal([]byte(c.PostForm("content")), &data)
	if err != nil {
		log.Fatal(err)
	}
	params["Type"] = data.Type
	params["Data"] = data.Data
	params["Success"] = data.Success
	params["client_id"] = c.PostForm("client_id")
	params["Message"] = ""
	clientId := c.PostForm("client_id")
	Websocket.WebConnection.SendToClient(params, clientId)
	Constant.Success(c, "查询成功!", params)
}
