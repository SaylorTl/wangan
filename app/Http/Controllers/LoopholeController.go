package Controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"log"
	"os/exec"
	"strings"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Jobs"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
)

type LoopholeController struct {
	BaseController
}

func (a LoopholeController) Check(c *gin.Context) {
	var params Request.LoopholeCheckValidate
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	//go Jobs.LoopholeStatusCheckJob.CreateLoopholePool(queryParams)

	threat := Models.ThreatsModel.Find(queryParams["threat_id"].(string))
	poc := Models.PocModel.Find(map[string]interface{}{"id": threat["poc_id"]})
	goScannerExecutableCmd := ""
	if Models.LOG4J2_FILENAME == poc.Filename || Models.SPRING_CORE_FILENAME == poc.Filename {
		goScannerExecutableCmd = Jobs.LoopholeStatusCheckJob.Log4j2(threat)
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
		log.Printf("任务 %v POC 扫描过程: %v\n", poc.Filename, buffer)
		arr := strings.Split(buffer, " ")
		//arr := strings.Fields(buffer)
		log.Printf("漏洞字符串切割：%q\n", arr)
		if arr[2] == "" {
			continue
		}
		var loopholeArray map[string]interface{}
		_ = json.Unmarshal([]byte(arr[2]), loopholeArray)
		log.Printf("漏洞扫描测试", loopholeArray)
		if loopholeArray == nil {
			continue
		}
		state := Constant.STATE_UN_REPAIR
		if loopholeArray["Vulnerable"] == true {
			state = Constant.STATE_REPAIRED
		}
		log.Print(state)
	}
	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
	Constant.Success(c, "msg", "")
	return
}
