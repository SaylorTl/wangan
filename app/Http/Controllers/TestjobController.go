package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"wangxin2.0/app/Jobs"
)

func Testjob(c *gin.Context) {
	// 需要2个管道
	// 1.job管道
	jobChan := make(chan *Jobs.Job, 128)
	// 2.结果管道
	resultChan := make(chan *Jobs.Result, 128)
	// 3.创建工作池
	Jobs.CreatePool(64, jobChan, resultChan)
	// 4.开个打印的协程
	go func(resultChan chan *Jobs.Result) {
		// 遍历结果管道打印
		for result := range resultChan {
			fmt.Printf("job id:%v randnum:%v result:%d\n", result.Job.Id,
				result.Job.RandNum, result.Sum)
		}
	}(resultChan)
	var id int
	// 循环创建job，输入到管道
	for {
		id++
		// 生成随机数
		r_num := rand.Int()
		job := &Jobs.Job{
			Id:      id,
			RandNum: r_num,
		}
		jobChan <- job
	}

}
