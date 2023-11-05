package Models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"strings"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var wanganeeserviceModel *wanganeeservicemodel

type wanganeeservicemodel struct {
	BaseModel
}

type wanganeeservice struct {
	Createtime     LocalTime `json:"createtime" from:"createtime" to:"createtime"`
	Lastupdatetime LocalTime `json:"lastupdatetime" from:"lastupdatetime" to:"lastupdatetime"`
}

type InsertwanganServicenAssetnFunc func()

func (a wanganeeservicemodel) InsertwanganServiceAsset(insertArr []map[string]interface{}) InsertwanganServicenAssetnFunc {
	return func() {
		dataArr := make([]map[string]interface{}, 0)
		for _, insertData := range insertArr {
			value := make(map[string]interface{})
			val := make(map[string]interface{})
			val = insertData["_source"].(map[string]interface{})
			value["_id"] = insertData["_id"].(string)
			if _, ok := val["base_protocol"]; ok {
				value["base_protocol"] = val["base_protocol"]
			}
			if _, ok := val["cert"]; ok {
				value["cert"] = val["cert"]
			}
			if _, ok := val["certs"]; ok {
				value["certs"] = val["certs"]
			}
			if _, ok := val["is_honeypot"]; ok {
				value["is_honeypot"] = val["is_honeypot"]
			}
			if _, ok := val["is_sensitive"]; ok {
				value["is_sensitive"] = val["is_sensitive"]
			}
			if _, ok := val["os"]; ok {
				value["os"] = val["os"]
			}
			if _, ok := val["banner"]; ok {
				value["banner"] = val["banner"]
			}
			if _, ok := val["geoip"]; ok {
				value["geoip"] = val["geoip"]
			}
			if _, ok := val["lastchecktime"]; ok {
				value["lastchecktime"] = val["lastchecktime"]
			}
			if _, ok := val["lastupdatetime"]; ok {
				value["lastupdatetime"] = val["lastupdatetime"]
			}
			value["is_ipv6"] = false
			if _, ok := val["is_ipv6"]; ok {
				value["is_ipv6"] = val["is_ipv6"]
			}
			value["ip"] = val["ip"].(string)
			if _, ok := val["domain"]; ok {
				value["domain"] = val["domain"]
			}
			value["port"] = val["port"]
			if _, ok := val["protocol"]; ok {
				value["protocol"] = val["protocol"]
			}
			if _, ok := val["product"]; ok {
				value["product"] = val["product"]
			} else {
				value["product"] = []interface{}{}
			}
			if _, ok := val["tags"]; ok {
				value["rule_tags"] = val["tags"]
			}
			if _, ok := val["server"]; ok {
				value["server"] = val["server"]
			}
			if _, ok := val["time"]; ok {
				value["time"] = val["time"]
			}
			dataArr = append(dataArr, value)
		}
		a.RetryInsertScanAims(dataArr, 0)
	}
}

func (a wanganeeservicemodel) InsertAsset(queryParams []map[string]interface{}) bool {
	dataArr := make([]map[string]interface{}, 0)
	for _, item := range queryParams {
		val := Utils.F.C(item)
		geoip := make(map[string]interface{})
		geoip["longitude"] = val["longitude"].(string)
		geoip["latitude"] = val["latitude"].(string)
		geoip["country_name"] = val["country_name"].(string)
		geoip["city_name"] = strings.ToLower(val["city"].(string))
		location := make(map[string]interface{})
		location["lon"] = val["longitude"].(string)
		location["lat"] = val["latitude"].(string)
		geoip["location"] = location

		value := make(map[string]interface{})
		value["lastchecktime"] = val["createtime"].(string)
		value["lastupdatetime"] = val["lastupdatetime"].(string)
		value["is_ipv6"] = false
		value["banner"] = val["banner"].(string)
		value["base_protocol"] = "tcp"
		value["ip"] = val["ip"].(string)
		value["geoip"] = geoip
		value["port"] = val["port"].(string)
		value["protocol"] = val["protocol"].(string)
		splitproducts := strings.Split(val["product"].(string), ",")
		productArr := make([]interface{}, 0)
		for _, tmpproduct := range splitproducts {
			productArr = append(productArr, tmpproduct)
		}
		value["product"] = productArr
		value["server"] = val["server"].(string)
		value["time"] = time.Now()
		value["_id"] = val["ip"].(string) + ":" + val["port"].(string)
		dataArr = append(dataArr, value)
	}
	if len(dataArr) <= 0 {
		return true
	}
	a.RetryInsertScanAims(dataArr, 0)

	return true
}

type serviceReturnFunc func()

