package Controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
)

type AssetsHistoryController struct {
}

type AssetsHistoryValidate struct {
	Ip   string `form:"ip" to:"ip" binding:"required,ip_valid_rule"`
	Page *int   `form:"page" to:"page" zh_cn:"页码" binding:"required,gt=0"`
	Size *int   `form:"size" to:"size" zh_cn:"页数" binding:"required,gt=0"`
}

func (a AssetsHistoryValidate) AssetsHistoryRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ip_valid_rule", Request.IpValidrule)
		if err != nil {
			log.Println("Error", err)
			return
		}
	}
}

// 查询历史记录
func (a AssetsHistoryController) Lists(c *gin.Context) { //带分页
	var params AssetsHistoryValidate
	params.AssetsHistoryRequest()
	if err := c.ShouldBindWith(&params, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	pagination := Models.GeneratePaginationFromRequest(c)
	res, countRes := Models.AssetsChangeHistoryModel.Lists(map[string]interface{}{"ip": params.Ip}, &pagination)
	lists := []interface{}{}
	for _, vv := range res {
		content := make(map[string]interface{})
		err := json.Unmarshal([]byte(vv.Content), &content)
		if err != nil {
			Constant.Error(c, err)
			return
		}
		created_at := vv.CreatedAt.String()
		temp := a.getContentJson(content, created_at)
		if len(temp) != 0 {
			lists = append(lists, temp)
		}
	}
	returnResult := make(map[string]interface{})
	returnResult["list"] = lists
	returnResult["page"] = params.Page
	returnResult["size"] = params.Size
	returnResult["total"] = countRes
	Constant.Success(c, "查询成功!", returnResult)
}

// 将历史记录json字符串转化为map
func (a AssetsHistoryController) getContentJson(content map[string]interface{}, created_at string) (data map[string]interface{}) {
	addPortList := []interface{}{}
	reducePortList := []interface{}{}
	addRuleList := []interface{}{}
	reduceRuleList := []interface{}{}
	temp := make(map[string]interface{})
	for kl, val := range content {
		if kl == "port_list" {
			for _, vk := range val.(map[string]interface{})["change_data"].([]interface{}) {
				if vk.(map[string]interface{})["change_type"] == "add" {
					addPortList = append(addPortList, vk.(map[string]interface{})["port"])
					temp["add_port_lists"] = addPortList
				}
				if vk.(map[string]interface{})["change_type"] == "reduce" {
					reducePortList = append(reducePortList, vk.(map[string]interface{})["port"])
					temp["reduce_port_lists"] = reducePortList
				}
			}
			if _, ok := val.(map[string]interface{})["current_port_lists"]; ok && val.(map[string]interface{})["current_port_lists"] != nil {
				temp["current_port_lists"] = val.(map[string]interface{})["current_port_lists"]
			}
		}
		if kl == "rule_list" {
			for _, vh := range val.(map[string]interface{})["change_data"].([]interface{}) {
				if vh.(map[string]interface{})["change_type"] == "add" {
					addRuleList = append(addRuleList, vh.(map[string]interface{})["title"])
					temp["add_rule_lists"] = addRuleList
				}
				if vh.(map[string]interface{})["change_type"] == "reduce" {
					reduceRuleList = append(reduceRuleList, vh.(map[string]interface{})["title"])
					temp["reduce_rule_lists"] = reduceRuleList
				}
			}
			if _, ok := val.(map[string]interface{})["current_rule_lists"]; ok && val.(map[string]interface{})["current_rule_lists"] != nil {
				temp["current_rule_lists"] = val.(map[string]interface{})["current_rule_lists"]
			}
		}
		if kl == "note_list" {
			temp["note_lists"] = val
		}
	}
	if len(temp) != 0 {
		temp["created_at"] = created_at
	}
	return temp
}
