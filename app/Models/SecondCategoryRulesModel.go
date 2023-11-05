package Models

import (
	"wangxin2.0/databases"
)

var SecondCategoryRulesModel *secondcategoryrulesmodel

type secondcategoryrulesmodel struct {
	BaseModel
}

type Second_category_rules struct {
	Id                 int
	Second_category_id string
	Rule_id            int
	WxSecondCategories Second_categories `gorm:"foreignkey:Second_category_id"`
}

func (w secondcategoryrulesmodel) Findip(rule_id int) (res *Second_category_rules) {
	var WxSecondCategoryRules Second_category_rules
	databases.Db.Where("rule_id=?", rule_id).Where("deleted_at IS NULL").First(&WxSecondCategoryRules)
	return &WxSecondCategoryRules
}
