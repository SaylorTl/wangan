package Request

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type GroupAddValidate struct {
	Ip         []string `form:"ip" to:"ip" binding:"omitempty"`
	Group_name string   `form:"group_name" zh_cn:"组名" to:"recom_group_name" binding:"required,alphanumunicode,max=30"`
}

type GroupListsValidate struct {
	Size             *int   `form:"size" zh_cn:"" to:"size" zh_cn:"页码" binding:"required,omitempty,min=10,max=500"`
	Page             *int   `form:"page" to:"page" zh_cn:"页数" binding:"required,omitempty,min=1"`
	Recom_group_name string `form:"recom_group_name" zh_cn:"组名" to:"recom_group_name" binding:"omitempty,max=30"`
}

type GroupShowValidate struct {
	Recom_group_name           string `form:"recom_group_name" zh_cn:"组名" to:"recom_group_name" binding:"omitempty"`
	Is_extract_clue_finished   *int   `form:"is_extract_clue_finished" zh_cn:"线索是否提炼完成" to:"is_extract_clue_finished" binding:"omitempty,oneof=0 1"`
	Is_extract_assets_finished *int   `form:"is_extract_assets_finished" zh_cn:"资产是否提炼完成" to:"is_extract_assets_finished" binding:"omitempty,oneof=0 1"`
	Is_clue_update             *int   `form:"is_clue_update" to:"is_clue_update" zh_cn:"线索是否更新" binding:"omitempty,oneof=0 1"`
	Id                         int    `form:"id" to:"id" zh_cn:"组id" binding:"omitempty"`
}

type RecommGroupDelValidate struct {
	Ids              []int  `form:"ids" to:"id" zh_cn:"组id" binding:"required"`
	Recom_group_name string `form:"recom_group_name" zh_cn:"组名" to:"recom_group_name" binding:"omitempty"`
}

func (i RecommGroupDelValidate) IdsValidateRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("IdsVal", IdsVal)
		if err != nil {
			return
		}
	}
}
