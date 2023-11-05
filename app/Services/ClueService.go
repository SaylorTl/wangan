package Services

import (
	"encoding/json"
	"log"
	"net"
	"reflect"
	"strconv"
	"sync"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/app/Utils/Websocket"
	"wangxin2.0/databases"
)

var ClueService clueservice

var wg sync.WaitGroup

type clueservice struct {
}

func (a clueservice) Extract(queryParams map[string]interface{}) (bool, string) {
	//构建异步队列，防并发场景
	websocketData := make(map[string]interface{})
	databases.Redis.Set("ExtractClue:"+strconv.Itoa(queryParams["asset_group_id"].(int)), 1, 0)
	websocketData["group_id"] = queryParams["asset_group_id"].(int)
	a.SendToWebsocket(Utils.F.C(websocketData), queryParams["user_id"].(int))
	updateData := make(map[string]interface{})
	if queryParams["is_add_clue_asset"].(int) == 1 {
		updateData["is_extract_assets_finished"] = 0
	}
	if queryParams["is_add_clue_asset"].(int) == 0 {
		updateData["is_clue_update"] = 1
	}
	updateData["is_extract_clue_finished"] = 0

	rangeData := queryParams["ip"]
	if nil == queryParams["ip"] {
		esData := Models.AssetModel.AssetExtractClue(Utils.F.C(queryParams))
		rangeData = esData
	} else {
		rangeData = queryParams["ip"]
		for _, vv := range rangeData.([]string) {
			if Utils.F.Isipv6(vv) {
				return false, Constant.IPV6_CANNOT_EXTRACT
			}
		}
	}
	if len(rangeData.([]string)) == 0 {
		return false, Constant.GROUP_CANNOT_EMPTY
	}
	AssetChan := make(chan map[string]interface{}, len(rangeData.([]string)))
	for _, inputIp := range rangeData.([]string) {
		//未知资产发异步发现
		AssetParams := make(map[string]interface{})
		AssetParams["inputIp"] = inputIp
		AssetParams["user_id"] = queryParams["user_id"].(int)
		AssetParams["asset_group_id"] = queryParams["asset_group_id"].(int)
		AssetParams["is_add_clue_asset"] = queryParams["is_add_clue_asset"].(int)
		AssetChan <- AssetParams
		wg.Add(1)
	}
	updateData["id"] = queryParams["asset_group_id"].(int)
	go a.createCluePool(updateData, queryParams, AssetChan)
	return true, Constant.EXTRACT_SUCCESS
}

// 创建提炼池
func (a clueservice) createCluePool(updateData map[string]interface{}, queryParams map[string]interface{}, AssetChan chan map[string]interface{}) {
	emptyData := make(map[string]interface{})
	emptyData["group_id"] = queryParams["asset_group_id"].(int)
	chanLen := len(AssetChan)
	var poolNum = chanLen
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("EXTRACT_CLUE_FIRL_POOL").Int()
	if chanLen > configPoolNum {
		poolNum = configPoolNum
	}
	for i := 0; i < poolNum; i++ {
		err := Utils.P.Submit(ClueJob.ExtractClue(updateData, queryParams, AssetChan, &wg))
		if err != nil {
			return
		}
	}
	//关闭管道
	wg.Wait()
	log.Print("test 执行完成")
	//更新分组线索数量
	emptyData["progress"] = "100%"
	updateData["is_extract_clue_finished"] = 1
	_, err := Models.RecomGroupModel.UpdateOrCreate(Utils.F.C(updateData))
	if err != nil {
		log.Print(err)
		return
	}
	//将分组状态，置为提炼完成
	if queryParams["is_add_clue_asset"].(int) == 1 {
		//提炼完线索之后提炼资产
		extractParams := make(map[string]interface{})
		extractParams["recom_group_id"] = queryParams["asset_group_id"].(int)
		ClueAssetsService.Find(Utils.F.C(extractParams))
	}
	a.SendToWebsocket(Utils.F.C(emptyData), queryParams["user_id"].(int))
	RecommGroupService.UpdateClueCount(queryParams["asset_group_id"].(int))
}

