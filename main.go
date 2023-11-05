package main

import (
	"database/sql"
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"wangxin2.0/app/Commands"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var File *os.File

func init() {
	//日志文件类初始化
	File = Utils.InitLog()
}

func main() {
	// 读取配置文件
	conf, err := ini.Load(Constant.GoAbsoulePath + "/config/my.ini")
	if err != nil {
		log.Fatal(Constant.CONF_READ_ERROR+", err = ", err)
	}
	// 初始化数据库
	db := databases.InitMysql(conf)
	sqlDB, _ := db.DB()
	//初始化工具类
	Utils.InitFunction()
	//初始化model构造函数
	Models.BaseModel{}.Init()

	antsPool := Utils.InitAntPool()

	SubdomainPool := Utils.InitSubdomainPool()

	ImportAssetPool := Utils.InitImportAssetPool()

	ImportSyncAssetPool := Utils.InitImportSyncAssetPool()

	databases.InitEsDb()

	//清楚日志
	go Utils.ClearLogs()
	//初始化参数校验器翻译功能
	transErr := Request.InitTrans("zh")
	if transErr != nil {
		log.Fatal(Constant.INIT_TRANS_ERROR+", err = ", err)
		return
	}

	//初始化redis
	if err := databases.InitRedisDb(conf); err != nil {
		log.Fatal(Constant.INIT_REDIS_ERROR+", err = ", err)
		return
	}
	Commands.Execute()

	defer SubdomainPool.Release()

	defer ImportAssetPool.Release()

	defer antsPool.Release()

	defer ImportSyncAssetPool.Release()
	//关闭文件读取类
	defer func(File *os.File) {
		err := File.Close()
		if err != nil {
			log.Fatal(Constant.CLOSE_FILE_ERROR+", err = ", err)
			return
		}
	}(File)
	//关闭数据库
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Fatal(Constant.CLOSE_DB_ERROR+", err = ", err)
			return
		}
	}(sqlDB)
	//关闭redis
	defer func(Redis *redis.Client) {
		err := Redis.Close()
		if err != nil {
			log.Fatal(Constant.CLOSE_REDIS_ERROR+", err = ", err)
			return
		}
	}(databases.Redis)
}
