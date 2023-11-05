package Models

import (
	"gorm.io/gorm"
	"wangxin2.0/databases"
)

var IpsModel *ipsmodel

type ipsmodel struct {
	BaseModel
}

type Wx_ips struct {
	Id         int
	Ip         string
	User_id    int
	Company_id int
	Status     int
	Ip_source  int
	Ip_type    int
}

func (w ipsmodel) Findip(parseIp string) (res *Wx_ips) {
	var Wxips Wx_ips
	databases.Db.Where("ip=?", parseIp).Where("deleted_at IS NULL").First(&Wxips)
	return &Wxips
}

func (w ipsmodel) UpdateOrCreate(params map[string]interface{}) (int, error) {
	var Wxips Wx_ips
	result := databases.Db.Where("company_id=? and ip=?", params["company_id"].(int), params["ip"].(string)).First(&Wxips)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 { //没有则插入
		Wxips.Ip = params["ip"].(string)
		Wxips.User_id = params["user_id"].(int)
		Wxips.Company_id = params["company_id"].(int)
		Wxips.Ip_type = params["ip_type"].(int)
		Wxips.Ip_source = params["ip_source"].(int)
		insertRes := databases.Db.Model(Wx_ips{}).Create(&Wxips)
		return Wxips.Id, insertRes.Error
	}
	//有则更新
	updataRes := databases.Db.Model(&Wxips).Updates(params)
	return Wxips.Id, updataRes.Error
}
