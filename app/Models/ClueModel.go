package Models

import (
	"gorm.io/gorm"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var ClueModel *cluemodel

type cluemodel struct {
	BaseModel
}

type Wx_assets_recommen struct {
	Id                 int                   `json:"id" gorm:"column:id;primary_key"`
	Uid                int                   `json:"uid" gorm:"column:uid"`
	Cid                int                   `json:"cid" gorm:"column:cid"`
	Assets_id          int                   `json:"assets_id" gorm:"column:assets_id"`
	Assets_keyword     string                `json:"assets_keyword" gorm:"column:assets_keyword"`
	Type               int                   `json:"type" gorm:"column:type"`
	Keyword            string                `json:"keyword" gorm:"column:keyword"`
	Is_device          int                   `json:"is_device" gorm:"column:is_device"`
	Is_sign            int                   `json:"is_sign" gorm:"column:is_sign"`
	Status             int                   `json:"status" gorm:"column:status"`
	CreatedAt          LocalTime             `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          LocalTime             `json:"updated_at" gorm:"column:updated_at"`
	Recom_group_id     int                   `json:"recom_group_id" gorm:"column:recom_group_id"`
	Wxassetsrecomgroup Wx_assets_recom_group `gorm:"foreignkey:recom_group_id"`
}

type Select_cluerecommen struct {
	Id      int    `json:"id" gorm:"column:id;primary_key"`
	Keyword string `json:"keyword" gorm:"column:keyword"`
}

func (w cluemodel) Lists(queryParams map[string]interface{}, pagination *Pagination) (returnResult []Wx_assets_recommen, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	delete(queryParams, "page")
	delete(queryParams, "size")
	returnData := []Wx_assets_recommen{}
	returnRes := databases.Db.Preload("Wxassetsrecomgroup")
	if _, ok := queryParams["keyword"]; ok {
		returnRes = returnRes.Where("locate('" + queryParams["keyword"].(string) + "', `keyword`)")
		delete(queryParams, "keyword")
	}
	returnRes.Where(queryParams).Limit(pagination.Size).Offset(offset).Order(pagination.Sort).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	var count int64
	returnRes.Offset(-1).Limit(-1).Count(&count)
	return returnData, count
}

func (w cluemodel) GetAllRecommand(queryParams map[string]interface{}, pagination *Pagination) (returnResult []Wx_assets_recommen, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	delete(queryParams, "page")
	delete(queryParams, "size")
	returnData := []Wx_assets_recommen{}
	returnRes := databases.Db.Preload("Wxassetsrecomgroup").Where(queryParams).Limit(pagination.Size).Offset(offset).Order(pagination.Sort).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	var count int64
	returnRes.Offset(-1).Limit(-1).Count(&count)
	return returnData, count
}

func (w cluemodel) Count(queryParams map[string]interface{}) int64 {
	var count int64
	databases.Db.Model(&Wx_assets_recommen{}).Preload("Wxassetsrecomgroup").Where(queryParams).Count(&count)
	return count
}

func (w cluemodel) UpdateOrCreate(params map[string]interface{}) (int, error) {
	var Wxassetsrecommen Wx_assets_recommen
	result := databases.Db.Where("keyword=? and recom_group_id=?", params["keyword"].(string), params["recom_group_id"].(int)).First(&Wxassetsrecommen)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		log.Print(params)
		Wxassetsrecommen.Recom_group_id = Utils.F.If(params["recom_group_id"].(int) == 0, "", params["recom_group_id"]).(int)
		Wxassetsrecommen.Keyword = Utils.F.If(params["keyword"].(string) == "", "", params["keyword"]).(string)
		Wxassetsrecommen.Assets_keyword = Utils.F.If(params["assets_keyword"].(string) == "", "", params["assets_keyword"]).(string)
		Wxassetsrecommen.Type = Utils.F.If(params["type"].(int) == 0, 1, params["type"]).(int)
		Wxassetsrecommen.Uid = Utils.F.If(params["uid"].(int) == 0, 0, params["uid"]).(int)
		Wxassetsrecommen.Cid = Utils.F.If(params["cid"].(int) == 0, 0, params["cid"]).(int)
		Wxassetsrecommen.Status = Utils.F.If(params["status"].(int) == 0, 0, params["status"]).(int)
		Wxassetsrecommen.Assets_id = Utils.F.If(params["assets_id"].(int) == 0, 0, params["assets_id"]).(int)
		insertRes := databases.Db.Create(&Wxassetsrecommen)
		return Wxassetsrecommen.Id, insertRes.Error
	}
	//有则更新
	updataRes := databases.Db.Model(&Wxassetsrecommen).Updates(params)
	return Wxassetsrecommen.Id, updataRes.Error
}

func (w cluemodel) InsertClue(params map[string]interface{}) error {
	var Wxassetsrecommen Wx_assets_recommen
	result := databases.Db.Where("keyword=? and recom_group_id=?", params["keyword"].(string), params["recom_group_id"].(int)).First(&Wxassetsrecommen)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxassetsrecommen.Recom_group_id = Utils.F.If(params["recom_group_id"].(int) == 0, "", params["recom_group_id"]).(int)
		Wxassetsrecommen.Keyword = Utils.F.If(params["keyword"].(string) == "", "", params["keyword"]).(string)
		Wxassetsrecommen.Assets_keyword = Utils.F.If(params["assets_keyword"].(string) == "", "", params["assets_keyword"]).(string)
		Wxassetsrecommen.Type = Utils.F.If(params["type"].(int) == 0, 1, params["type"]).(int)
		Wxassetsrecommen.Uid = Utils.F.If(params["uid"].(int) == 0, 0, params["uid"]).(int)
		Wxassetsrecommen.Cid = Utils.F.If(params["cid"].(int) == 0, 0, params["cid"]).(int)
		Wxassetsrecommen.Status = Utils.F.If(params["status"].(int) == 0, 0, params["status"]).(int)
		Wxassetsrecommen.Assets_id = Utils.F.If(params["assets_id"].(int) == 0, 0, params["assets_id"]).(int)
		databases.Db.Create(&Wxassetsrecommen)
	}
	return nil
}

func (w cluemodel) GetAll(queryParams map[string]interface{}, whereInParams ...map[string]interface{}) (returnResult []Wx_assets_recommen) {
	returnData := []Wx_assets_recommen{}
	result := databases.Db
	if queryParams != nil {
		result = result.Where(queryParams)
	}
	if whereInParams != nil {
		for key, val := range whereInParams[0] {
			result = result.Where(key+" IN (?)", val)
		}
	}
	result = result.Find(&returnData)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		panic(result.Error)
	}
	return returnData

}

func (w cluemodel) BuildD01Param(culeType int, content interface{}) (returnres map[string]interface{}) {
	returnData := make(map[string]interface{})
	currentType := 0
	if culeType == Constant.TYPE_DOMAIN {
		currentType = Constant.AssetEnrichment_TYPE_DOMAIN
	}
	if culeType == Constant.TYPE_ICP {
		currentType = Constant.AssetEnrichment_TYPE_CERT
	}
	if culeType == Constant.TYPE_CERT {
		currentType = Constant.AssetEnrichment_TYPE_ICP
	}
	returnData["type"] = currentType
	returnData["content"] = content
	return returnData

}

func (w cluemodel) Delete(queryParams map[string]interface{}) (returnResult int64) {
	result := databases.Db
	if _, ok := queryParams["id"]; ok && queryParams["id"].([]int)[0] == 0 {
		delete(queryParams, "id")
		result = result.Where("1=1")
	}
	if _, ok := queryParams["keyword"]; ok {
		result = result.Where("locate('" + queryParams["keyword"].(string) + "', `keyword`)")
		delete(queryParams, "keyword")
	}
	result = result.Where(queryParams)
	result.Delete(&Wx_assets_recommen{})
	return result.RowsAffected
}

func (w cluemodel) DeleteNullGoupClue() {
	returnData := []int{}
	databases.Db.Model(&Wx_assets_recom_group{}).Select("id").Find(&returnData)
	result := databases.Db
	result = result.Where(" recom_group_id not in ?", returnData)
	result.Delete(&Wx_assets_recommen{})
}
