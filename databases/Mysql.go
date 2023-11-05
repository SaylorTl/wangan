package databases

import (
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

var Db *gorm.DB
var Conf *ini.File

func InitMysql(conf *ini.File) *gorm.DB {
	var connErr error
	Conf = conf
	// 数据库的类型
	dbType := conf.Section("").Key("DB_TYPE").String()
	// Mysql配置信息
	mysqlName := conf.Section("mysql").Key("DB_NAME").String()
	mysqlUser := conf.Section("mysql").Key("DB_USER").String()
	mysqlPwd := conf.Section("mysql").Key("DB_PWD").String()
	mysqlHost := conf.Section("mysql").Key("DB_HOST").String()
	mysqlPort := conf.Section("mysql").Key("DB_PORT").String()
	mysqlCharset := conf.Section("mysql").Key("DB_CHATSET").String()
	// sqlite3配置信息
	var dataSource string
	switch dbType {
	case "mysql":
		dataSource = mysqlUser + ":" + mysqlPwd + "@tcp(" + mysqlHost + ":" +
			mysqlPort + ")/" + mysqlName + "?parseTime=true&loc=Local&charset=" + mysqlCharset
		Db, connErr = gorm.Open(mysql.Open(dataSource),
			&gorm.Config{Logger: logger.Default.LogMode(logger.Info), NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			}})
	}
	if connErr != nil {
		log.Fatal("连接失败", connErr)
	}
	sqlDB, dbErr := Db.DB()
	if dbErr != nil {
		sqlDB.Close()
		log.Fatal("连接失败", connErr)
	}
	// 设置连接池，空闲连接
	// 打开链接

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// 表明禁用后缀加s
	return Db
}
