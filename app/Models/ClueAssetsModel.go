package Models

import (
	"gorm.io/gorm"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var ClueAssetsModel *clueassetsmodel

type clueassetsmodel struct {
	BaseModel
}

type Wx_clue_assets struct {
	Recomm_id          int                   `json:"recomm_id" gorm:"column:recomm_id;primary_key"`
	Company_id         int                   `json:"company_id" gorm:"column:company_id"`
	Ip                 string                `json:"ip" gorm:"column:ip"`
	Recom_group_id     int                   `json:"recom_group_id" gorm:"column:recom_group_id"`
	Clue_id            int                   `json:"clue_id" gorm:"column:clue_id"`
	Is_import          int                   `json:"is_import" gorm:"column:is_import"`
	Is_scan            int                   `json:"is_scan" gorm:"column:is_scan"`
	Ip_type            int                   `json:"ip_type" gorm:"column:ip_type"`
	Assets_keyword     string                `json:"assets_keyword" gorm:"column:assets_keyword"`
	User_id            int                   `json:"user_id" gorm:"column:user_id"`
	CreatedAt          LocalTime             `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          LocalTime             `json:"updated_at" gorm:"column:updated_at"`
	Wxassetsrecomgroup Wx_assets_recom_group `gorm:"foreignkey:recom_group_id"`
	Wxassetsrecommen   Wx_assets_recommen    `gorm:"foreignkey:clue_id"`
	Wxcompanies        Wx_companies          `gorm:"foreignkey:company_id"`
}

type Select_clue_assets struct {
	Recomm_id          int                     `json:"recomm_id"`
	Company_id         int                     `json:"company_id"`
	Ip                 string                  `json:"ip"`
	Clue_id            int                     `json:"clue_id"`
	Ip_type            int                     `json:"ip_type" gorm:"column:ip_type"`
	Recom_group_id     int                     `json:"recom_group_id,omitempty" gorm:"column:recom_group_id"`
	Is_import          int                     `json:"is_import"`
	Is_scan            int                     `json:"is_scan"`
	CreatedAt          LocalTime               `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          LocalTime               `json:"updated_at" gorm:"column:updated_at"`
	Wxassetsrecomgroup Select_assetsrecomgroup `gorm:"foreignkey:recom_group_id"`
	Wxassetsrecommen   Select_cluerecommen     `gorm:"foreignkey:clue_id"`
	Wxcompanies        Select_companies        `gorm:"foreignkey:company_id"`
}

func (w clueassetsmodel) GetAllClueAsset(queryParams map[string]interface{}, pagination *Pagination) (returnResult []Select_clue_assets, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	returnData := []Select_clue_assets{}
	delete(queryParams, "page")
	delete(queryParams, "size")
	returnDB := databases.Db.Model(&Wx_clue_assets{}).
		Preload("Wxassetsrecomgroup", func(db *gorm.DB) *gorm.DB { return db.Model(&Wx_assets_recom_group{}) }).
		Preload("Wxassetsrecommen", func(db *gorm.DB) *gorm.DB { return db.Model(&Wx_assets_recommen{}) }).
		Preload("Wxcompanies", func(db *gorm.DB) *gorm.DB { return db.Model(&Wx_companies{}).Where("deleted_at IS NULL") })
	if _, ok := queryParams["ip"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["ip"].(string) + "', `ip`)")
		delete(queryParams, "ip")
	}
	if _, ok := queryParams["keyword"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["keyword"].(string) + "', `ip`)")
		delete(queryParams, "keyword")
	}
	if nil != queryParams["created_at_range"] {
		created_at_range := queryParams["created_at_range"].([]string)
		returnDB.Where("created_at > ? and created_at < ?", created_at_range[0], created_at_range[1])
		delete(queryParams, "created_at_range")
	}
	if nil != queryParams["updated_at_range"] {
		update_at_range := queryParams["updated_at_range"].([]string)
		returnDB.Where("updated_at > ? and updated_at < ?", update_at_range[0], update_at_range[1])
		delete(queryParams, "updated_at_range")
	}
	returnDB.Where(queryParams).Limit(pagination.Size).Offset(offset).Order(pagination.Sort).Find(&returnData)
	if returnDB.Error != nil && returnDB.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	var count int64
	returnDB.Offset(-1).Limit(-1).Count(&count)
	return returnData, count

}

