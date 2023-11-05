package Services

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

type scanService struct {
}

var ScanService scanService

func (s scanService) Goscanner(goScannerExecutableCmd string, task *Models.Wx_tasks, poc *Models.Wx_pocs) {
	cmd := exec.Command("bash", "-c", "/wangxin/web/goscanner/goscanner -m vnc_unauth_pwd.json -t http://172.16.20.88:5900 2>&1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("1 %#v", err)
	}
	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
		}
	}(stdout)

	if err := cmd.Start(); err != nil {
		log.Printf("2 %#v", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanQueue := make(chan string, 1000)
	go ScanService.CreatePool(scanQueue, task)
	for scanner.Scan() {
		buffer := scanner.Text()
		scanQueue <- buffer
		log.Printf("任务 %d POC %s 扫描过程: %s\n", task.Id, poc.Filename, buffer)
	}
	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
	//log.Print(fmt.Sprintf("loopholeStatusCheck:%v", goScannerExecutableCmd))
	//findCmd := cmd.NewCmd("bash", "-c", goScannerExecutableCmd)
	//statusChan := findCmd.Start()
	//findCmd.Status().Stdout
	//ticker := time.NewTicker(1 * time.Second)
	//go func() {
	//	for range ticker.C { //推送websocket消息
	//		status := findCmd.Status()
	//		n := len(status.Stdout)
	//		fmt.Println(status.Stdout[n-1])
	//		if status.Stdout[n-1] == "true" {
	//			findCmd.Stop()
	//			ticker.Stop()
	//		}
	//		log.Println(status.Stdout[n-1])
	//	}
	//}()
	//select {
	//case finalStatus := <-statusChan:
	//	log.Print(finalStatus) // done
	//}

}

func (s scanService) CreatePool(scanQueue <-chan string, task *Models.Wx_tasks) {
	configPoolNum := 100
	for j := 0; j < configPoolNum; j++ {
		select {
		case scan := <-scanQueue:
			// 将任务提交到协程池中处理
			_ = Utils.P.Submit(func() {
				saveOrUpdateLoopholeData(scan, task)
			})
		}
	}
}
func saveOrUpdateLoopholeData(resultString string, task *Models.Wx_tasks) {
	splitArray := strings.SplitN(resultString, " ", 3)
	if len(splitArray) < 3 {
		return
	}

	var array map[string]interface{}
	_ = json.Unmarshal([]byte(splitArray[2]), &array)

	if array == nil || array["Vulnerable"] == nil {
		log.Println("saveOrUpdateLoopholeData 没有扫到漏洞数据:" + resultString)
		return
	}

	if !array["Vulnerable"].(bool) && (Constant.LOG4J2_FILENAME == array["FileName"] || Constant.SPRING_CORE_FILENAME == array["FileName"]) {
		//LOG4J2 不支持自动修复漏洞状态
		return
	}
	poc := Models.PocModel.FindPocByName(array["filename"].(string))
	if poc == nil {
		log.Println("saveOrUpdateLoopholeData poc 不存在")
		return
	}

	// 构造搜索条件
	boolQ := elastic.NewBoolQuery()
	boolQ = boolQ.Filter(
		elastic.NewTermQuery("hostinfo", array["HostInfo"]),
		elastic.NewTermQuery("common_title.raw", array["Name"]),
		elastic.NewTermQuery("source", Constant.SOURCE_PLATFORM),
	)
	conf := databases.Conf
	ctx := context.Background()
	searchResult, err := databases.Esdb.Es.Search().
		Index(conf.Section("es").Key("ES_wangxin_LOOPHOLE_INDEX").String()).
		Type(conf.Section("es").Key("ES_wangxin_LOOPHOLE_TYPE").String()).
		Size(1).
		Query(boolQ).
		Pretty(true).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 如果记录不存在且 Vulnerable 字段为 false，直接返回 true
	if !array["Vulnerable"].(bool) && searchResult.TotalHits() == 0 {
		log.Println("saveOrUpdateLoopholeData 漏洞扫描结果为:FALSE;并且之前ES数据库中没有此漏洞；直接end", array)
		return
	}
	ip, port := Utils.F.ParseIpAndPortByHost(array["HostInfo"].(string))
	lastResponse := array["LastResponse"]
	if lastResponse == nil {
		lastResponse = ""
	}
	lastResponse = strings.ReplaceAll(lastResponse.(string), "'", "\\'")
	res, err := databases.Esdb.Es.Get().
		Index(conf.Section("es").Key("ES_wangxin_ASSETS_INDEX").String()).
		Type(conf.Section("es").Key("ES_wangxin_ASSETS_TYPE").String()).
		Id(ip).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("company_id", "geoatla_lft_id", "industry", "geoip")).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	var ipObj map[string]interface{}
	_ = json.Unmarshal(*res.Source, &ipObj)
	companyId := ipObj["company_id"]
	industry := ipObj["industry"]
	geoatlaLftId := ipObj["geoatla_lft_id"]
	var geoip = map[string]interface{}{}
	if ipObj["geoip"] != nil {
		geoip = ipObj["geoip"].(map[string]interface{})
	}
	if !array["Vulnerable"].(bool) && searchResult.TotalHits() > 0 {
		var doc = map[string]interface{}{
			"state":          Constant.STATE_REPAIRED,
			"geoatla_lft_id": geoatlaLftId,
			"company_id":     companyId,
			"industry":       industry,
			"last_response":  lastResponse,
			"repaired_time":  time.Now().Format("2006-01-02 15:04:05"),
		}
		_, err = databases.Esdb.Es.Update().
			Index(conf.Section("es").Key("ES_wangxin_LOOPHOLE_INDEX").String()).
			Type(conf.Section("es").Key("ES_wangxin_LOOPHOLE_TYPE").String()).
			Id(searchResult.Hits.Hits[0].Id).
			Doc(doc).
			Do(ctx)
		return
	}
	log.Println("saveOrUpdateLoopholeData 漏洞扫描结果为:TRUE", array)
	var url string
	if array["VulURL"] != nil {
		url = array["VulURL"].(string)
	} else {
		url = array["HostInfo"].(string)
	}

	data := map[string]interface{}{
		"ip":                 ip,
		"name":               nil,
		"cveId":              nil,
		"port":               port,
		"hostinfo":           array["HostInfo"].(string),
		"last_response":      lastResponse,
		"url":                url,
		"common_title":       array["Name"].(string),
		"common_description": poc.Description,
		"common_impact":      poc.Impact,
		"recommandation":     poc.Recommendation,
		"belong_user_id":     task.User_id,
		"poc_id":             poc.Id,
		"level":              array["Level"].(int),
		"is_ipv6":            Utils.F.GetIpType(ip),
		"vulType":            poc.Vul_type,
		"createtime":         time.Now().Format("2006-01-02 15:04:05"),
		"lastchecktime":      time.Now().Format("2006-01-02 15:04:05"),
		"lastupdatetime":     time.Now().Format("2006-01-02 15:04:05"),
		"noticetime":         nil,
		"repaired_time":      nil,
		"company_id":         companyId,
		"industry":           industry,
		"geoatla_lft_id":     geoatlaLftId,
		"geoip":              geoip,
		"state":              Constant.STATE_UN_REPAIR,
		"source":             Constant.SOURCE_PLATFORM,
		"has_exp":            poc.Has_exp,
		"task_ids":           []int{int(task.Id)},
	}

	if searchResult.TotalHits() > 0 {
		var source map[string]interface{}
		_ = json.Unmarshal(*searchResult.Hits.Hits[0].Source, &source)
		if source["task_ids"] != nil {
			delete(data, "createtime")
			data["task_ids"] = append(data["task_ids"].([]int), source["task_ids"].([]int)...)
		} else {
			data["task_ids"] = append(data["task_ids"].([]int), []int{0}...)
		}

		_, err = databases.Esdb.Es.Update().
			Index(conf.Section("es").Key("ES_wangxin_LOOPHOLE_INDEX").String()).
			Type(conf.Section("es").Key("ES_wangxin_LOOPHOLE_TYPE").String()).
			Id(searchResult.Hits.Hits[0].Id).
			Doc(data).
			Do(ctx)
		return
	}
	data["task_ids"] = append(data["task_ids"].([]int), []int{0}...)

	_, err = databases.Esdb.Es.Index().
		Index(conf.Section("es").Key("ES_wangxin_LOOPHOLE_INDEX").String()).
		Type(conf.Section("es").Key("ES_wangxin_LOOPHOLE_TYPE").String()).
		BodyJson(data).
		Do(ctx)
}
