package test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"wangxin2.0/app/Constant"
	"wangxin2.0/databases"
)

// map转字符串
func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

func getVerifyCodeIndustry() (code string) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	response, err := client.Get("https://10.10.11.74/api/auth/getVerifyCode")
	if err != nil {
		log.Println("get error", err)
	}
	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		log.Println("ioutil read error")
	}
	m := make(map[string]map[string]interface{})
	json.Unmarshal(body, &m)
	return m["data"]["code_key"].(string)

}

func getJwtToken() (token string) {
	log.Println("get error1111111")
	vcode_key := getVerifyCodeIndustry()
	//初始化redis
	AbsoulePath, _ := os.Getwd()
	parentPath := filepath.Dir(AbsoulePath)
	conf, err := ini.Load(parentPath + "/config/my.ini")
	if err != nil {
		log.Fatal(Constant.CONF_READ_ERROR+", err = ", err)
	}
	if err := databases.InitRedisDb(conf); err != nil {
		log.Fatal(Constant.INIT_REDIS_ERROR+", err = ", err)
		return
	}
	vcode, _ := databases.Redis.Get("user:verify_code:" + vcode_key).Result()
	params := make(map[string]interface{})
	params["code"] = vcode
	params["code_key"] = vcode_key
	params["password"] = "5d871c087456508df661ec9cce820ac0"
	params["username"] = "admin"
	strParam, _ := json.Marshal(params)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	response, err := client.Post("https://10.10.11.74/api/auth/login", "application/json", bytes.NewReader(strParam))
	if err != nil {
		log.Println("get error", err)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	m := make(map[string]map[string]interface{})
	json.Unmarshal(body, &m)
	return m["data"]["access_token"].(string)
}

// get access controller
func Get(uri string, router *gin.Engine) *httptest.ResponseRecorder {
	jwtToken := getJwtToken()
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "Bearer "+jwtToken)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

func Post(uri string, data []string, router *gin.Engine) *httptest.ResponseRecorder {
	jwtToken := getJwtToken()
	// 构造get请求
	strParam, _ := json.Marshal(data)
	req := httptest.NewRequest("Post", uri, bytes.NewReader(strParam))
	req.Header.Add("Authorization", "Bearer "+jwtToken)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

// post access controller
func PostForm(uri string, param map[string]string, router *gin.Engine) *httptest.ResponseRecorder {
	jwtToken := getJwtToken()
	req := httptest.NewRequest("POST", uri+ParseToStr(param), nil)
	req.Header.Add("Authorization", "Bearer "+jwtToken)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	return w
}
