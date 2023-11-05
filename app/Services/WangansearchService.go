package Services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/thinkeridea/go-extend/exnet"
	"log"
	"strconv"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/app/Utils/Websocket"
	"wangxin2.0/databases"
)

var wangansearchService wangansearchservice

type wangansearchservice struct {
}

var current_search_ip_list = []string{}

func (f wangansearchservice) Find(queryParams map[string]interface{}) {
	databases.Redis.Set("wangansearch:currentsize", 0, 3600*time.Second)
	databases.Redis.Set("wangansearch:totalsize", 0, 3600*time.Second)
	databases.Redis.Set("wangansearch:is_finish", 1, 3600*time.Second)
	pageNum := 0
	size := 500
	f.SetProcess(float32(1))
	scroll_id := ""
	page := 0
	for {
		wangantotalTmpsize, _ := databases.Redis.Get("wangansearch:totalsize").Result()
		if "" == wangantotalTmpsize {
			break
		}
		totalsize, _ := strconv.Atoi(wangantotalTmpsize)
		if totalsize > 0 && pageNum >= (totalsize/size) {
			continue
		}
		importStatus := f.Import(Utils.F.C(queryParams), size, scroll_id)
		if false == importStatus["status"] {
			f.SetFailProcess()
			f.clearSyncAssetStatus()
			return
		}
		if true == importStatus["status"] && "" == importStatus["next"] {
			break
		}
		scroll_id = importStatus["next"].(string)
		page++
		log.Print(fmt.Sprintf("当前运行到第%v页", page), scroll_id)
	}
	for {
		runningNums := Utils.ImportAssetPool.Running()
		if 0 == runningNums {
			break
		}
	}
	//计算聚合总数
	user_id := queryParams["user_id"].(int)
	company_id := queryParams["company_id"].(int)
	totalNum := 0
	ipArr, _ := Utils.F.Unique(current_search_ip_list)
	log.Print("current_search_ip_list total size", len(ipArr.([]string)))
	for _, ipval := range ipArr.([]string) {
		Utils.ImportAssetPool.Submit(Models.AssetModel.InsertAsset(ipval, user_id, company_id))
		totalNum++
		strprogress, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (float32(totalNum)/float32(len(ipArr.([]string))*2))*100), 32)
		progress := float32(50) + float32(strprogress)
		if totalNum >= len(ipArr.([]string)) {
			f.SetProcess(float32(100))
			break
		}
		if totalNum > 500 && totalNum%500 == 0 {
			f.SetProcess(progress)
			log.Print("insert asset process", progress)
		}
		if totalNum < 500 {
			f.SetProcess(progress)
			log.Print("insert asset process", progress)
		}
	}
	f.clearSyncStatus()
	log.Print("wangan sync finish")
	//f.setProcess(float64(100))
}

func (f wangansearchservice) Import(params map[string]interface{}, size int, scroll_id string) map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Print(fmt.Sprintf("第%v页查询失败2", scroll_id), err)
			f.clearSyncAssetStatus()
			return
		}
	}()
	returnStatus := make(map[string]interface{})
	returnStatus["status"] = true
	returnStatus["next"] = ""
	data := f.retrywanganSearch(params, size, scroll_id, 1)
	if _, ok := data["next"]; !ok {
		log.Print("无next数据", data)
		returnStatus["status"] = false
		return returnStatus
	}
	if _, ok := data["results"]; !ok || (ok && 0 == len(data["results"].([]interface{}))) {
		log.Print("空result", data)
		return returnStatus
	}
	Utils.ImportAssetPool.Submit(f.insertData(data, params))
	returnStatus["next"] = data["next"].(string)
	return returnStatus
}

