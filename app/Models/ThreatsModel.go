package Models

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/databases"
)

var ThreatsModel *threatsmodel

type threatsmodel struct {
	BaseModel
}

func (a threatsmodel) Find(id string) map[string]interface{} {
	conf := databases.Conf
	ctx := context.Background()
	searchResult, err := databases.Esdb.Es.Get().
		Index(conf.Section("es").Key("ES_wangxin_LOOPHOLE_INDEX").String()).
		Type(conf.Section("es").Key("ES_wangxin_LOOPHOLE_TYPE").String()).
		Id(id).       // 设置查询条件
		Pretty(true). // 查询结果返回可读性较好的JSON格式
		Do(ctx)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	temData := make(map[string]interface{})
	umerr := json.Unmarshal(*searchResult.Source, &temData)
	if umerr != nil {
		panic(Constant.JSON_UNENCODE_FAIL)
	}
	return temData
}

func (a threatsmodel) Count(queryParams map[string]interface{}) int64 {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(elastic.NewTermQuery("state", Constant.STATE_UN_REPAIR))
	conf := databases.Conf
	ctx := context.Background()
	searchResult, err := databases.Esdb.Es.Count().
		Index(conf.Section("es").Key("ES_wangxin_LOOPHOLE_INDEX").String()).
		Type(conf.Section("es").Key("ES_wangxin_LOOPHOLE_TYPE").String()).
		Query(boolQuery). // 设置查询条件
		Pretty(true).     // 查询结果返回可读性较好的JSON格式
		Do(ctx)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	log.Print(searchResult)
	//temData := make(map[string]interface{})
	//umerr := json.Unmarshal(*searchResult, &temData)
	//if umerr != nil {
	//	panic(Constant.JSON_UNENCODE_FAIL)
	//}
	return searchResult
}
