package Models

import (
	"gorm.io/gorm"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/databases"
)

var ClueSubDomainModel *cluesubdomainmodel

type cluesubdomainmodel struct {
	BaseModel
}

type Wx_subdomain_clue_list struct {
	Subdomain_clue_id int       `json:"subdomain_clue_id" gorm:"column:subdomain_clue_id;primary_key"`
	Domain            string    `json:"domain" gorm:"column:domain"`
	Ip                string    `json:"ip" gorm:"column:ip"`
	Status            int       `json:"status" gorm:"column:status"`
	Path              string    `json:"path" gorm:"column:path"`
	CreatedAt         LocalTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         LocalTime `json:"updated_at" gorm:"column:updated_at"`
}

func (w cluesubdomainmodel) Lists(queryParams map[string]interface{}, pagination *Pagination) (returnResult []Wx_subdomain_clue_list, countResult int64) {
	// 分页查询
	offset := (pagination.Page - 1) * pagination.Size
	returnData := []Wx_subdomain_clue_list{}
	delete(queryParams, "page")
	delete(queryParams, "size")
	returnDB := databases.Db.Model(&Wx_subdomain_clue_list{})
	if _, ok := queryParams["domain"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["domain"].(string) + "', `domain`)")
		delete(queryParams, "domain")
	}
	if _, ok := queryParams["keyword"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["keyword"].(string) + "', `domain`)")
		delete(queryParams, "keyword")
	}
	if _, ok := queryParams["status"]; ok {
		returnDB.Where("status= ?", queryParams["status"])
		delete(queryParams, "status")
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
	returnDB.Where(queryParams).Limit(pagination.Size).Offset(offset).Order("FIELD(status, 1, 0, 2,3),created_at DESC").Find(&returnData)
	if returnDB.Error != nil && returnDB.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	var count int64
	returnDB.Offset(-1).Limit(-1).Count(&count)
	return returnData, count

}

func (w cluesubdomainmodel) Count(queryParams map[string]interface{}) int64 {
	// 分页查询
	returnDB := databases.Db.Model(&Wx_subdomain_clue_list{})
	if _, ok := queryParams["domain"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["domain"].(string) + "', `domain`)")
		delete(queryParams, "domain")
	}
	if _, ok := queryParams["keyword"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["keyword"].(string) + "', `domain`)")
		delete(queryParams, "keyword")
	}
	if _, ok := queryParams["status"]; ok {
		returnDB.Where("status= ?", queryParams["status"])
		delete(queryParams, "status")
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
	var count int64
	returnDB.Where(queryParams).Count(&count)
	if returnDB.Error != nil && returnDB.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	return count

}

func (w cluesubdomainmodel) Getall(queryParams map[string]interface{}) (returnResult []Wx_subdomain_clue_list) {
	// 分页查询
	returnData := []Wx_subdomain_clue_list{}
	returnDB := databases.Db.Model(&Wx_subdomain_clue_list{})
	if _, ok := queryParams["domain"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["domain"].(string) + "', `domain`)")
		delete(queryParams, "domain")
	}
	if _, ok := queryParams["keyword"]; ok {
		returnDB = returnDB.Where("locate('" + queryParams["keyword"].(string) + "', `domain`)")
		delete(queryParams, "keyword")
	}
	if _, ok := queryParams["status"]; ok {
		returnDB.Where("status= ?", queryParams["status"])
		delete(queryParams, "status")
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
	returnDB.Where(queryParams).Find(&returnData)
	if returnDB.Error != nil && returnDB.Error != gorm.ErrRecordNotFound {
		log.Print(Constant.DATABASE_QUERY_ERR)
		panic(Constant.DATABASE_QUERY_ERR)
	}
	return returnData

}