func (f wangansearchservice) retrywanganSearch(params map[string]interface{}, size int, scroll_id string, retryNum int) map[string]interface{} {
	if retryNum > 20 {
		return map[string]interface{}{}
	}
	qbase64 := base64.StdEncoding.EncodeToString([]byte(params["keyword"].(string)))
	queryParams := map[string]string{
		"email":   params["email"].(string),
		"key":     params["key"].(string),
		"qbase64": qbase64,
		"fields":  params["fields"].(string),
		"full":    strconv.FormatBool(params["full"].(bool)),
		"size":    strconv.Itoa(size),
	}
	if "" != scroll_id {
		queryParams["next"] = scroll_id
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	response, err := grequests.Get(params["base_url"].(string), &grequests.RequestOptions{
		Params: queryParams,
	})
	if err != nil {
		log.Print(fmt.Sprintf("第%v页查询失败。重试%v次", scroll_id, retryNum), err)
		retryNum++
		time.Sleep(60 * time.Second)
		retrydata := f.retrywanganSearch(params, size, scroll_id, retryNum)
		return retrydata
	}
	body := response.Bytes()
	data := make(map[string]interface{})
	_ = json.Unmarshal(body, &data)
	if errstatus, errok := data["error"]; errok && errstatus.(bool) == true {
		log.Print(fmt.Sprintf("第%v页查询失败。重试%v", scroll_id, retryNum), data)
		log.Print(fmt.Sprintf("第%v页查询失败。重试%v", scroll_id, retryNum), data["error"])
		retryNum++
		time.Sleep(1 * time.Second)
		retrydata := f.retrywanganSearch(params, size, scroll_id, retryNum)
		return retrydata
	}
	if 500 != len(data["results"].([]interface{})) {
		log.Print("results数量不等于500", queryParams)
	}
	return data
}

type wanganserviceFunc func()

func (f wangansearchservice) insertData(data map[string]interface{}, params map[string]interface{}) wanganserviceFunc {
	return func() {
		Subdomain_resolve := make([]map[string]interface{}, 0)
		i := 0
		for _, value := range data["results"].([]interface{}) {
			tmp := map[string]interface{}{}
			tmp["ip_segment"] = value.([]interface{})[0].(string)
			tmp["ip"] = value.([]interface{})[0].(string)
			tmp["ip_type"] = Constant.IP_TYPE_V4
			tmp["createtime"] = time.Now().Format("2006-01-02 15:04:05")
			tmp["belong_user_id"] = params["user_id"].(int)
			tmp["subdomain_clue_id"] = 0
			start_ip, _ := exnet.IPString2Long(value.([]interface{})[0].(string))
			tmp["start_ip"] = int64(start_ip)
			tmp["end_ip"] = int64(start_ip)
			tmp["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS
			tmp["company_id"] = params["company_id"].(int)
			tmp["user_id"] = params["user_id"].(int)
			tmp["port"] = value.([]interface{})[1].(string)
			tmp["protocol"] = value.([]interface{})[2].(string)
			tmp["country"] = value.([]interface{})[3].(string)
			tmp["city"] = value.([]interface{})[4].(string)
			tmp["title"] = value.([]interface{})[5].(string)
			tmp["banner"] = value.([]interface{})[6].(string)
			tmp["body"] = value.([]interface{})[7].(string)
			tmp["server"] = value.([]interface{})[8].(string)
			tmp["domain"] = value.([]interface{})[9].(string)

			tmp["country_name"] = value.([]interface{})[10].(string)
			tmp["region"] = value.([]interface{})[11].(string)
			tmp["longitude"] = value.([]interface{})[12].(string)
			tmp["latitude"] = value.([]interface{})[13].(string)
			tmp["product"] = value.([]interface{})[14].(string)
			tmp["product_category"] = value.([]interface{})[15].(string)
			tmp["header"] = value.([]interface{})[16].(string)
			tmp["fid"] = value.([]interface{})[17].(string)
			tmp["structinfo"] = value.([]interface{})[18].(string)
			tmp["icp"] = value.([]interface{})[19].(string)

			tmp["os"] = value.([]interface{})[20].(string)
			tmp["link"] = value.([]interface{})[21].(string)
			tmp["lastupdatetime"] = value.([]interface{})[22].(string)
			Subdomain_resolve = append(Subdomain_resolve, tmp)
			databases.Redis.Incr("wangansearch:currentsize")
			current_search_ip_list = append(current_search_ip_list, tmp["ip_segment"].(string))
			iparr, _ := Utils.F.Unique(current_search_ip_list)
			current_search_ip_list = iparr.([]string)
			i++
		}
		databases.Redis.Expire("wangansearch:is_finish", 3600*time.Second)
		databases.Redis.Expire("wangansearch:currentsize", 3600*time.Second)
		is_finish_tmp, _ := databases.Redis.Get("wangansearch:is_finish").Result()
		if "" != is_finish_tmp {
			databases.Redis.Set("wangansearch:totalsize", data["size"], 3600*time.Second)
		}
		Models.wanganeesubdomainModel.InsertAsset(Subdomain_resolve)
		Models.wanganeeserviceModel.InsertAsset(Subdomain_resolve)
		//推送进度
		currentTmpsize, _ := databases.Redis.Get("wangansearch:currentsize").Result()
		currentsize, _ := strconv.Atoi(currentTmpsize)
		totalTmpsize, totalTmpErr := databases.Redis.Get("wangansearch:totalsize").Result()
		totalsize, _ := strconv.Atoi(totalTmpsize)
		if 500000 < totalsize {
			return
		}
		if "" != totalTmpsize && 0 != totalsize {
			strprogress, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(currentsize)/float64(totalsize*2))*100), 32)
			progress := strprogress
			log.Print("progressaaa ", progress)
			log.Print("totalTmpsize ", totalTmpsize, totalTmpErr, currentsize)
			f.SetProcess(float32(progress))
		}
		if 0 != totalsize && currentsize == totalsize {
			log.Print(fmt.Sprintf("第%v页查询成功，service和subdomain同步完成 "), totalsize, currentsize)
			f.clearSyncStatus()
		}
	}
}
func (f wangansearchservice) clearSyncStatus() {
	databases.Redis.Del("wangansearch:is_finish")
	databases.Redis.Del("wangansearch:page")
	databases.Redis.Del("wangansearch:currentsize")
	databases.Redis.Del("wangansearch:totalsize")
}

