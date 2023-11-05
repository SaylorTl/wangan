package Services

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var cluewg sync.WaitGroup

type clueassetsservice struct {
}

var ClueAssetsService clueassetsservice

func (c clueassetsservice) Find(queryParams map[string]interface{}) bool {
	updateGroupData := make(map[string]interface{})
	AssetData := Models.ClueModel.GetAll(nil, queryParams)
	updateGroupData["id"] = queryParams["recom_group_id"]
	updateGroupData["is_extract_assets_finished"] = 1
	updateGroupData["is_clue_update"] = 0
	if 0 == len(AssetData) {
		Models.RecomGroupModel.UpdateOrCreate(Utils.F.C(updateGroupData))
		return true
	}
	clueAssetChan := make(chan Models.Wx_assets_recommen, len(AssetData))
	for _, recommend := range AssetData {
		clueAssetChan <- recommend
		cluewg.Add(1)
	}
	go c.createFirstPool(updateGroupData, clueAssetChan)
	return true
}

// 创建异步池
func (c clueassetsservice) createFirstPool(updateGroupData map[string]interface{}, clueAssetChan chan Models.Wx_assets_recommen) {
	var poolNum = len(clueAssetChan)
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("EXTRACT_CLUEASSET_FIRL_POOL").Int()
	if poolNum > configPoolNum {
		poolNum = configPoolNum
	}
	for i := 0; i < poolNum; i++ {
		err := Utils.P.Submit(ClueAssetJob.ExtractClueAssets(updateGroupData, clueAssetChan))
		if err != nil {
			return
		}
	}
	cluewg.Wait()
	//关闭管道
	log.Print("数据提取完成")
	_, updateErr := Models.RecomGroupModel.UpdateOrCreate(Utils.F.C(updateGroupData))
	if updateErr != nil {
		log.Print("更新失败")
		return
	}
}

func (c clueassetsservice) getClueAsset(recommend Models.Wx_assets_recommen) {
	params := []interface{}{}
	clueArr := Models.ClueModel.BuildD01Param(recommend.Type, recommend.Keyword)
	params = append(params, clueArr)
	postParams := make(map[string][]interface{})
	postParams["param"] = params
	strParam, _ := json.Marshal(postParams)
	AssetQueryData := D01ICiiApiService.GetAssetEnrichmentData(strParam)
	if _, ok := AssetQueryData["data"]; !ok || reflect.TypeOf(AssetQueryData["data"]).String() != "map[string]interface {}" || AssetQueryData["success"] == false {
		return
	}
	var childWg sync.WaitGroup
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("EXTRACT_CLUEASSET_SECEND_POOL").Int()
	for j := 0; j < configPoolNum; j++ {
		childWg.Add(1)
		err := Utils.P.Submit(ClueAssetJob.ExtractInnerClueAssets(AssetQueryData, recommend, params, &childWg))
		if err != nil {
			return
		}
	}
	childWg.Wait()
}

// 提取任务隔离
var ClueAssetJob extractclueassetjob

type extractclueassetjob struct {
}
type clueassetFunc func()

func (a extractclueassetjob) ExtractClueAssets(updateGroupData map[string]interface{}, clueAssetChan chan Models.Wx_assets_recommen) clueassetFunc {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				//出现异常时，不管提练完与否，都将分组状态置为提炼完成
				_, updateErr := Models.RecomGroupModel.UpdateOrCreate(Utils.F.C(updateGroupData))
				if updateErr != nil {
					return
				}
				return
			}
		}()
		for {
			if len(clueAssetChan) == 0 {
				close(clueAssetChan)
				break
			}
			// 读数据
			AssetP, ok := <-clueAssetChan
			if !ok {
				close(clueAssetChan)
				break
			}
			//从通道取值，直到取完关闭通道
			ClueAssetsService.getClueAsset(AssetP)
			cluewg.Done()
		}
	}
}

func (a extractclueassetjob) ExtractInnerClueAssets(AssetQueryData map[string]interface{}, recommend Models.Wx_assets_recommen, params []interface{}, childWg *sync.WaitGroup) clueassetFunc {
	return func() {
		defer func() {
			childWg.Done()
			if err := recover(); err != nil {
				updateGroupData := make(map[string]interface{})
				updateGroupData["id"] = recommend.Recom_group_id
				updateGroupData["is_extract_assets_finished"] = 1
				updateGroupData["is_clue_update"] = 0
				//出现异常时，不管提练完与否，都将分组状态置为提炼完成
				_, updateErr := Models.RecomGroupModel.UpdateOrCreate(Utils.F.C(updateGroupData))
				log.Println("test Exception CHildClueAsset", err)
				if updateErr != nil {
					return
				}
				return
			}
		}()
		for {
			if _, ok := AssetQueryData["data"]; !ok || reflect.TypeOf(AssetQueryData["data"]).String() != "map[string]interface {}" || AssetQueryData["success"] == false {
				log.Print("循环结束1！")
				return
			}
			AssetQueryMapData := AssetQueryData["data"].(map[string]interface{})
			if _, ok := AssetQueryMapData["data"]; !ok || reflect.TypeOf(AssetQueryMapData["data"]).String() != "map[string]interface {}" {
				log.Print("循环结束2！")
				return
			}
			for ip, ipcontent := range AssetQueryMapData["data"].(map[string]interface{}) {
				insertData := make(map[string]interface{})
				switch recommend.Type {
				case 1:
					insertData["assets_keyword"] = ipcontent.([]interface{})[0].(map[string]interface{})["cert"].(string)
				case 2:
					insertData["assets_keyword"] = ipcontent.([]interface{})[0].(map[string]interface{})["icp"].(string)
				case 3:
					insertData["assets_keyword"] = ipcontent.([]interface{})[0].(map[string]interface{})["domain"].(string)
				}
				insertData["company_id"] = recommend.Cid
				insertData["ip"] = ip
				insertData["recom_group_id"] = recommend.Recom_group_id
				insertData["clue_id"] = recommend.Id
				insertData["user_id"] = recommend.Uid
				_, err := Models.ClueAssetsModel.InsertClueAsset(Utils.F.C(insertData))
				if err != nil {
					continue
				}
			}
			cPostParams := make(map[string]interface{})
			cPostParams["param"] = params
			cPostParams["scroll_id"] = AssetQueryMapData["scroll_id"].(string)
			strParam, _ := json.Marshal(cPostParams)
			AssetQueryData = D01ICiiApiService.GetAssetEnrichmentData(strParam)
		}
	}
}
