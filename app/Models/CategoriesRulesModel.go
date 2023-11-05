package Models

import (
	"wangxin2.0/databases"
)

var CategoriesRulesModel *categoriesrulesmodel

type categoriesrulesmodel struct {
	BaseModel
}

type Categories_rules struct {
	Id           int
	Rule_id      int
	Category_id  int
	WxCategories Categories `gorm:"foreignkey:Category_id"`
}

func (w categoriesrulesmodel) Find(rule_id int) (res *Categories_rules) {
	var Wxcategoriesrules Categories_rules
	databases.Db.Where("rule_id=?", rule_id).First(&Wxcategoriesrules)
	return &Wxcategoriesrules
}
