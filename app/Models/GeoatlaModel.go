package Models

import (
	"gorm.io/gorm"
	"wangxin2.0/databases"
)

var GeoatlaModel *geoatlamodel

type geoatlamodel struct {
	BaseModel
}

type Wx_geoatlas struct {
	Id        int       `json:"id" gorm:"column:id"`
	Parent_id int       `json:"parent_id" gorm:"column:parent_id"`
	Lft       int       `json:"lft" gorm:"column:lft"`
	Rgt       int       `json:"rgt" gorm:"column:rgt"`
	Depth     int16     `json:"depth" gorm:"column:depth"`
	Name      string    `json:"name" gorm:"column:name"`
	Level     int       `json:"level" gorm:"column:level"`
	Adcode    int       `json:"adcode" gorm:"column:adcode"`
	Initials  string    `json:"initials" gorm:"column:initials"`
	Pinyin    string    `json:"pinyin" gorm:"column:pinyin"`
	CreatedAt LocalTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt LocalTime `json:"updated_at" gorm:"column:updated_at"`
}

func (w geoatlamodel) Find(queryParams map[string]interface{}, whereInParams ...map[string]interface{}) (returnResult Wx_geoatlas) {
	returnData := Wx_geoatlas{}
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

func (w geoatlamodel) Findgeotlas() (*Wx_geoatlas, error) {
	var wxGeoatlas Wx_geoatlas
	err := databases.Db.First(&wxGeoatlas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &wxGeoatlas, nil
}

func (w geoatlamodel) Addgeotlas(Wxgeoatlas *Wx_geoatlas) (id int, err error) {
	result := databases.Db.Where("name=?", Wxgeoatlas.Name).First(&Wxgeoatlas)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		result = databases.Db.Create(Wxgeoatlas)
	}
	return Wxgeoatlas.Id, result.Error
}
func (w geoatlamodel) Updatageotlas(id int, data Wx_geoatlas) (num int64, err error) {
	var wxGeoatlas Wx_geoatlas
	res := databases.Db.First(&wxGeoatlas, id).Updates(data)
	return res.RowsAffected, res.Error
}
