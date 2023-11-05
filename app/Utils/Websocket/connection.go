package Websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"sync"
	"time"
	"wangxin2.0/app/Utils"
)

var WebConnection *connection

type connection struct {
	ws   *websocket.Conn
	sc   chan Data
	data Data
	Mux  sync.RWMutex
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var wu = &websocket.Upgrader{ReadBufferSize: 512, HandshakeTimeout: 120 * time.Second,
	WriteBufferSize: 512, CheckOrigin: func(r *http.Request) bool { return true }}

func myws(r *gin.Context) {
	id := uuid.NewV4()
	client_id := id.String()
	ws, err := wu.Upgrade(r.Writer, r.Request, nil)
	if err != nil {
		return
	}
	c := &connection{sc: make(chan Data, 256), ws: ws, data: Data{}}
	c.data.ClientId = client_id
	client_list = append(client_list, c.data.ClientId)
	h.c.Store(client_id, c)
	h.r <- c
	go c.writer()
	c.reader()
	defer func() {
		client_list = del(client_list, c.data.ClientId)
		delUserInfo(c.data.ClientId)
	}()
}

func (c *connection) writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case clientData, ok := <-c.sc:
			send_client_conn, hcok := h.c.Load(clientData.ClientId)
			if hcok {
				send_client_conn.(*connection).Mux.Lock()
				send_client_conn.(*connection).ws.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					send_client_conn.(*connection).ws.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				message, _ := json.Marshal(clientData)
				err := send_client_conn.(*connection).ws.WriteMessage(websocket.TextMessage, message)
				send_client_conn.(*connection).Mux.Unlock()
				if err != nil {
					return
				}
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			data := make(map[string]string)
			data["type"] = "ping"
			message, _ := json.Marshal(data)
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

var client_list = []string{}
var user_client_list = map[int][]string{}
var client_user_list = map[string]int{}

func (c *connection) reader() {
	defer func() {
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			h.r <- c
			break
		}
		if string(message) == "pong" {
			c.ws.SetReadDeadline(time.Now().Add(pongWait))
			continue
		}
		unmarErr := json.Unmarshal(message, &c.data)
		if unmarErr != nil {
			continue
		}
		switch c.data.Type {
		case "login":
			if Utils.F.InArray(c.data.ClientId, client_list) {
				client_list = append(client_list, c.data.ClientId)
				c.data.Message = "init," + c.data.ClientId
				h.c.Store(c.data.ClientId, c)
				h.b <- c.data
			}
		case "user":
			c.data.Type = "user"
			h.b <- c.data
		case "logout":
			client_list = del(client_list, c.data.ClientId)
			delUserInfo(c.data.ClientId)
		default:
			fmt.Print("========default================")
		}
	}
}

func (c *connection) Login(user_id int, client_id string) bool {
	var cdata Data
	if _, ok := client_user_list[client_id]; !ok {
		client_user_list[client_id] = user_id
	}
	cdata.Type = "login"
	cdata.Message = "login," + client_id
	if _, ok := user_client_list[user_id]; !ok {
		var n_slice = []string{}
		user_client_list[user_id] = append(n_slice, client_id)
	} else {
		user_client_list[user_id] = append(user_client_list[user_id], client_id)
	}
	if !Utils.F.InArray(client_id, client_list) {
		return false
	}
	cdata.ClientId = client_id
	h.b <- cdata
	return true
}

func (c *connection) SendToAll(params map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("websocket sendToall出现异常", err)
		}
	}()
	for _, vv := range client_list {
		var cdata Data
		if params["Type"] == nil || params["Type"] == "" {
			params["Type"] = "user"
		}
		cdata.Type = params["Type"].(string)
		if params["Data"] == nil || params["Data"] == "" {
			params["Data"] = ""
		}
		cdata.Data = params["Data"]
		if params["Success"] == nil || params["Success"] == false {
			params["Success"] = false
		}
		cdata.Success = params["Success"].(bool)
		cdata.ClientId = vv
		if params["Message"] == nil || params["Message"] == "" {
			params["Message"] = ""
		}
		cdata.Message = params["Message"].(string)
		h.b <- cdata
	}
}

func (c *connection) SendToClient(params map[string]interface{}, client_id string) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("websocket sendToClient出现异常", err)
		}
	}()
	if !Utils.F.InArray(client_id, client_list) {
		return
	}
	var cdata Data
	cdata.Type = "user"
	if params["Data"] == nil || params["Data"] == "" {
		params["Data"] = ""
	}
	cdata.ClientId = client_id
	cdata.Data = params["Data"]
	cdata.Message = params["Message"].(string)
	h.b <- cdata
}

func (c *connection) SendToUser(params map[string]interface{}, user_id int) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("websocket sendToUser出现异常", err)
		}
	}()
	if _, ok := user_client_list[user_id]; !ok {
		return
	}
	for _, vv := range user_client_list[user_id] {
		var cdata Data
		if params["Type"] == nil || params["Type"] == "" {
			params["Type"] = "user"
		}
		cdata.Type = params["Type"].(string)
		if params["Data"] == nil || params["Data"] == "" {
			params["Data"] = ""
		}
		cdata.Data = params["Data"]
		if params["Success"] == nil || params["Success"] == false {
			params["Success"] = false
		}
		cdata.Success = params["Success"].(bool)
		cdata.ClientId = vv
		if params["Message"] == nil || params["Message"] == "" {
			params["Message"] = ""
		}
		cdata.Message = params["Message"].(string)
		h.b <- cdata
	}
}

func del(slice []string, client_id string) []string {
	for i := 0; i < len(slice); i++ {
		if client_id == slice[i] {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

func delUserInfo(client_id string) {
	if user_id, ok := client_user_list[client_id]; ok {
		if _, user_ok := user_client_list[user_id]; user_ok {
			var dstList = []string{}
			for _, vv := range user_client_list[user_id] {
				if vv != client_id {
					dstList = append(dstList, vv)
				}
				user_client_list[user_id] = dstList
			}
			if 0 == len(user_client_list[user_id]) {
				delete(user_client_list, user_id)
			}
		}
	}
	if val, ok := h.c.Load(client_id); ok {
		close(val.(*connection).sc)
	}
	h.c.Delete(client_id)
	delete(client_user_list, client_id)
}
