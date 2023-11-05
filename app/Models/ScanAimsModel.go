package Models

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/olivere/elastic"
	"log"
	"strconv"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var ScanAimsModel *scanaimsmodel

type scanaimsmodel struct {
	BaseModel
}

func (a scanaimsmodel) Show(queryParams map[string]interface{}) map[string]interface{} {
	conf := databases.Conf
	ctx := context.Background()
	boolQuery := elastic.NewBoolQuery()
	if domain, ok := queryParams["ip_segment"]; ok {
		boolQuery.Filter(elastic.NewTermQuery("ip_segment", domain))
	}
	searchResult, err := databases.Esdb.Es.Scroll().
		Index(conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_INDEX").String()). // 设置索引名
		Type(conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_TYPE").String()).
		Query(boolQuery). // 设置查询条件
		Size(1).
		Pretty(true). // 查询结果返回可读性较好的JSON格式
		Do(ctx)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	returnData := map[string]interface{}{}
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	if 0 == searchResult.Hits.TotalHits {
		return returnData
	}
	for _, val := range searchResult.Hits.Hits {
		temData := make(map[string]interface{})
		err := json.Unmarshal(*val.Source, &temData)
		if err != nil {
			panic(Constant.JSON_UNENCODE_FAIL)
		}
		returnData = temData
	}
	return returnData
}

func (a scanaimsmodel) SyncScanAims(queryParams []map[string]interface{}) bool {
	conf := databases.Conf
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_TYPE").String()
	for _, val := range queryParams {
		value := make(map[string]interface{})
		value["ip_segment"] = val["ip_segment"].(string)
		value["ip_type"] = val["ip_type"].(int)
		value["createtime"] = val["createtime"].(string)
		value["lastupdatetime"] = val["lastupdatetime"].(string)
		value["belong_user_id"] = val["belong_user_id"].(int)
		value["start_ip"] = val["start_ip"].(int64)
		value["end_ip"] = val["end_ip"].(int64)
		value["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS
		value["company_id"] = val["company_id"].(int)
		value["hosts"] = val["hosts"]
		m := md5.Sum([]byte(value["ip_segment"].(string)))
		_id := hex.EncodeToString(m[:])

		scerr := databases.Esdb.Upsert(context.Background(), SCAN_INDEX, SCAN_TYPE, _id, value)
		if scerr != nil {
			return false
		}
	}
	return true
}

func (a scanaimsmodel) InsertScanAims(queryParams []map[string]interface{}) bool {
	for _, val := range queryParams {
		value := make(map[string]interface{})
		value["ip_segment"] = val["ip_segment"].(string)
		value["ip_type"] = val["ip_type"].(int)
		value["createtime"] = val["createtime"].(string)
		value["lastupdatetime"] = val["lastupdatetime"].(string)
		value["belong_user_id"] = val["belong_user_id"].(int)
		value["start_ip"] = val["start_ip"].(int64)
		value["end_ip"] = val["end_ip"].(int64)
		value["ip_source"] = Constant.SOURCE_DOMAIN_DETECT
		value["company_id"] = val["company_id"].(int)
		m := md5.Sum([]byte(value["ip_segment"].(string)))
		_id := hex.EncodeToString(m[:])
		a.RetryInsertScanAims(_id, Utils.F.C(val), value, 0)
	}
	return true
}

func (a scanaimsmodel) RetryInsertScanAims(_id string, queryParams map[string]interface{}, value map[string]interface{}, retrynum int) bool {
	//if retrynum >= 10 {
	//	return false
	//}
	conf := databases.Conf
	ctx := context.Background()
	updateParams := Utils.F.C(value)
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_TYPE").String()
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(elastic.NewTermQuery("_id", _id))
	searchResult, err := databases.Esdb.Es.Scroll().
		Index(SCAN_INDEX). // 设置索引名
		Type(SCAN_TYPE).
		Query(boolQuery). // 设置查询条件
		Pretty(true).     // 查询结果返回可读性较好的JSON格式
		Do(ctx)
	if err != nil && err.Error() != "EOF" {
		log.Print("search", err)
		return false
	}
	if 0 != searchResult.Hits.TotalHits {
		temData := make(map[string]interface{})
		err := json.Unmarshal(*searchResult.Hits.Hits[0].Source, &temData)
		if err != nil {
			log.Print(Constant.JSON_UNENCODE_FAIL)
			return false
		}
		hosts := temData["hosts"].([]interface{})
		if "" != queryParams["domain"].(string) && !Utils.F.InArray(queryParams["domain"].(string), temData["hosts"].([]interface{})) {
			hosts = append(temData["hosts"].([]interface{}), queryParams["domain"].(string))
		}
		unique, err := Utils.F.Unique(hosts)
		if err != nil {
			return false
		}
		updateParams["hosts"] = unique
	} else {
		if _, ok := queryParams["domain"]; ok && queryParams["domain"] != "" {
			updateParams["hosts"] = append([]string{}, queryParams["domain"].(string))
		} else {
			updateParams["hosts"] = []string{}
		}

	}

	scerr := databases.Esdb.Upsert(context.Background(), SCAN_INDEX, SCAN_TYPE, _id, updateParams)
	if scerr != nil {
		retrynum++
		log.Print("Upsert_scan_aims_index_retry"+strconv.Itoa(retrynum)+":", scerr)
		time.Sleep(5 * time.Second)
		a.RetryInsertScanAims(_id, queryParams, value, retrynum)
	}
	return true
}