// 提炼线索
func (c clueservice) Rich(AssetP map[string]interface{}) (reponseData any) {
	addr, err := net.ResolveIPAddr("ip", AssetP["inputIp"].(string))
	if err != nil {
		panic(err.Error())
	}
	parseIp := addr.String()
	params := make(map[string][]map[string]string)
	arr := make([]map[string]string, 1)
	subMapC := make(map[string]string)
	subMapC["type"] = strconv.Itoa(Constant.AssetEnrichment_TYPE_IP)
	subMapC["content"] = parseIp
	arr[0] = subMapC
	params["param"] = arr
	strParam, _ := json.Marshal(params)
	responseArray := D01ICiiApiService.GetAssetEnrichmentData(strParam)
	if _, ok := responseArray["data"]; !ok || reflect.TypeOf(responseArray["data"]).String() != "map[string]interface {}" || responseArray["success"] == false {
		return
	}
	// 创建提炼池
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("EXTRACT_CLUE_SECEND_POOL").Int()
	var childWg sync.WaitGroup
	for i := 0; i < configPoolNum; i++ {
		childWg.Add(1)
		err := Utils.P.Submit(ClueJob.ExtractInnerClue(responseArray, AssetP, parseIp, &childWg))
		if err != nil {
			return
		}
	}
	childWg.Wait()
	return parseIp
}

// 发送websocket信息
func (a clueservice) SendToWebsocket(listData map[string]interface{}, user_id int) {
	if _, ok := listData["progress"]; ok {
		databases.Redis.Del("ExtractClue:" + strconv.Itoa(listData["group_id"].(int)))
	} else {
		currentProcess := databases.Redis.Get("ExtractClue:" + strconv.Itoa(listData["group_id"].(int)))
		curProcess, _ := currentProcess.Int()
		if curProcess > 100 {
			curProcess = 99
		}
		listData["progress"] = strconv.Itoa(curProcess) + "%"
	}
	mapdata := make(map[string]interface{})
	listData["type"] = Constant.EXTRACT_CLUE_PROCESS
	mapdata["success"] = true
	mapdata["Data"] = listData
	mapdata["Message"] = "查询成功"
	go Websocket.WebConnection.SendToUser(Utils.F.C(mapdata), user_id)
	go Websocket.WebConnection.SendToUser(Utils.F.C(mapdata), 1)
}

// 查询数据库线索信息
func (a clueservice) Find(queryParams map[string]interface{}, pagination Models.Pagination) (res []Models.Wx_assets_recommen, countRes int64) {
	// 这里可以组装查询条件
	recommenRes, countRes := Models.ClueModel.Lists(queryParams, &pagination)
	return recommenRes, countRes
}

// 线索入库
func (a clueservice) ExtractClues(arr []interface{}, AssetP map[string]interface{}, stype int, parseIp string) {
	extractArr, _ := Utils.F.Unique(arr)
	ipRes := Models.IpsModel.Findip(parseIp)
	for _, value := range extractArr.([]interface{}) {
		if value == nil || value == "" {
			continue
		}
		user_id, _ := AssetP["user_id"].(int)
		asset_group_id, _ := AssetP["asset_group_id"].(int)
		insertData := make(map[string]interface{})
		insertData["uid"] = user_id
		insertData["cid"] = ipRes.Company_id
		insertData["type"] = stype
		insertData["assets_id"] = ipRes.Id
		insertData["status"] = Constant.STATUS_INIT
		insertData["keyword"] = value.(string)
		insertData["assets_keyword"] = parseIp
		insertData["recom_group_id"] = asset_group_id
		err := Models.ClueModel.InsertClue(Utils.F.C(insertData))
		if err != nil {
			log.Print("test assetEnrichment:", err)
			panic(err)
		}
	}
}

// 提取任务隔离
var ClueJob extractcluejob

type extractcluejob struct {
}
type taskFunc func()

