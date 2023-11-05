package Controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Services"
	"wangxin2.0/app/Utils"
)

type ClueController struct {
	BaseController
}

func (a ClueController) Add(c *gin.Context) {
	var params Request.ClueAddValidate
	params.ClueRequest()
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	queryParams["user_id"] = a.GetUserId(c)
	tmpParams := make(map[string]interface{})
	tmpParams["id"] = queryParams["asset_group_id"].(int)
	returnData := Models.RecomGroupModel.Show(tmpParams)
	if returnData["is_extract_clue_finished"] == 0 {
		Constant.Error(c, Constant.CLUE_EXTRACTING)
		return
	}
	res, msg := Services.ClueService.Extract(queryParams)
	if !res {
		Constant.Error(c, msg)
		return
	}
	Constant.Success(c, msg, "")
	return
}

func (a ClueController) Lists(c *gin.Context) {
	var params Request.GetClueListValidate
	params.GetClueValidateRequest()
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	if !a.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["uid"] = a.GetUserId(c)
	}
	pagination := Models.GeneratePaginationFromRequest(c)
	res, countRes := Services.ClueService.Find(queryParams, pagination)
	returnResult := make(map[string]interface{})
	returnResult["list"] = res
	returnResult["page"] = params.Page
	returnResult["size"] = params.Size
	returnResult["total"] = countRes
	Constant.Success(c, "查询成功!", returnResult)
	return
}

func (a ClueController) Show(c *gin.Context) {
	var assetParams Request.GetClueShowValidate
	if err := c.ShouldBindWith(&assetParams, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(assetParams)
	if !a.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["uid"] = a.GetUserId(c)
	}
	returnData := Models.ClueModel.GetAll(queryParams)
	Constant.Success(c, "查询成功!", returnData)
	return
}

func (a ClueController) Delete(c *gin.Context) {
	var assetParams Request.ClueDelValidate
	assetParams.ClueDelRequest()
	if err := c.ShouldBindWith(&assetParams, binding.FormPost); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(assetParams)
	if !a.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["uid"] = a.GetUserId(c)
	}
	returnNum := Models.ClueModel.Delete(queryParams)
	if returnNum == 0 {
		Constant.Error(c, Constant.DEL_ERROR)
		return
	}
	//更新所有分组线索数量
	Services.RecommGroupService.UpdateAllClueCount()
	Constant.Success(c, "删除成功!", "")
}
