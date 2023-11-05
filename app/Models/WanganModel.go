package Models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var wanganeesubdomainModel *wanganeesubdomainmodel

type wanganeesubdomainmodel struct {
	BaseModel
}

type wanganeeSubdomain struct {
	Createtime     LocalTime `json:"createtime" from:"createtime" to:"createtime"`
	Lastupdatetime LocalTime `json:"lastupdatetime" from:"lastupdatetime" to:"lastupdatetime"`
}

type InsertwanganSubdomainAssetnFunc func()

func (a wanganeesubdomainmodel) InsertwanganSubdomainAsset(insertArr []map[string]interface{}) InsertwanganSubdomainAssetnFunc {
	return func() {
		dataArr := make([]map[string]interface{}, 0)
		for _, insertData := range insertArr {
			value := make(map[string]interface{})
			value["_id"] = insertData["_id"].(string)
			val := make(map[string]interface{})
			val = insertData["_source"].(map[string]interface{})
			value["lastupdatetime"] = val["lastupdatetime"].(string)
			value["ip"] = val["ip"]
			if strings.Hwangxinrefix(value["_id"].(string), "https://") {
				value["protocol"] = "https"
			} else {
				value["protocol"] = "http"
			}
			if _, ok := val["header"]; ok {
				value["header"] = val["header"]
			}
			if _, ok := val["os"]; ok {
				value["os"] = val["os"]
			}
			if _, ok := val["gid"]; ok {
				value["gid"] = val["gid"]
			}
			if _, ok := val["dom"]; ok {
				value["dom"] = val["dom"]
			}
			if _, ok := val["fid"]; ok {
				value["fid"] = val["fid"]
			}
			if _, ok := val["host"]; ok {
				value["host"] = val["host"]
			}
			if _, ok := val["tags"]; ok {
				value["rule_tags"] = val["tags"]
			}
			if _, ok := val["is_rabbish"]; ok {
				value["is_rabbish"] = val["is_rabbish"]
			}
			if _, ok := val["is_sensitive"]; ok {
				value["is_sensitive"] = val["is_sensitive"]
			}
			if _, ok := val["language"]; ok {
				value["language"] = val["language"]
			}
			if _, ok := val["js_info"]; ok {
				value["js_info"] = val["js_info"]
			}
			if _, ok := val["jarm"]; ok {
				value["jarm"] = val["jarm"]
			}
			if _, ok := val["is_fraud"]; ok {
				value["is_fraud"] = val["is_fraud"]
			}
			if _, ok := val["is_honeypot"]; ok {
				value["is_honeypot"] = val["is_honeypot"]
			}
			if _, ok := val["appserver"]; ok {
				value["appserver"] = val["appserver"]
			}
			if _, ok := val["body"]; ok {
				value["body"] = val["body"]
			}
			value["is_ipv6"] = false
			if _, ok := val["is_ipv6"]; ok {
				value["is_ipv6"] = val["is_ipv6"]
			}
			value["port"] = val["port"]
			if _, ok := val["server"]; ok {
				value["server"] = val["server"]
			}
			if _, ok := val["title"]; ok {
				value["title"] = val["title"]
			}
			if _, ok := val["subdomain"]; ok {
				value["subdomain"] = val["subdomain"]
			}
			value["isdomain"] = false
			if _, ok := val["fid"]; ok {
				value["fid"] = val["fid"]
			}
			if _, ok := val["products"]; ok {
				value["products"] = val["products"]
			}
			if _, ok := val["product"]; ok {
				value["product"] = val["product"]
			} else {
				value["product"] = []interface{}{}
			}
			// Host
			if _, ok := val["domain"]; ok && "" != val["domain"] {
				value["domain"] = val["domain"]
				value["isdomain"] = true
			}
			// Host
			if _, ok := val["geoip"]; ok {
				value["geoip"] = val["geoip"]
			}
			dataArr = append(dataArr, value)
		}
		a.RetryInsertScanAims(dataArr, 0)
	}
}
func (a wanganeesubdomainmodel) InsertAsset(queryParams []map[string]interface{}) bool {
	dataArr := make([]map[string]interface{}, 0)
	for _, item := range queryParams {
		val := Utils.F.C(item)
		if "http" != val["protocol"] && "https" != val["protocol"] {
			continue
		}
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
		value["ip"] = val["ip"].(string)
		value["header"] = val["header"].(string)
		value["body"] = val["body"].(string)
		value["domain"] = val["domain"].(string)
		value["is_ipv6"] = false
		value["port"] = val["port"].(string)
		value["protocol"] = val["protocol"].(string)
		value["server"] = val["server"].(string)
		value["title"] = val["title"].(string)
		value["isdomain"] = false
		value["fid"] = val["fid"].(string)
		splitproducts := strings.Split(val["product"].(string), ",")
		productArr := make([]interface{}, 0)
		for _, tmpproduct := range splitproducts {
			productArr = append(productArr, tmpproduct)
		}
		value["product"] = productArr

		// Host
		newUrl, _ := url.Parse(val["link"].(string))
		newUrlHost := newUrl.Hostname()
		value["host"] = newUrlHost
		if Utils.F.IsDomain(newUrlHost) {
			value["isdomain"] = true
		}
		value["geoip"] = geoip
		value["_id"] = val["link"].(string)
		dataArr = append(dataArr, value)
	}
	if len(dataArr) <= 0 {
		return true
	}
	a.RetryInsertScanAims(dataArr, 0)
	return true
}

type subdomainReturnFunc func()