// 创建提炼池
func (a extractcluejob) ExtractClue(updateData map[string]interface{}, queryParams map[string]interface{}, AssetChan chan map[string]interface{}, wg *sync.WaitGroup) taskFunc {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				//出现异常时，不管提练完与否，都将分组状态置为提炼完成
				databases.Redis.Del("ExtractClue:" + strconv.Itoa(queryParams["asset_group_id"].(int)))
				updateData["id"] = queryParams["asset_group_id"].(int)
				updateData["is_extract_clue_finished"] = 1
				updateData["is_extract_assets_finished"] = 1
				update_id, _ := Models.RecomGroupModel.UpdateOrCreate(updateData)
				log.Println("test Exception clue"+strconv.Itoa(update_id), err)
				return
			}
		}()
		for {
			if len(AssetChan) == 0 {
				close(AssetChan)
				break
			}
			// 写数据
			AssetP, ok := <-AssetChan
			if !ok {
				close(AssetChan)
				break
			}
			//将分组状态置为未提炼完
			updateData["id"] = queryParams["asset_group_id"].(int)
			_, uerr := Models.RecomGroupModel.UpdateOrCreate(updateData)
			if uerr != nil {
				return
			}
			log.Print("test 提炼开始", AssetP)
			ClueService.Rich(AssetP)
			wg.Done()
			log.Print("test 提炼结束", len(AssetChan))
		}
	}
}

// 创建提炼池
func (a extractcluejob) ExtractInnerClue(responseArray map[string]interface{}, AssetP map[string]interface{}, parseIp string, childWg *sync.WaitGroup) taskFunc {
	return func() {
		user_id, _ := AssetP["user_id"].(int)
		listData := make(map[string]interface{})
		listData["group_id"] = AssetP["asset_group_id"]
		defer func() {
			childWg.Done()
			if err := recover(); err != nil {
				//出现异常时，不管提练完与否，都将分组状态置为提炼完成
				log.Println("test Exception CHildClue"+parseIp, err)
				return
			}
		}()
		for {
			if _, ok := responseArray["data"]; !ok || reflect.TypeOf(responseArray["data"]).String() != "map[string]interface {}" || responseArray["success"] == false {
				return
			}
			data := responseArray["data"].(map[string]interface{})
			if _, ok := data["data"]; !ok || reflect.TypeOf(data["data"]).String() != "map[string]interface {}" {
				return
			}
			for _, enrichValue := range data["data"].(map[string]interface{}) {
				if reflect.TypeOf(enrichValue).String() != "[]interface {}" {
					return
				}
				domainArray := Utils.F.MapColumn(enrichValue.([]interface{}), "domain")
				certArray := Utils.F.MapColumn(enrichValue.([]interface{}), "cert")
				icpArray := Utils.F.MapColumn(enrichValue.([]interface{}), "icp")
				ClueService.ExtractClues(domainArray, AssetP, Constant.TYPE_DOMAIN, parseIp)
				ClueService.ExtractClues(certArray, AssetP, Constant.TYPE_CERT, parseIp)
				ClueService.ExtractClues(icpArray, AssetP, Constant.TYPE_ICP, parseIp)
			}
			scrollparams := make(map[string]interface{})
			scrollarr := make([]map[string]string, 1)
			scrollsubMapC := make(map[string]string)
			scrollsubMapC["type"] = strconv.Itoa(Constant.AssetEnrichment_TYPE_IP)
			scrollsubMapC["content"] = parseIp
			scrollarr[0] = scrollsubMapC
			scrollparams["param"] = scrollarr
			scrollparams["scroll_id"] = data["scroll_id"].(string)
			scrollstrParam, _ := json.Marshal(scrollparams)
			//推送至websocker
			queryParams := make(map[string]interface{})
			queryParams["Assets_keyword"] = parseIp
			queryParams["Status"] = Constant.STATUS_INIT
			currentProcess, _ := databases.Redis.Get("ExtractClue:" + strconv.Itoa(AssetP["asset_group_id"].(int))).Result()
			curProcess, _ := strconv.Atoi(currentProcess)
			process := curProcess + 1
			databases.Redis.SetXX("ExtractClue:"+strconv.Itoa(AssetP["asset_group_id"].(int)), process, 0)
			responseArray = D01ICiiApiService.GetAssetEnrichmentData(scrollstrParam)
			//是map时继续执行
			ClueService.SendToWebsocket(Utils.F.C(listData), user_id)
		}
	}
}