func (w cluesubdomainmodel) Show(queryParams map[string]interface{}, whereInParams ...map[string]interface{}) (returnResult []Wx_subdomain_clue_list) {
	returnData := []Wx_subdomain_clue_list{}
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

func (w cluesubdomainmodel) Add(params map[string]interface{}) (int, error) {
	var Wxsubdomaincluelist Wx_subdomain_clue_list
	result := databases.Db.Select("subdomain_clue_id").Where("domain=? and path=?", params["domain"].(int), params["path"].(string)).First(&Wxsubdomaincluelist)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxsubdomaincluelist.Domain = params["domain"].(string)
		Wxsubdomaincluelist.Path = params["path"].(string)
		Wxsubdomaincluelist.Status = params["status"].(int)
		insertRes := databases.Db.Model(Wx_subdomain_clue_list{}).Create(&Wxsubdomaincluelist)
		return Wxsubdomaincluelist.Subdomain_clue_id, insertRes.Error
	}
	//有则更新
	return Wxsubdomaincluelist.Subdomain_clue_id, nil
}

func (w cluesubdomainmodel) Delete(queryParams map[string]interface{}) (returnResult int64) {
	result := databases.Db
	if _, ok := queryParams["subdomain_clue_id"]; ok && queryParams["subdomain_clue_id"].([]int)[0] == 0 {
		result = result.Where("1=1")
		delete(queryParams, "subdomain_clue_id")
	}
	if _, ok := queryParams["subdomain_clue_ids"]; ok {
		result = result.Where("subdomain_clue_id IN (?)", queryParams["subdomain_clue_ids"])
		delete(queryParams, "subdomain_clue_ids")
	}
	if _, ok := queryParams["domain"]; ok {
		result = result.Where("locate('" + queryParams["domain"].(string) + "', `domain`)")
		delete(queryParams, "domain")
	}
	if _, ok := queryParams["keyword"]; ok {
		result = result.Where("locate('" + queryParams["keyword"].(string) + "', `domain`)")
		delete(queryParams, "keyword")
	}
	if _, ok := queryParams["status"]; ok {
		result = result.Where("status= ?", queryParams["status"])
		delete(queryParams, "status")
	}
	result = result.Where(queryParams)
	result.Delete(&Wx_subdomain_clue_list{})
	return result.RowsAffected
}

func (w cluesubdomainmodel) UpdateOrCreate(params map[string]interface{}) (int, error) {
	var Wxsubdomaincluelist Wx_subdomain_clue_list
	result := databases.Db
	if _, ok := params["subdomain_clue_id"]; ok {
		result = result.Where("subdomain_clue_id=?", params["subdomain_clue_id"].(int))
	}
	if _, ok := params["domain"]; ok {
		result = result.Where("domain=?", params["domain"].(string))
	}
	if _, ok := params["path"]; ok {
		result = result.Where("path=?", params["path"].(string))
	}
	result.First(&Wxsubdomaincluelist)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxsubdomaincluelist.Domain = params["domain"].(string)
		Wxsubdomaincluelist.Path = params["path"].(string)
		Wxsubdomaincluelist.Ip = params["ip"].(string)
		Wxsubdomaincluelist.Status = 0
		if _, ok := params["status"]; ok {
			Wxsubdomaincluelist.Status = params["status"].(int)
		}
		insertRes := databases.Db.Model(Wx_subdomain_clue_list{}).Create(&Wxsubdomaincluelist)
		return Wxsubdomaincluelist.Subdomain_clue_id, insertRes.Error
	}
	//有则更新
	updataRes := databases.Db.Model(&Wxsubdomaincluelist).Updates(params)
	return Wxsubdomaincluelist.Subdomain_clue_id, updataRes.Error
}

func (w cluesubdomainmodel) TxUpdateOrCreate(tx *gorm.DB, params map[string]interface{}) (int, error) {
	var Wxsubdomaincluelist Wx_subdomain_clue_list
	result := databases.Db.Where("subdomain_clue_id=?", params["subdomain_clue_id"].(int)).First(&Wxsubdomaincluelist)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxsubdomaincluelist.Domain = params["domain"].(string)
		Wxsubdomaincluelist.Path = params["path"].(string)
		Wxsubdomaincluelist.Status = 0
		insertRes := tx.Model(Wx_subdomain_clue_list{}).Create(&Wxsubdomaincluelist)
		return Wxsubdomaincluelist.Subdomain_clue_id, insertRes.Error
	}
	//有则更新
	updataRes := tx.Model(&Wxsubdomaincluelist).Updates(params)
	return Wxsubdomaincluelist.Subdomain_clue_id, updataRes.Error
}
