package Models

import (
	"gorm.io/gorm"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var RecomGroupModel *recomgroupmodel

type recomgroupmodel struct {
	BaseModel
}

type Wx_assets_recom_group struct {
	Id                         int       `json:"id" gorm:"column:id;primary_key"`
	Recom_group_name           string    `json:"recom_group_name" gorm:"column:recom_group_name"`
	Is_extract_clue_finished   int       `json:"is_extract_clue_finished" gorm:"column:is_extract_clue_finished"`
	Is_extract_assets_finished int       `json:"is_extract_assets_finished" gorm:"column:is_extract_assets_finished"`
	Is_clue_update             int       `json:"is_clue_update" gorm:"column:is_clue_update"`
	Clue_count                 int       `json:"clue_count" gorm:"column:clue_count"`
	User_id                    int       `json:"user_id" gorm:"column:user_id"`
	CreatedAt                  LocalTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                  LocalTime `json:"updated_at" gorm:"column:updated_at"`
}
type Select_assetsrecomgroup struct {
	Id               int           `json:"id,omitempty" gorm:"column:id;primary_key"`
	Recom_group_name string        `json:"recom_group_name"`
	Select_clue      []Select_clue `gorm:"foreignkey:recom_group_id"`
}
type Select_clue struct {
	Assets_keyword string `json:"assets_keyword" gorm:"column:assets_keyword"`
	Recom_group_id int    `json:"recom_group_id," gorm:"column:recom_group_id"`
}

func (w recomgroupmodel) Show(queryParams map[string]interface{}) (returnResult map[string]interface{}) {
	returnData := map[string]interface{}{}
	returnRes := databases.Db.Model(&Wx_assets_recom_group{})
	if _, ok := queryParams["recom_group_name"]; ok {
		returnRes = returnRes.Where("recom_group_name = ? ", "%"+queryParams["recom_group_name"].(string)+"%")
		delete(queryParams, "recom_group_name")
	}
	returnRes = returnRes.Where(queryParams).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		panic(returnRes.Error)
	}
	return returnData
}

func (w recomgroupmodel) Find(queryParams map[string]interface{}) (returnResult map[string]interface{}) {
	returnData := map[string]interface{}{}
	returnRes := databases.Db.Model(&Wx_assets_recom_group{})
	returnRes = returnRes.Where(queryParams).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		panic(returnRes.Error)
	}
	return returnData
}

func (w recomgroupmodel) GetAllGroup(queryParams map[string]interface{}) (returnResult []map[string]interface{}) {
	returnData := []map[string]interface{}{}
	returnRes := databases.Db.Model(&Wx_assets_recom_group{})
	if _, ok := queryParams["recom_group_name"]; ok {
		returnRes = returnRes.Where("recom_group_name = ? ", "%"+queryParams["recom_group_name"].(string)+"%")
		delete(queryParams, "recom_group_name")
	}
	returnRes = returnRes.Where(queryParams).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		panic(returnRes.Error)
	}
	return returnData
}

func (w recomgroupmodel) Lists(queryParams map[string]interface{}, pagination *Pagination) (returnResult []map[string]interface{}, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	delete(queryParams, "page")
	delete(queryParams, "size")
	returnData := []map[string]interface{}{}
	returnRes := databases.Db.Model(&Wx_assets_recom_group{})
	if _, ok := queryParams["recom_group_name"]; ok {
		returnRes = returnRes.Where("recom_group_name = ? ", "%"+queryParams["recom_group_name"].(string)+"%")
		delete(queryParams, "recom_group_name")
	}
	returnRes = returnRes.Where(queryParams).Limit(pagination.Size).Offset(offset).Order(pagination.Sort).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		panic(returnRes.Error)
	}
	var count int64
	databases.Db.Model(&Wx_assets_recom_group{}).Where(queryParams).Count(&count)
	return returnData, count
}

func (w recomgroupmodel) StructLists(queryParams map[string]interface{}, pagination *Pagination) (returnResult []Wx_assets_recom_group, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	delete(queryParams, "page")
	delete(queryParams, "size")
	returnData := []Wx_assets_recom_group{}
	returnRes := databases.Db.Model(&Wx_assets_recom_group{})
	if _, ok := queryParams["recom_group_name"]; ok {
		returnRes = returnRes.Where("locate('" + queryParams["recom_group_name"].(string) + "', `recom_group_name`)")
		delete(queryParams, "recom_group_name")
	}
	returnRes = returnRes.Where(queryParams).Limit(pagination.Size).Offset(offset).Order(pagination.Sort).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		panic(returnRes.Error)
	}
	var count int64
	returnRes.Offset(-1).Limit(-1).Count(&count)
	return returnData, count
}

func (w recomgroupmodel) UpdateOrCreate(params map[string]interface{}) (int, error) {
	var selectAssets Wx_assets_recom_group
	result := databases.Db
	if params["recom_group_name"] != nil {
		result = result.Where("recom_group_name=?", params["recom_group_name"].(string))
	}
	if params["id"] != nil {
		result = result.Where("id=?", params["id"].(int))
	}
	result.Find(&selectAssets)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		selectAssets.Recom_group_name = Utils.F.If(params["recom_group_name"] == nil, "", params["recom_group_name"]).(string)
		selectAssets.Is_extract_clue_finished = Utils.F.If(params["is_extract_clue_finished"] == nil, 1, params["is_extract_clue_finished"]).(int)
		selectAssets.Is_extract_assets_finished = Utils.F.If(params["is_extract_assets_finished"] == nil, 1, params["is_extract_assets_finished"]).(int)
		selectAssets.Is_clue_update = Utils.F.If(params["is_clue_update"] == nil, 0, params["is_clue_update"]).(int)
		selectAssets.User_id = Utils.F.If(params["user_id"] == nil, 0, params["user_id"]).(int)
		insertresult := databases.Db.Create(&selectAssets)
		return selectAssets.Id, insertresult.Error
	}
	//有则更新
	updataRes := databases.Db.Model(&selectAssets).Updates(params)
	return selectAssets.Id, updataRes.Error
}

func (w recomgroupmodel) GetAll(WxAssets *Wx_assets_recom_group) (returnResult []map[string]interface{}, err error) {
	var returnData []map[string]interface{}
	returnRes := databases.Db.Model(&Wx_assets_recom_group{}).Where(&WxAssets).Find(&returnData)
	if returnRes.Error != nil && returnRes.Error != gorm.ErrRecordNotFound {
		panic(returnRes.Error)
	}
	var count int64
	databases.Db.Model(&Wx_assets_recommen{}).Where(&WxAssets).Count(&count)
	return returnData, nil
}

func (w recomgroupmodel) Delete(queryParams map[string]interface{}) (returnResult int64) {
	result := databases.Db
	if _, ok := queryParams["recom_group_name"]; ok {
		result = result.Where("locate('" + queryParams["recom_group_name"].(string) + "', `recom_group_name`)")
		delete(queryParams, "recom_group_name")
	}
	if _, ok := queryParams["id"]; ok && len(queryParams["id"].([]int)) == 0 {
		delete(queryParams, "id")
		result = result.Where("1=1")
	}
	result = result.Where(queryParams)
	result.Delete(&Wx_assets_recom_group{})
	return result.RowsAffected
}
