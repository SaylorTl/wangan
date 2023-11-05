package Models

import (
	"wangxin2.0/databases"
)

var RulesModel *rulesmodel

type rulesmodel struct {
	BaseModel
}

type Rules struct {
	Id                    int
	Company               string
	Product               string
	Rule                  int
	Taged                 int
	Company_taged         int
	Level_code            int
	Soft_hard_code        int
	En_product            string
	En_company            string
	WxSecondCategoryRules Second_category_rules `gorm:"foreignKey:rule_id"`
}

func (w rulesmodel) Findip(product_name string) (res *Rules) {
	var WxRules Rules
	databases.Db.Where("product=?", product_name).First(&WxRules)
	return &WxRules
}

func (w rulesmodel) FindOne(product_name string) (res *Rules) {
	var WxRules Rules
	databases.Db.Preload("WxSecondCategoryRules").
		Preload("WxSecondCategoryRules.WxSecondCategories").
		Where("product=?", product_name).First(&WxRules)
	return &WxRules
}