func (w clueassetsmodel) UpdateOrCreate(params map[string]interface{}) (int, error) {
	var Wxclueassets Wx_clue_assets
	result := databases.Db.Where("clue_id=? and ip=?", params["clue_id"].(int), params["ip"].(string)).First(&Wxclueassets)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxclueassets.Company_id = Utils.F.If(params["company_id"].(int) == 0, 0, params["company_id"]).(int)
		Wxclueassets.Ip = Utils.F.If(params["ip"].(string) == "", "", params["ip"]).(string)
		Wxclueassets.Recom_group_id = Utils.F.If(params["recom_group_id"].(int) == 0, "", params["recom_group_id"]).(int)
		Wxclueassets.Clue_id = Utils.F.If(params["clue_id"].(int) == 0, "", params["clue_id"]).(int)
		insertRes := databases.Db.Model(Wx_clue_assets{}).Create(&Wxclueassets)
		return Wxclueassets.Recomm_id, insertRes.Error
	}
	//有则更新
	updataRes := databases.Db.Model(&Wxclueassets).Updates(params)
	return Wxclueassets.Recomm_id, updataRes.Error
}

func (w clueassetsmodel) InsertClueAsset(params map[string]interface{}) (int, error) {
	var Wxclueassets Wx_clue_assets
	result := databases.Db.Select("clue_id").Where("clue_id=? and ip=?", params["clue_id"].(int), params["ip"].(string)).First(&Wxclueassets)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxclueassets.Company_id = Utils.F.If(params["company_id"].(int) == 0, 0, params["company_id"]).(int)
		Wxclueassets.Ip = Utils.F.If(params["ip"].(string) == "", "", params["ip"]).(string)
		Wxclueassets.Recom_group_id = Utils.F.If(params["recom_group_id"].(int) == 0, "", params["recom_group_id"]).(int)
		Wxclueassets.Clue_id = Utils.F.If(params["clue_id"].(int) == 0, "", params["clue_id"]).(int)
		isipv6 := Utils.F.Isipv6(params["ip"].(string))
		Wxclueassets.Ip_type = Utils.F.If(isipv6 == true, 1, 0).(int)
		Wxclueassets.User_id = Utils.F.If(params["user_id"].(int) == 0, 0, params["user_id"]).(int)
		Wxclueassets.Assets_keyword = Utils.F.If(params["assets_keyword"].(string) == "", "", params["assets_keyword"]).(string)
		insertRes := databases.Db.Model(Wx_clue_assets{}).Create(&Wxclueassets)
		return Wxclueassets.Recomm_id, insertRes.Error
	}
	//有则更新
	return Wxclueassets.Recomm_id, nil
}

func (w clueassetsmodel) GetAll(queryParams map[string]interface{}) (returnResult []Select_clue_assets) {
	returnData := []Select_clue_assets{}
	result := databases.Db.Model(&Wx_clue_assets{}).
		Preload("Wxassetsrecomgroup", func(db *gorm.DB) *gorm.DB { return db.Model(&Wx_assets_recom_group{}) }).
		Preload("Wxassetsrecommen", func(db *gorm.DB) *gorm.DB { return db.Model(&Wx_assets_recommen{}) }).
		Preload("Wxcompanies", func(db *gorm.DB) *gorm.DB { return db.Model(&Wx_companies{}) }).
		Where(queryParams).First(&returnData)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	return returnData
}

func (w clueassetsmodel) Delete(queryParams map[string]interface{}) (returnResult int64) {
	result := databases.Db
	if _, ok := queryParams["recomm_id"]; ok && queryParams["recomm_id"].([]int)[0] == 0 {
		result = result.Where("1=1")
		delete(queryParams, "recomm_id")
	}
	if _, ok := queryParams["ip"]; ok {
		result = result.Where("locate('" + queryParams["ip"].(string) + "', `ip`)")
		delete(queryParams, "ip")
	}
	result = result.Where(queryParams)
	result.Delete(&Wx_clue_assets{})
	return result.RowsAffected
}
