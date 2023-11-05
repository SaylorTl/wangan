package Controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Services"
	"wangxin2.0/app/Utils"
)

type ClueAssetsController struct {
	BaseController
}

// 展示列表
func (cl ClueAssetsController) Lists(c *gin.Context) { //带分页
	var params Request.ClueAssetListValidate
	params.ClueAssetValidateRequest()
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	//将struct结果转换为map
	queryParams := Utils.F.StructTomap(params)
	log.Print(queryParams)
	if !cl.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["user_id"] = cl.GetUserId(c)
	}
	pagination := Models.GeneratePaginationFromRequest(c)
	res, countRes := Models.ClueAssetsModel.GetAllClueAsset(queryParams, &pagination)
	returnResult := make(map[string]interface{})
	returnResult["list"] = res
	returnResult["page"] = params.Page
	returnResult["size"] = params.Size
	returnResult["total"] = countRes
	Constant.Success(c, "查询成功!", returnResult)

}

// 展示单个结果
func (cl ClueAssetsController) Show(c *gin.Context) { //不带分页
	var params Request.ClueAssetShowValidate
	if err := c.ShouldBind(&params); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	if !cl.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["user_id"] = cl.GetUserId(c)
	}
	returnData := Models.ClueAssetsModel.GetAll(queryParams)
	Constant.Success(c, "查询成功!", returnData)

}

// 提炼资产
func (cl ClueAssetsController) AddAssetByClue(c *gin.Context) {
	var params Request.ClueAssetsAddValidate
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	updateGroupData := make(map[string]interface{})
	updateGroupData["is_extract_assets_finished"] = 0
	updateGroupData["id"] = queryParams["recom_group_id"]
	_, err := Models.RecomGroupModel.UpdateOrCreate(updateGroupData)
	if err != nil {
		Constant.Error(c, err)
		return
	}
	is_success := Services.ClueAssetsService.Find(queryParams)
	if !is_success {
		Constant.Error(c, Constant.FIND_CLUE_ERROR)
		return
	}
	Constant.Success(c, "推荐成功!", "")
}

// 删除推荐资产
func (cl ClueAssetsController) Delete(c *gin.Context) { //不带分页
	var params Request.ClueAssetsDelValidate
	if err := c.ShouldBindWith(&params, binding.FormPost); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	if !cl.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["user_id"] = cl.GetUserId(c)
	}
	returnNum := Models.ClueAssetsModel.Delete(queryParams)
	if returnNum == 0 {
		Constant.Error(c, Constant.DEL_ERROR)
		return
	}
	Constant.Success(c, "删除成功!", "")
}
