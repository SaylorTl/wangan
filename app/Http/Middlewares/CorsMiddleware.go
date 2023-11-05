package Middlewares

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
)

// 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		Utils.InitLog()
		method := c.Request.Method
		if strings.ToLower(method) == "get" {
			getParams := GetQueryParams(c)
			getStrParam, _ := json.Marshal(getParams)
			log.Print("请求参数get:", string(getStrParam))
		}
		if strings.ToLower(method) == "post" || strings.ToLower(method) == "delete" || strings.ToLower(method) == "put" {
			postParams, _ := GetPostFormParams(c)
			postStrParam, _ := json.Marshal(postParams)
			log.Print("请求参数post:", string(postStrParam))
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
		c.Set("content-type", "application/json")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		defer func() {
			if err := recover(); err != nil {
				log.Println("Exception", err)
				Constant.Error(c, err)
				return
			}
		}()
		c.Next()
	}
}

func GetQueryParams(c *gin.Context) map[string]any {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]any, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}

func GetPostFormParams(c *gin.Context) (map[string]any, error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return nil, err
		}
	}
	var postMap = make(map[string]any, len(c.Request.PostForm))
	for k, v := range c.Request.PostForm {
		if len(v) > 1 {
			postMap[k] = v
		} else if len(v) == 1 {
			postMap[k] = v[0]
		}
	}
	return postMap, nil
}
