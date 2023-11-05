package Models

import (
	"wangxin2.0/databases"
)

var CategoriesModel *categoriesmodel

type categoriesmodel struct {
	BaseModel
}

type Categories struct {
	Id    int
	Title string `gorm:"foreignkey:Second_category_id"`
}

func (w categoriesmodel) Find(catid int) (res *Categories) {
	var Wxcategories Categories
	databases.Db.Where("id=?", catid).First(&Wxcategories)
	return &Wxcategories
}