func (f wangansearchservice) clearSyncAssetStatus() {
	databases.Redis.Del("wangansearch:is_finish")
	databases.Redis.Del("wangansearch:page")
	databases.Redis.Del("wangansearch:currentsize")
	databases.Redis.Del("wangansearch:totalsize")
	current_search_ip_list = []string{}
}

func (f wangansearchservice) Sync(params []map[string]interface{}) {
	//Models.ScanAimsModel.InsertScanAims(params)
	//Models.ScanSubdomainsModel.InsertScanSubdomains(params)
	//Models.AssetModel.InsertAsset(Subdomain_resolve)
	//Models.wanganeeserviceModel.InsertAsset(Subdomain_resolve)
	//Models.wanganeesubdomainModel.InsertAsset(Subdomain_resolve)
}
func (f wangansearchservice) SetProcess(progress float32) {
	processdata := make(map[string]interface{})
	listData := make(map[string]interface{})
	listData["progress"] = progress
	listData["type"] = Constant.wanganSEARCH_ASSETS
	processdata["Type"] = "user"
	processdata["Data"] = listData
	processdata["Message"] = "查询成功"
	processdata["Success"] = true
	Websocket.WebConnection.SendToAll(Utils.F.C(processdata))
}

func (f wangansearchservice) SetFailProcess() {
	processdata := make(map[string]interface{})
	listData := make(map[string]interface{})
	listData["type"] = Constant.wanganSEARCH_ASSETS
	processdata["Type"] = "user"
	processdata["Data"] = listData
	processdata["Message"] = "wangan同步失败"
	processdata["Success"] = false
	Websocket.WebConnection.SendToAll(Utils.F.C(processdata))
}
