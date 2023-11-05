package Models

import (
	"gorm.io/gorm"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/databases"
)

var AssetsChangeHistoryModel *assetschangehistorymodel

type assetschangehistorymodel struct {
	BaseModel
}

type Wx_assets_change_history struct {
	Id        int       `json:"id" gorm:"column:id;primary_key"`
	Ip        string    `json:"ip" gorm:"column:ip"`
	Content   string    `json:"content" gorm:"column:content"`
	CreatedAt LocalTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt LocalTime `json:"updated_at" gorm:"column:updated_at"`
}

func (w assetschangehistorymodel) Lists(queryParams map[string]interface{}, pagination *Pagination) (returnResult []Wx_assets_change_history, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	returnData := []Wx_assets_change_history{}
	returnRes := databases.Db.Model(&Wx_assets_change_history{}).Where(queryParams).Limit(pagination.Size).Offset(offset).Order(pagination.Sort).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	var count int64
	databases.Db.Model(&Wx_assets_change_history{}).Where(queryParams).Count(&count)
	return returnData, count

}
