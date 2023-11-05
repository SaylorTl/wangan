package Controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Services"
	"wangxin2.0/app/Utils"
)

type ClueSubdomainController struct {
	BaseController
}

type ClueSubdomainAddValidate struct {
	Domains []map[string]string `json:"domains"  form:"domains" to:"domains" binding:"required"`
}

type ClueSubdomainListsValidate struct {
	Keyword string `form:"keyword" to:"keyword"`
	Page    int    `form:"page" to:"page" zh_cn:"页数" binding:"required,omitempty,min=1"`
	Size    int    `form:"size" to:"size" zh_cn:"页码" binding:"required,omitempty,min=10,max=500"`
}
type ClueSubdomainShowValidate struct {
	Domain string `form:"domain" to:"domain"`
}
type ClueSubdomainDelValidate struct {
	Keyword            string `form:"keyword" to:"keyword"`
	Subdomain_clue_ids []int  `form:"subdomain_clue_ids" to:"subdomain_clue_ids"`
}

// 域名探测添加
func (cs ClueSubdomainController) Add(c *gin.Context) { //添加
	var addparams ClueSubdomainAddValidate
	if err := c.ShouldBindJSON(&addparams); err != nil {
		log.Print(err)
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	insertParams := Utils.F.StructTomap(addparams)
	date := time.Now().Format("01-02-2006")
	for _, val := range insertParams["domains"].([]map[string]string) {
		for domain, ip := range val {
			insertData := make(map[string]interface{})
			insertData["domain"] = domain
			insertData["ip"] = ip
			insertData["status"] = 1
			insertData["path"] = Constant.SubdomainPath + date + "/" + domain
			count := Models.ClueSubDomainModel.Count(Utils.F.C(insertData))
			insertData["status"] = 0
			wait_count := Models.ClueSubDomainModel.Count(Utils.F.C(insertData))
			if count > 0 || wait_count > 0 {
				continue
			}
			insertData["status"] = 3
			insertData["created_at"] = time.Now().Format("2006-01-02 15:04:05")
			_, err := Models.ClueSubDomainModel.UpdateOrCreate(insertData)
			if err != nil {
				Constant.Error(c, Constant.ADD_FAIL)
				return
			}
		}
	}
	user_id := cs.GetUserId(c)
	Services.ClueSubdomainService.Extract(user_id)
	Constant.Success(c, Constant.ADD_SUCCESS, 0)
	return

}

// 域名探测记录列表
func (cs ClueSubdomainController) Lists(c *gin.Context) { //带分页
	var params ClueSubdomainListsValidate
	if err := c.ShouldBindWith(&params, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	pagination := Models.GeneratePaginationFromRequest(c)
	res, countRes := Models.ClueSubDomainModel.Lists(queryParams, &pagination)
	returnResult := make(map[string]interface{})
	returnResult["list"] = res
	returnResult["page"] = params.Page
	returnResult["size"] = params.Size
	returnResult["total"] = countRes
	Constant.Success(c, Constant.QUERY_SUCCESS, returnResult)
	return
}

// 域名探测展示
func (cs ClueSubdomainController) Show(c *gin.Context) { //不带分页
	var Params ClueSubdomainShowValidate
	if err := c.ShouldBindWith(&Params, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(Params)
	returnData := Models.ClueSubDomainModel.Show(queryParams)
	Constant.Success(c, Constant.QUERY_SUCCESS, returnData)
	return
}

// 域名探测记录删除
func (cs ClueSubdomainController) Delete(c *gin.Context) { //删除
	var params ClueSubdomainDelValidate
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	returnNum := Models.ClueSubDomainModel.Delete(queryParams)
	if returnNum == 0 {
		Constant.Error(c, Constant.DEL_ERROR)
		return
	}
	//更新所有分组线索数量
	Services.RecommGroupService.UpdateAllClueCount()
	Constant.Success(c, Constant.DEL_SUCCESS, "")
}
