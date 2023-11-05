package Models

import (
	"gorm.io/gorm"
	"wangxin2.0/databases"
)

var CompaniesModel *companiesmodel

type companiesmodel struct {
	BaseModel
}
type Wx_companies struct {
	Id                              int         `json:"id" gorm:"column:id;primary_key"`
	User_id                         int         `json:"user_id" gorm:"column:user_id"`
	Industry_id                     int         `json:"industry_id" gorm:"column:industry_id"`
	Geoatla_id                      int         `json:"geoatla_id" gorm:"column:geoatla_id"`
	Geoatla_lft_id                  int         `json:"geoatla_lft_id" gorm:"column:geoatla_lft_id"`
	Name                            string      `json:"name" gorm:"column:name"`
	Area                            string      `json:"area" gorm:"column:area"`
	Risk_coefficient                float64     `json:"risk_coefficient" gorm:"column:risk_coefficient"`
	Longitude                       float64     `json:"longitude" gorm:"column:longitude"`
	Asset_count                     string      `json:"asset_count" gorm:"column:asset_count"`
	Ipv4_asset_count                string      `json:"ipv4_asset_count" gorm:"column:ipv4_asset_count"`
	Ipv6_asset_count                string      `json:"ipv6_asset_count" gorm:"column:ipv6_asset_count"`
	Service_asset_count             string      `json:"service_asset_count" gorm:"column:service_asset_count"`
	Loophole_count                  string      `json:"loophole_count" gorm:"column:loophole_count"`
	Ipv4_loophole_count             string      `json:"ipv4_loophole_count" gorm:"column:ipv4_loophole_count"`
	Ipv6_loophole_count             string      `json:"ipv6_loophole_count" gorm:"column:ipv6_loophole_count"`
	Announced_loophole_count        string      `json:"announced_loophole_count" gorm:"column:announced_loophole_count"`
	Confirm_repaired_loophole_count string      `json:"confirm_repaired_loophole_count" gorm:"column:confirm_repaired_loophole_count"`
	Not_repaired_loophole_count     string      `json:"not_repaired_loophole_count" gorm:"column:not_repaired_loophole_count"`
	Repaired_total_second           string      `json:"repaired_total_second" gorm:"column:repaired_total_second"`
	Announcement_count              string      `json:"announcement_count" gorm:"column:announcement_count"`
	Announced_at                    string      `json:"announced_at" gorm:"column:announced_at"`
	Repair_rate                     string      `json:"repair_rate" gorm:"column:repair_rate"`
	CreatedAt                       LocalTime   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                       LocalTime   `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt                       LocalTime   `json:"deleted_at" gorm:"column:deleted_at"`
	Industry                        Wx_industry `gorm:"foreignkey:industry_id"`
}

type Select_companies struct {
	Id   int    `json:"id" gorm:"column:id;primary_key"`
	Name string `json:"name" gorm:"column:name"`
}

func (w companiesmodel) Find(queryParams map[string]interface{}, whereInParams ...map[string]interface{}) (returnResult Wx_companies) {
	returnData := Wx_companies{}
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
