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

var ScanSubdomainsModel *scansubdomainsmodel

type scansubdomainsmodel struct {
	BaseModel
}

func (a scansubdomainsmodel) Show(queryParams map[string]interface{}) map[string]interface{} {
	conf := databases.Conf
	ctx := context.Background()
	boolQuery := elastic.NewBoolQuery()
	if domain, ok := queryParams["domain"]; ok {
		boolQuery.Filter(elastic.NewTermQuery("domain", domain))
	}
	searchResult, err := databases.Esdb.Es.Scroll().
		Index(conf.Section("es").Key("ES_wangxin_WX_SCAN_SUBDOMAINS_INDEX").String()). // 设置索引名
		Type(conf.Section("es").Key("ES_wangxin_WX_SCAN_SUBDOMAINS_TYPE").String()).
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
func (a scansubdomainsmodel) SyncScanSubdomains(queryParams []map[string]interface{}) bool {
	conf := databases.Conf
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_WX_SCAN_SUBDOMAINS_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_WX_SCAN_SUBDOMAINS_TYPE").String()
	for _, val := range queryParams {
		if val["domain"].(string) == "" {
			continue
		}
		value := make(map[string]interface{})
		value["domain"] = val["domain"].(string)
		value["createtime"] = val["createtime"].(string)
		value["lastupdatetime"] = val["lastupdatetime"].(string)
		value["belong_user_id"] = val["belong_user_id"].(int)
		value["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS
		value["company_id"] = val["company_id"].(int)
		value["ip"] = val["ip"].(string)
		m := md5.Sum([]byte(value["domain"].(string) + val["ip"].(string)))
		_id := hex.EncodeToString(m[:])
		scerr := databases.Esdb.Upsert(context.Background(), SCAN_INDEX, SCAN_TYPE, _id, value)
		if scerr != nil {
			return false
		}
	}
	return true
}

func (a scansubdomainsmodel) InsertScanSubdomains(queryParams []map[string]interface{}) bool {
	for _, val := range queryParams {
		if val["domain"].(string) == "" {
			continue
		}
		value := make(map[string]interface{})
		value["domain"] = val["domain"].(string)
		value["createtime"] = val["createtime"].(string)
		value["lastupdatetime"] = val["lastupdatetime"].(string)
		value["belong_user_id"] = val["belong_user_id"].(int)
		value["ip_source"] = Constant.SOURCE_DOMAIN_DETECT
		value["company_id"] = val["company_id"].(int)
		value["ip"] = val["ip"].(string)
		m := md5.Sum([]byte(value["domain"].(string) + val["ip"].(string)))
		_id := hex.EncodeToString(m[:])
		a.RetryInsertScanAims(_id, Utils.F.C(val), value, 0)
	}
	return true
}

func (a scansubdomainsmodel) RetryInsertScanAims(_id string, queryParams map[string]interface{}, value map[string]interface{}, retrynum int) bool {
	//if retrynum >= 10 {
	//	return false
	//}
	conf := databases.Conf
	ctx := context.Background()
	updateParams := Utils.F.C(value)
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_WX_SCAN_SUBDOMAINS_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_WX_SCAN_SUBDOMAINS_TYPE").String()
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
		subdomain_clue_ids := []int{queryParams["subdomain_clue_id"].(int)}
		switch temData["subdomain_clue_ids"].(type) {
		case []interface{}:
			for _, tmpid := range temData["subdomain_clue_ids"].([]interface{}) {
				subdomain_clue_ids = append(subdomain_clue_ids, int(tmpid.(float64)))
			}
		}
		sunique, serr := Utils.F.Unique(subdomain_clue_ids)
		if serr != nil {
			return false
		}
		updateParams["subdomain_clue_ids"] = sunique
	} else {
		updateParams["subdomain_clue_ids"] = append([]int{}, queryParams["subdomain_clue_id"].(int))
	}

	scerr := databases.Esdb.Upsert(context.Background(), SCAN_INDEX, SCAN_TYPE, _id, updateParams)
	if scerr != nil {
		time.Sleep(5 * time.Second)
		retrynum++
		log.Print("Upsert_scan_subdomain_retry"+strconv.Itoa(retrynum)+":", scerr)
		a.RetryInsertScanAims(_id, queryParams, value, retrynum)
	}
	return true
}
