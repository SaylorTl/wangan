package Jobs

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var LoopholeStatusCheckJob *loopholestatuscheckjob

type loopholestatuscheckjob struct {
	BaseJobs
}

// 参数1：开几个协程
func (L loopholestatuscheckjob) CreateLoopholePool(queryParams map[string]interface{}) {
	var threat_num = Models.ThreatsModel.Count(map[string]interface{}{})
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("EXTRACT_CLUEASSET_FIRL_POOL").Int()
	strInt64 := strconv.FormatInt(threat_num, 10)
	poolNum, _ := strconv.Atoi(strInt64)
	if poolNum > configPoolNum {
		poolNum = configPoolNum
	}
	for i := 0; i < poolNum; i++ {
		err := Utils.P.Submit(L.BatchCheckStatus(queryParams))
		if err != nil {
			return
		}
	}
}

type batchCheckFunc func()

func (L loopholestatuscheckjob) BatchCheckStatus(queryParams map[string]interface{}) batchCheckFunc {
	return func() {
		threat := Models.ThreatsModel.Find(queryParams["threat_id"].(string))
		poc := Models.PocModel.Find(map[string]interface{}{"id": threat["poc_id"]})
		goScannerExecutableCmd := ""
		if Models.LOG4J2_FILENAME == poc.Filename || Models.SPRING_CORE_FILENAME == poc.Filename {
			goScannerExecutableCmd = L.Log4j2(threat)
		} else {
			scanTarget := threat["hostinfo"]
			goScannerExecutableCmd = fmt.Sprintf("%v -m %v -t %v 2>&1", Utils.F.GoScannerPath("goscanner"), poc.Filename, scanTarget)
		}
		log.Print(fmt.Sprintf("loopholeStatusCheck:%v", goScannerExecutableCmd))
		cmd := exec.Command("bash", "-c", goScannerExecutableCmd)
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
		for scanner.Scan() {
			buffer := scanner.Text()
			scanQueue <- buffer
			log.Printf("任务 %d POC %s 扫描过程: %s\n", poc.Filename, buffer)
		}
		if err := cmd.Wait(); err != nil {
			log.Println(err)
		}
	}
}

func (L loopholestatuscheckjob) Log4j2(threat map[string]interface{}) string {
	task_ids := Utils.F.SliceInterfaceToInt(threat["task_ids"].([]interface{}))
	task_id, _ := Utils.F.MaxNum(task_ids)
	task := Models.TaskModel.Find(map[string]interface{}{"id": task_id})
	poc := Models.PocModel.Find(map[string]interface{}{"id": threat["poc_id"]})
	host := strings.Split(threat["hostinfo"].(string), "?")[0]
	newhost := strings.TrimRight(host, "/")
	crawlerBase64 := fmt.Sprintf("crawler_base64_%s#", string(newhost))
	crawlerBase64 = strings.Replace(crawlerBase64, "/", "_", -1)
	files, err := ioutil.ReadDir(task.Asset_scan_complete_filename)
	log.Print(task.Asset_scan_complete_filename)
	log.Print(files)
	if err != nil {
		return ""
	}
	for _, file := range files {
		if false != strings.Hwangxinrefix(file.Name(), crawlerBase64) {
			target := fmt.Sprintf("%v/%v", host, base64.StdEncoding.EncodeToString([]byte(task.Poc_scan_absolute_path+file.Name())))
			return Models.TaskModel.GetGoScannerExecutableCmd(task, poc, target)
		}
	}
	return ""
}
