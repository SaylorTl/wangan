package Constant

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const READ_JSON_ERROR string = "json读取失败"
const CONF_READ_ERROR string = "配置文件读取失败"
const INIT_REDIS_ERROR string = "init redis client failed"
const INIT_ROUTE_ERROR string = "init route client failed"
const CLOSE_FILE_ERROR string = "close file client failed"
const CLOSE_REDIS_ERROR string = "close redis client failed"
const INIT_TRANS_ERROR string = "init translate client failed"
const CLOSE_DB_ERROR string = "close db client failed"
const ADD_SUCCESS string = "添加成功！"
const ADD_FAIL string = "添加失败！"
const UPDATE_FAIL string = "更新失败！"
const QUERY_FAIL string = "查询失败！"
const QUERY_SUCCESS string = "查询成功！"
const INIT_ES_ERROR string = "es初始化失败"
const FIND_CLUE_ERROR string = "提炼失败"
const DEL_ERROR string = "删除失败"
const DEL_SUCCESS string = "删除成功"
const PAREMS_ERROR string = "参数错误"
const CLUE_EXTRACTING string = "线索正在提炼中"
const CLUE_ASSETS_EXTRACTING string = "资产正在推荐中"
const GROUP_NOT_EXISTS string = "分组不存在！"
const GROUP_IS_EXISTS string = "分组已存在！"
const ES_QUERY_FAIL string = "ES查询失败！"
const JSON_UNENCODE_FAIL string = "Json解析失败！"
const READ_YMAL_FAIL string = "读取文件config.yml发生错误"
const UNENCODE_YMAL_FAIL string = "解析文件config.yml发生错误"
const CLUE_RICH_CLOSE string = "资产富化功能关闭"

const IPV6_CANNOT_EXTRACT string = "ipv6地址无法提取"
const GROUP_CANNOT_EMPTY string = "分组数据不能为空"
const EXTRACT_SUCCESS string = "智能发现成功"
const wangan_SYNC_TASK_EXISTS string = "当前有任务正在执行"

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

func Error(c *gin.Context, msg any) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": msg,
		"data":    struct{}{},
	})
}
