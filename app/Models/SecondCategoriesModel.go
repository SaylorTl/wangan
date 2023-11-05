package Models

import (
	"wangxin2.0/databases"
)

var SecondCategoriesModel *secondcategoriesmodel

type secondcategoriesmodel struct {
	BaseModel
}

type Second_categories struct {
	Id       int
	Title    string `gorm:"foreignkey:Second_category_id"`
	Ancestry string
	En_title string
}

type Select_second_categories struct {
	Id       int
	Title    string
	Ancestry string
	En_title string
}

func (w secondcategoriesmodel) Find(sec_cat_id int) (res *Select_second_categories) {
	var Wxsecondcategories Select_second_categories
	databases.Db.Model(&Second_categories{}).Where("id=?", sec_cat_id).First(&Wxsecondcategories)
	return &Wxsecondcategories
}
