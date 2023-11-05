package Request

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ClueAssetListValidate struct {
	Company_id       *int     `json:"cid,omitempty" form:"company_id" to:"company_id" zh_cn:"公司" binding:"omitempty"`
	Clud_id          *int     `json:"clud_id,omitempty" form:"clud_id" to:"clud_id"  zh_cn:"线索" binding:"omitempty"`
	Recom_group_id   *int     `json:"recom_group_id,omitempty" form:"recom_group_id" zh_cn:"分组" to:"recom_group_id" binding:"omitempty"`
	Page             *int     `json:"page,omitempty" form:"page" to:"page" zh_cn:"页数" binding:"required,min=1"`
	Size             *int     `json:"size,omitempty" form:"size" to:"size" zh_cn:"页码" binding:"required,min=10,max=500"`
	Is_import        *int     `json:"is_import,omitempty" form:"is_import" to:"is_import"  zh_cn:"是否导入" binding:"omitempty,oneof=0 1"`
	Is_scan          *int     `json:"is_scan,omitempty" form:"is_scan" to:"is_scan" zh_cn:"是否扫描" binding:"omitempty,oneof=0 1"`
	Ip_type          *int     `json:"ip_type,omitempty" form:"ip_type" to:"ip_type" zh_cn:"ip类型" binding:"omitempty,oneof=0 1 2"`
	Ip               string   `json:"ip,omitempty" form:"ip" to:"ip" binding:"omitempty"`
	Keyword          string   `json:"keyword,omitempty" form:"keyword" to:"keyword" binding:"omitempty"`
	Created_at_range []string `json:"created_at_range,omitempty" form:"created_at_range" to:"created_at_range"  zh_cn:"创建时间" binding:"omitempty,TimeRangeVal"`
	Updated_at_range []string `json:"updated_at_range,omitempty" form:"updated_at_range" to:"updated_at_range" zh_cn:"更新时间" binding:"omitempty,TimeRangeVal"`
}

type ClueAssetShowValidate struct {
	Cid            int  `form:"cid" to:"cid" zh_cn:"单位" binding:"omitempty"`
	Clud_id        int  `form:"clud_id" to:"clud_id" zh_cn:"线索" binding:"omitempty"`
	Recom_group_id int  `form:"recom_group_id" to:"recom_group_id" zh_cn:"分组" binding:"required,omitempty"`
	Is_import      *int `form:"is_import" to:"is_import" zh_cn:"是否导入" binding:"omitempty,oneof=0 1"`
	Is_scan        *int `form:"is_scan" to:"is_scan" zh_cn:"是否扫描" binding:"omitempty,oneof=0 1"`
}

type ClueAssetsAddValidate struct {
	Recom_group_id *int  `form:"recom_group_id" to:"recom_group_id" zh_cn:"分组" binding:"required,omitempty"`
	Clue_ids       []int `form:"clue_ids" to:"clue_id" zh_cn:"线索" binding:"omitempty"`
}

type ClueAssetsDelValidate struct {
	Recomm_ids     []int `form:"recomm_ids" to:"recomm_id" zh_cn:"分组" binding:"omitempty"`
	Cid            int   `form:"cid" to:"cid" zh_cn:"单位" binding:"omitempty"`
	Clud_id        int   `form:"clud_id" to:"clud_id" zh_cn:"线索" binding:"omitempty"`
	Recom_group_id int   `form:"recom_group_id" to:"recom_group_id" zh_cn:"分组" binding:"omitempty"`
	Is_import      *int  `form:"is_import" to:"is_import" zh_cn:"是否导入" binding:"omitempty,oneof=0 1"`
	Is_scan        *int  `form:"is_scan" to:"is_scan" zh_cn:"是否扫描" binding:"omitempty,oneof=0 1"`
}

type Clud_idsValidate struct {
	Clud_ids []int `form:"clud_ids" to:"clud_id" zh_cn:"线索" binding:"required,IdsVal"`
}

func (b ClueAssetListValidate) ClueAssetValidateRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("IpValidrule", IpValidrule)
		if err != nil {
			return
		}
		timeerr := v.RegisterValidation("TimeRangeVal", TimeRangeVal)
		if timeerr != nil {
			return
		}
	}
}

func (b Clud_idsValidate) Clud_idsValidateRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("IdsVal", IdsVal)
		if err != nil {
			return
		}
	}
}
