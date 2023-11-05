package Request

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

type ClueAddValidate struct {
	Ip                []string `form:"ip" to:"ip" binding:"omitempty"`
	Is_add_clue_asset *int     `form:"is_add_clue_asset"  to:"is_add_clue_asset" zh_cn:"线索资产"  binding:"required,oneof=0 1"`
	Asset_group_id    *int     `form:"asset_group_id" to:"asset_group_id" zh_cn:"分组"  binding:"required"`
	Keyword           string   `form:"keyword" to:"keyword"  zh_cn:"关键字" binding:"omitempty"`
	Ip_type           bool     `form:"ip_type" to:"ip_type" zh_cn:"ip类型" binding:"omitempty"`
	Has_log4j         string   `form:"has_log4j" to:"has_log4j" binding:"omitempty"`
	Created_at_range  []string `form:"created_at_range" to:"created_at_range" zh_cn:"创建时间" binding:"omitempty"`
	Update_at_range   []string `form:"update_at_range" to:"update_at_range" zh_cn:"更新时间" binding:"omitempty"`
}

type GetClueListValidate struct {
	Page           int    `form:"page" to:"page" zh_cn:"页数" binding:"required,omitempty,min=1"`
	Size           int    `form:"size" to:"size" zh_cn:"页码" binding:"required,omitempty,min=10,max=500"`
	Cid            int    `form:"cid" to:"cid" zh_cn:"单位"  binding:"omitempty"`
	Ip             string `form:"ip" to:"ip" binding:"omitempty,ip_valid_rule"`
	Ip_type        *int   `form:"ip_type" to:"ip_type" zh_cn:"ip类型" binding:"omitempty,oneof=0 1 2"`
	Is_device      *int   `form:"is_device" to:"is_device" zh_cn:"设备" binding:"omitempty,oneof=0 1"`
	Is_sign        *int   `form:"is_sign" to:"is_sign" zh_cn:"签名" binding:"omitempty,oneof=0 1"`
	Keyword        string `form:"keyword" to:"keyword" zh_cn:"关键字" binding:"omitempty,keyword_val"`
	Recom_group_id int    `form:"recom_group_id" to:"recom_group_id" zh_cn:"分组" binding:"omitempty"`
	Assettype      *int   `form:"assettype" to:"type" zh_cn:"资产类型" binding:"omitempty"`
}

type GetClueShowValidate struct {
	Cid            int    `form:"cid" to:"cid" zh_cn:"单位" binding:"omitempty"`
	Ip             string `form:"ip" to:"ip" binding:"omitempty,ip_valid_rule"`
	Ip_type        *int   `form:"ip_type"  to:"ip_type" zh_cn:"ip类型" binding:"omitempty,oneof=0 1 2"`
	Is_device      *int   `form:"is_device" to:"is_device" zh_cn:"设备" binding:"omitempty,oneof=0 1"`
	Is_sign        *int   `form:"is_sign" to:"is_sign" zh_cn:"签名" binding:"omitempty,oneof=0 1"`
	Keyword        string `form:"keyword" to:"keyword" zh_cn:"关键字" binding:"omitempty,keyword_val"`
	Recom_group_id int    `form:"recom_group_id" to:"recom_group_id" zh_cn:"分组" binding:"omitempty"`
	Assettype      *int   `form:"assettype"  to:"type" zh_cn:"资产类型" binding:"omitempty"`
}

type ClueDelValidate struct {
	Ids            []int  `form:"ids"  to:"id" zh_cn:"线索" binding:"required"`
	Cid            int    `form:"cid" to:"cid" zh_cn:"单位" binding:"omitempty"`
	Ip             string `form:"ip" to:"ip" binding:"omitempty,ip_valid_rule"`
	Ip_type        *int   `form:"ip_type"  to:"ip_type" zh_cn:"ip类型" binding:"omitempty,oneof=0 1 2"`
	Is_device      *int   `form:"is_device" to:"is_device" zh_cn:"设备" binding:"omitempty,oneof=0 1"`
	Is_sign        *int   `form:"is_sign" to:"is_sign" zh_cn:"签名" binding:"omitempty,oneof=0 1"`
	Keyword        string `form:"keyword" to:"keyword" zh_cn:"关键字" binding:"omitempty,keyword_val"`
	Recom_group_id int    `form:"recom_group_id" to:"recom_group_id" zh_cn:"分组" binding:"omitempty"`
	Assettype      *int   `form:"assettype"  to:"type" zh_cn:"资产类型" binding:"omitempty"`
}

func (a GetClueShowValidate) ClueRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ip_valid_rule", IpDataValidrule)
		if err != nil {
			log.Println("Error", err)
			return
		}
		keyerr := v.RegisterValidation("keyword_val", KeywordVal)
		if keyerr != nil {
			log.Println("Error", err)
			return
		}
	}
}

func (a ClueAddValidate) ClueRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ip_valid_rule", IpDataValidrule)
		if err != nil {
			log.Println("Error", err)
			return
		}
		keyerr := v.RegisterValidation("keyword_val", KeywordVal)
		if keyerr != nil {
			log.Println("Error", err)
			return
		}
	}
}

func (a GetClueListValidate) GetClueValidateRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ip_valid_rule", IpDataValidrule)
		if err != nil {
			log.Println("Error", err)
			return
		}
		keyerr := v.RegisterValidation("keyword_val", KeywordVal)
		if keyerr != nil {
			log.Println("Error", err)
			return
		}
	}
}

func (a ClueDelValidate) ClueDelRequest() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ip_valid_rule", IpDataValidrule)
		if err != nil {
			log.Println("Error", err)
			return
		}
		keyerr := v.RegisterValidation("keyword_val", KeywordVal)
		if keyerr != nil {
			log.Println("Error", err)
			return
		}
	}
}
