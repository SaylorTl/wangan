package Utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"wangxin2.0/app/Constant"
)

func InitLog() *os.File {
	gin.DisableConsoleColor()
	daytime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Day())
	logfile := Constant.GoAbsoulePath + "/storage/gin-" + daytime + ".log"
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Print(err)
		file, err = os.Create(logfile)
	}
	if runtime.GOOS == "windows" {
		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
		log.SetOutput(gin.DefaultWriter)
	} else {
		// 如果需要同时将日志写入文件和控制台，请使用以下代码。
		if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
			fmt.Println(err)
			return file
		}
		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
		log.SetOutput(gin.DefaultWriter)
	}
	return file
}

// 时间 7天
const clearLogForwardUnix int64 = 60 * 60 * 24 * 7

// 定时删除过期日志
func ClearLogs() {
	//开启一个12h的定时器
	ticker := time.Tick(time.Second * 12)
	for range ticker {
		//获取当前目录//获取文件或目录相关信息
		wdd := Constant.GoAbsoulePath
		wd := wdd + "/storage"

		fileInfoList, _ := ioutil.ReadDir(wd)

		for i := range fileInfoList {
			// 判断文件是否空的
			if fileInfoList[i].Size() == 64 || (time.Now().Unix()-fileInfoList[i].ModTime().Unix()) > clearLogForwardUnix {
				if !fileInfoList[i].IsDir() {
					delFile(wd + "/" + fileInfoList[i].Name())
				}
				continue
			}
		}
		phpLogDir(Constant.wangxinAbsoulePath+"/storage", "lumen")
	}

}

// 删除文件
func delFile(str string) {
	// 文件夹
	_ = os.RemoveAll(str)
}

func phpLogDir(dir, prefix string) {

	rd, _ := ioutil.ReadDir(dir)
	for _, fi := range rd {
		if fi.IsDir() {
			phpLogDir(dir+string(filepath.Separator)+fi.Name(), prefix)
		} else {
			if strings.Hwangxinrefix(fi.Name(), prefix) {
				if (time.Now().Unix() - fi.ModTime().Unix()) > clearLogForwardUnix {
					toDelFile := dir + string(filepath.Separator) + fi.Name()
					delFile(toDelFile)
					log.Printf("删除文件[%s]\n", toDelFile)
				}
			}
		}
	}
}
