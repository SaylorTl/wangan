package Services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var D01ICiiApiService d01iciiapiservice

type d01iciiapiservice struct {
}

// 封装D01查询组件
func (a d01iciiapiservice) GetAssetEnrichmentData(strParam []byte) (data map[string]interface{}) {
	conf := databases.Conf
	base_url := a.bindUrlParam(conf.Section("d01").Key("D01_I_CII_API_BASE_URL").String())
	//构建http请求
	client := &http.Client{}
	response, err := client.Post(base_url, "application/json", bytes.NewBuffer([]byte(strParam)))
	if err != nil {
		log.Println("D01 get error", err)
		panic(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	returnData := make(map[string]interface{})
	unmarErr := json.Unmarshal(body, &returnData)
	if unmarErr != nil {
		return nil
	}
	return returnData
}

func (a d01iciiapiservice) bindUrlParam(api string) string {
	params := a.sign()
	return api + "?" + Utils.F.HttpBuildQuery(params, "")
}

// 封装D01查询签名模块
func (a d01iciiapiservice) sign() map[string]interface{} {
	conf := databases.Conf
	app_id, err := strconv.Atoi(conf.Section("d01").Key("APP_ID").String())
	if err != nil {
		log.Println("d01 appid查找失败", err)
		panic(err)
	}
	params := make(map[string]interface{})
	params["secret"] = conf.Section("d01").Key("SECRET").String()
	params["appId"] = GarBleService.IdToString(app_id)
	params["nonceStr"] = strconv.Itoa(rand.Intn(90000) + 9999)
	params["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	md5Str := strings.Join([]string{params["secret"].(string), params["appId"].(string), params["nonceStr"].(string),
		params["timestamp"].(string)}, "&")
	m := md5.Sum([]byte(md5Str))
	returnSign := hex.EncodeToString(m[:])
	params["sign"] = returnSign
	delete(params, "secret")
	return params
}