func (a wanganeeservicemodel) wanganSyncAsset(queryParams []map[string]interface{}) serviceReturnFunc {
	return func() {
		dataArr := make([]map[string]interface{}, 0)
		for _, val := range queryParams {
			geoip := make(map[string]interface{})
			geoip["longitude"] = val["longitude"].(string)
			geoip["latitude"] = val["latitude"].(string)
			geoip["country_name"] = val["country_name"].(string)
			geoip["city_name"] = strings.ToLower(val["city"].(string))
			location := make(map[string]interface{})
			location["lon"] = val["longitude"].(string)
			location["lat"] = val["latitude"].(string)
			geoip["location"] = location

			value := make(map[string]interface{})
			value["lastchecktime"] = val["createtime"].(string)
			value["lastupdatetime"] = val["lastupdatetime"].(string)
			value["is_ipv6"] = false
			value["banner"] = val["banner"].(string)
			value["base_protocol"] = "tcp"
			value["ip"] = val["ip"].(string)
			value["geoip"] = geoip
			value["port"] = val["port"].(string)
			value["protocol"] = val["protocol"].(string)
			splitproducts := strings.Split(val["product"].(string), ",")
			value["product"] = splitproducts
			value["server"] = val["server"].(string)
			value["time"] = time.Now()
			value["_id"] = val["ip"].(string) + ":" + val["port"].(string)
			dataArr = append(dataArr, value)
		}
		if len(dataArr) <= 0 {
			return
		}
		a.RetryInsertScanAims(dataArr, 0)
		return
	}
}

func (a wanganeeservicemodel) RetryInsertScanAims(dataArr []map[string]interface{}, retrynum int) bool {
	if retrynum >= 10 {
		return false
	}
	conf := databases.Conf
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_SERVICE_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_SERVICE_TYPE").String()
	bulkRequest := databases.Esdb.Es.Bulk()
	ctx := context.Background()
	//去重其中重复数据，取更新时间最新的一条
	checkRepeatObj := make(map[string]map[string]interface{})
	for _, repeatVal := range dataArr {
		_id := repeatVal["_id"].(string)
		if _, ok := checkRepeatObj[_id]; ok {
			log.Print("checkrepeat", _id, " ", checkRepeatObj[_id]["lastupdatetime"])
			updatetime, updateerr := time.Parse("2006-01-02 15:04:05", repeatVal["lastupdatetime"].(string))
			if nil != updateerr {
				continue
			}
			existtime, existerr := time.Parse("2006-01-02 15:04:05", checkRepeatObj[_id]["lastupdatetime"].(string))
			if nil != existerr {
				continue
			}
			if updatetime.Before(existtime) {
				//log.Print("first check this wanganee_subdomain record is old version")
				continue
			}
		}
		checkRepeatObj[_id] = repeatVal
	}
	//入库
	for _, val := range checkRepeatObj {
		_id := val["_id"].(string)
		delete(val, "_id")
		boolQuery := elastic.NewBoolQuery()
		boolQuery.Filter(elastic.NewTermQuery("_id", _id))
		searchResult, err := databases.Esdb.Es.Search().
			Index(SCAN_INDEX). // 设置索引名
			Type(SCAN_TYPE).
			FetchSourceContext(elastic.NewFetchSourceContext(true).Include("lastupdatetime")).
			Query(boolQuery). // 设置查询条件
			Pretty(true).     // 查询结果返回可读性较好的JSON格式
			Do(ctx)
		if err != nil && err.Error() != "EOF" {
			log.Print("search", err)
			return false
		}
		val["rule_tags"] = []interface{}{}
		if _, ok := val["product"]; ok && len(val["product"].([]interface{})) > 0 {
			for i := 0; i < len(val["product"].([]interface{})); i++ {
				rule_val := RulesModel.FindOne(val["product"].([]interface{})[i].(string))
				if "" == rule_val.Product {
					val["product"] = append(val["product"].([]interface{})[:i], val["product"].([]interface{})[i+1:]...)
					i--
				}
			}
			val["rule_tags"] = GetRuleInfo(val["product"].([]interface{}))
		}
		if 0 != searchResult.Hits.TotalHits {
			temData := make(map[string]interface{})
			err := json.Unmarshal(*searchResult.Hits.Hits[0].Source, &temData)
			if err != nil {
				log.Print(Constant.JSON_UNENCODE_FAIL)
				return false
			}
			updatetime, err := time.Parse("2006-01-02 15:04:05", val["lastupdatetime"].(string))
			existtime, err := time.Parse("2006-01-02 15:04:05", temData["lastupdatetime"].(string))
			if err == nil && updatetime.Before(existtime) {
				//log.Print("second check  this wanganee_service record is old version")
				continue
			}
			updatereq := elastic.NewBulkUpdateRequest().Index(SCAN_INDEX).Type(SCAN_TYPE).Id(_id).Doc(val)
			bulkRequest = bulkRequest.Add(updatereq)
		} else {
			addreq := elastic.NewBulkIndexRequest().Index(SCAN_INDEX).Type(SCAN_TYPE).Id(_id).Doc(val)
			bulkRequest = bulkRequest.Add(addreq)
		}
	}
	bulkResponse, serviceErr := bulkRequest.Refresh("true").Do(ctx)
	if nil != serviceErr {
		log.Print("serviceErr", serviceErr)
		return false
	}
	if len(bulkResponse.Failed()) > 0 {
		filedArr := make([]map[string]interface{}, 0)
		for _, failItem := range bulkResponse.Failed() {
			log.Print("servicefailItem.Id", failItem.Id)
			log.Print("servicefailItem.Result", failItem.Result)
			log.Print("servicefailItem.Error", failItem.Error)
			for _, val := range dataArr {
				port_str := Utils.F.Ts(val["port"])
				if failItem.Id == val["ip"].(string)+":"+port_str {
					filedArr = append(filedArr, val)
				}
			}
		}
		retrynum++
		log.Print(fmt.Sprintf("retry service 第%v次", retrynum))
		a.RetryInsertScanAims(filedArr, retrynum)
	}
	return true
}
