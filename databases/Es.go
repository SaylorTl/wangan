package databases

import (
	"context"
	"errors"
	"github.com/olivere/elastic"
	"log"
)

//var Es *elastic.Client

var Esdb *esdb

type esdb struct {
	Es *elastic.Client
}

func InitEsDb() (client *esdb, err error) {
	// es配置信息
	Esdb = &esdb{}
	esHost := Conf.Section("es").Key("ES_HOST").String()
	esPort := Conf.Section("es").Key("ES_PORT").String()
	log.Print("http://" + esHost + ":" + esPort)
	Esdb.Es, err = elastic.NewClient(elastic.SetURL("http://" + esHost + ":" + esPort))
	if err != nil {
		log.Println("连接失败!!! http://"+esHost+":"+esPort, err)
		return nil, err
	}
	return Esdb, nil
}

// Update 修改文档
// param: index 索引; id 文档ID; m 待更新的字段键值结构
func (e esdb) Update(ctx context.Context, index, doctype, id string, doc interface{}) error {
	_, err := e.Es.
		Update().
		Index(index).
		Type(doctype).
		Id(id).
		Doc(doc).
		Refresh("true").
		Do(ctx)
	return err
}

// Upsert 修改文档（不存在则插入）
// params index 索引; id 文档ID; m 待更新的字段键值结构
func (e esdb) Upsert(ctx context.Context, index, doctype, id string, doc interface{}) error {
	_, err := e.Es.
		Update().
		Index(index).
		Type(doctype).
		Id(id).
		Doc(doc).
		Refresh("true").
		Upsert(doc).
		Do(ctx)
	return err
}

// Upsert 修改文档
// params index 索引; id 文档ID; m 待更新的字段键值结构
func (e esdb) UpdateById(ctx context.Context, index, doctype, id string, doc interface{}) error {
	_, err := e.Es.
		Update().
		Index(index).
		Type(doctype).
		Id(id).
		Doc(doc).
		Refresh("true").
		Do(ctx)
	return err
}

// UpdateByQuery 根据条件修改文档
// param: index 索引; query 条件; script 脚本指定待更新的字段与值
func (e esdb) UpdateByQuery(ctx context.Context, index string, query elastic.Query, script *elastic.Script) (int64, error) {
	rsp, err := e.Es.
		UpdateByQuery(index).
		Query(query).
		Script(script).
		Refresh("true").
		Do(ctx)
	if err != nil {
		return 0, err
	}
	return rsp.Updated, nil
}

// UpdateBulk 批量修改文档
func (e esdb) UpdateBulk(ctx context.Context, index string, ids []string, docs []interface{}) error {
	bulkService := e.Es.Bulk().Index(index).Refresh("true")
	for i := range ids {
		doc := elastic.NewBulkUpdateRequest().
			Id(ids[i]).
			Doc(docs[i])
		bulkService.Add(doc)
	}
	res, err := bulkService.Do(ctx)
	if err != nil {
		return err
	}
	if len(res.Failed()) > 0 {
		return errors.New(res.Failed()[0].Error.Reason)
	}
	return nil
}

// UpsertBulk 批量修改文档（不存在则插入）
func (e esdb) UpsertBulk(ctx context.Context, index string, ids []string, docs []interface{}) error {
	bulkService := e.Es.Bulk().Index(index).Refresh("true")
	for i := range ids {
		doc := elastic.NewBulkUpdateRequest().
			Id(ids[i]).
			Doc(docs[i]).
			Upsert(docs[i])
		bulkService.Add(doc)
	}
	res, err := bulkService.Do(ctx)
	if err != nil {
		return err
	}
	if len(res.Failed()) > 0 {
		return errors.New(res.Failed()[0].Error.Reason)
	}
	return nil
}
