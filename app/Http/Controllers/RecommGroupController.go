package Controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
)

type RecommGroupController struct {
	BaseController
}

func (group RecommGroupController) Add(c *gin.Context) {
	var params Request.GroupAddValidate
	if err := c.ShouldBind(&params); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	if nil != queryParams["ip"] {
		rangeData := queryParams["ip"]
		for _, vv := range rangeData.([]string) {
			if Utils.F.Isipv6(vv) {
				Constant.Error(c, Constant.IPV6_CANNOT_EXTRACT)
				return
			}
		}
	}
	delete(queryParams, "ip")
	queryData := Models.RecomGroupModel.Find(queryParams)
	if 0 < len(queryData) {
		Constant.Error(c, Constant.GROUP_IS_EXISTS)
		return
	}
	jwt_payload := jwt.ExtractClaims(c)
	insertData := Utils.F.StructTomap(params)
	insertData["user_id"] = int(jwt_payload["sub"].(float64))
	insertId, insertError := Models.RecomGroupModel.UpdateOrCreate(insertData)
	if insertError != nil {
		Constant.Error(c, Constant.UPDATE_FAIL)
		return
	}
	Constant.Success(c, Constant.ADD_SUCCESS, insertId)
	return
}

func (g RecommGroupController) Show(c *gin.Context) {
	var params Request.GroupShowValidate
	if err := c.ShouldBind(&params); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	insertData := Utils.F.StructTomap(params)
	if !g.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		insertData["user_id"] = g.GetUserId(c)
	}
	AllData := Models.RecomGroupModel.Show(insertData)
	Constant.Success(c, Constant.QUERY_SUCCESS, AllData)
	return
}

func (g RecommGroupController) GetAll(c *gin.Context) {
	var params Request.GroupShowValidate
	if err := c.ShouldBind(&params); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	insertData := Utils.F.StructTomap(params)
	if !g.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		insertData["user_id"] = g.GetUserId(c)
	}
	AllData := Models.RecomGroupModel.GetAllGroup(insertData)
	Constant.Success(c, Constant.QUERY_SUCCESS, AllData)
	return
}

func (g RecommGroupController) Delete(c *gin.Context) {
	var params Request.RecommGroupDelValidate
	params.IdsValidateRequest()
	if err := c.ShouldBindWith(&params, binding.FormPost); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	selectQueryParams := make(map[string]interface{})
	if !g.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["user_id"] = g.GetUserId(c)
		selectQueryParams["user_id"] = g.GetUserId(c)
	}
	if _, ok := queryParams["recom_group_name"]; ok {
		selectQueryParams["recom_group_name"] = queryParams["recom_group_name"]
	}
	if _, ok := queryParams["id"]; ok && queryParams["id"].([]int)[0] != 0 {
		queryData := Models.RecomGroupModel.Show(Utils.F.C(queryParams))
		if 0 == len(queryData) {
			Constant.Error(c, Constant.GROUP_NOT_EXISTS)
			return
		}
		if 0 == queryData["is_extract_clue_finished"].(int) {
			Constant.Error(c, Constant.CLUE_EXTRACTING)
			return
		}
		if 0 == queryData["is_extract_assets_finished"].(int) {
			Constant.Error(c, Constant.CLUE_ASSETS_EXTRACTING)
			return
		}
	} else {
		selectQueryParams["is_extract_clue_finished"] = 0
		extract_clue_arr := Models.RecomGroupModel.Show(Utils.F.C(selectQueryParams))
		if _, ok := extract_clue_arr["id"]; ok {
			Constant.Error(c, Constant.CLUE_EXTRACTING)
			return
		}
		if _, ok := selectQueryParams["is_extract_clue_finished"]; ok {
			delete(selectQueryParams, "is_extract_clue_finished")
		}
		selectQueryParams["is_extract_assets_finished"] = 0
		extract_assets_arr := Models.RecomGroupModel.Show(Utils.F.C(selectQueryParams))
		if _, ok := extract_assets_arr["id"]; ok {
			Constant.Error(c, Constant.CLUE_ASSETS_EXTRACTING)
			return
		}
		queryParams["id"] = []int{}
	}
	returnNum := Models.RecomGroupModel.Delete(Utils.F.C(queryParams))
	if returnNum == 0 {
		Constant.Error(c, Constant.DEL_ERROR)
		return
	}
	Models.ClueModel.DeleteNullGoupClue()
	Constant.Success(c, "删除成功!", "")
}

func (g RecommGroupController) Lists(c *gin.Context) {
	var params Request.GroupListsValidate
	if err := c.ShouldBindWith(&params, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	queryParams := Utils.F.StructTomap(params)
	if !g.HasRole(c, []string{Constant.ROLE_ADMIN_NAME}) {
		queryParams["user_id"] = g.GetUserId(c)
	}
	pagination := Models.GeneratePaginationFromRequest(c)
	res, countRes := Models.RecomGroupModel.StructLists(queryParams, &pagination)
	returnResult := make(map[string]interface{})
	returnResult["list"] = res
	returnResult["page"] = params.Page
	returnResult["size"] = params.Size
	returnResult["total"] = countRes
	Constant.Success(c, "查询成功!", returnResult)
	return
}
