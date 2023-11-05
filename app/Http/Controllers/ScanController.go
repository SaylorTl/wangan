package Controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
)

type ScanController struct {
	BaseController
}

func (s ScanController) GoScanner(c *gin.Context) {
	param := c.PostForm("task_id")
	taskId, _ := strconv.Atoi(param)
	param = c.PostForm("poc_id")
	pocId, _ := strconv.Atoi(param)
	task := Models.TaskModel.FindTaskById(taskId)
	poc := Models.PocModel.FindPocById(pocId)
	goScannerExecutableCmd := Models.TaskModel.GetGoScannerExecutableCmd(*task, *poc, "")
	log.Printf("%s", goScannerExecutableCmd)
	if task.Step != Constant.STEP_LOOPHOLE {
		Models.TaskModel.UpdateColumns(task, map[string]interface{}{"step": Constant.STEP_LOOPHOLE})
	}
	//go Services.ScanService.Goscanner(goScannerExecutableCmd, task, poc)
	Constant.Success(c, "操作成功!", "")
}