func (a wanganeesubdomainmodel) wanganSyncAsset(queryParams []map[string]interface{}) subdomainReturnFunc {
	return func() {
		dataArr := make([]map[string]interface{}, 0)
		for _, val := range queryParams {
			if "http" != val["protocol"] && "https" != val["protocol"] {
				continue
			}
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
			value["ip"] = val["ip"].(string)
			value["header"] = val["header"].(string)
			value["body"] = val["body"].(string)
			value["domain"] = val["domain"].(string)
			value["is_ipv6"] = false
			value["port"] = val["port"].(string)
			value["protocol"] = val["protocol"].(string)
			value["server"] = val["server"].(string)
			value["title"] = val["title"].(string)
			value["isdomain"] = false
			value["fid"] = val["fid"].(string)
			splitproducts := strings.Split(val["product"].(string), ",")
			value["product"] = splitproducts
			// Host
			newUrl, _ := url.Parse(val["link"].(string))
			newUrlHost := newUrl.Hostname()
			value["host"] = newUrlHost
			if Utils.F.IsDomain(newUrlHost) {
				value["isdomain"] = true
			}

			value["geoip"] = geoip
			value["_id"] = val["link"].(string)
			dataArr = append(dataArr, value)
		}
		if len(dataArr) <= 0 {
			return
		}
		a.RetryInsertScanAims(dataArr, 0)
		return
	}
}

func (a wanganeesubdomainmodel) RetryInsertScanAims(dataArr []map[string]interface{}, retrynum int) bool {
	if retrynum >= 10 {
		return false
	}
	conf := databases.Conf
	SCAN_SUBDOMAIN_INDEX := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_INDEX").String()
	SCAN_SUBDOMAIN_TYPE := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_TYPE").String()
	bulkRequest := databases.Esdb.Es.Bulk()
	ctx := context.Background()
	//去重其中重复数据，取更新时间最新的一条
	checkRepeatObj := make(map[string]map[string]interface{})
	for _, repeatVal := range dataArr {
		_id := repeatVal["_id"].(string)
		if _, ok := checkRepeatObj[_id]; ok {
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
	for _, val := range checkRepeatObj {
		_id := val["_id"].(string)
		delete(val, "_id")
		boolQuery := elastic.NewBoolQuery()
		boolQuery.Filter(elastic.NewTermQuery("_id", _id))
		searchResult, err := databases.Esdb.Es.Search().
			Index(SCAN_SUBDOMAIN_INDEX). // 设置索引名
			Type(SCAN_SUBDOMAIN_TYPE).
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
				//log.Print("second check this wanganee_subdomain record is old version")
				continue
			}
			updatereq := elastic.NewBulkUpdateRequest().Index(SCAN_SUBDOMAIN_INDEX).Type(SCAN_SUBDOMAIN_TYPE).Id(_id).Doc(val)
			bulkRequest = bulkRequest.Add(updatereq)
		} else {
			addreq := elastic.NewBulkIndexRequest().Index(SCAN_SUBDOMAIN_INDEX).Type(SCAN_SUBDOMAIN_TYPE).Id(_id).Doc(val)
			bulkRequest = bulkRequest.Add(addreq)
		}
	}
	bulkResponse, subdomainErr := bulkRequest.Refresh("true").Do(ctx)
	if nil != subdomainErr {
		log.Print("subdomainErr", subdomainErr)
		return false
	}
	if len(bulkResponse.Failed()) > 0 {
		filedArr := make([]map[string]interface{}, 0)
		for _, failItem := range bulkResponse.Failed() {
			log.Print("subdomainfailItem.Id", failItem.Id)
			log.Print("subdomainfailItem.Result", failItem.Result)
			log.Print("subdomainfailItem.Error", failItem.Error)
			for _, val := range dataArr {
				if failItem.Id == val["link"].(string) {
					filedArr = append(filedArr, val)
				}
			}
		}
		log.Print(fmt.Sprintf("retry subdomain 第%v次", retrynum))
		retrynum++
		a.RetryInsertScanAims(filedArr, retrynum)
	}
	return true
}

func GetRuleInfo(products []interface{}) []map[string]interface{} {
	rule_infos := make([]map[string]interface{}, 0)
	if len(products) > 0 {
		//如果是subdomain单独判断是否是https或http
		for _, product := range products {
			if "" == product {
				continue
			}
			rule_item := make(map[string]interface{})
			rule_val := RulesModel.FindOne(product.(string))
			rule_item["category"] = rule_val.WxSecondCategoryRules.WxSecondCategories.En_title
			rule_item["cn_category"] = rule_val.WxSecondCategoryRules.WxSecondCategories.Title
			rule_item["company"] = rule_val.En_company
			rule_item["cn_company"] = rule_val.Company
			rule_item["product"] = rule_val.En_product
			rule_item["cn_product"] = rule_val.Product
			rule_item["rule_id"] = rule_val.Id
			rule_item["softhard"] = rule_val.Soft_hard_code
			rule_item["level"] = rule_val.Level_code
			rule_item["cn_parent_category"] = ""
			rule_item["parent_category"] = ""
			if "" != rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry {
				first_cat_id, _ := strconv.Atoi(rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry)
				WxFirstCategories := SecondCategoriesModel.Find(first_cat_id)
				rule_item["cn_parent_category"] = WxFirstCategories.Title
				rule_item["parent_category"] = WxFirstCategories.En_title
			}
			rule_infos = append(rule_infos, rule_item)
		}
	}
	temp_rule_infos, _ := Utils.F.UniqueMap(rule_infos, "rule_id")
	return temp_rule_infos.([]map[string]interface{})
}
