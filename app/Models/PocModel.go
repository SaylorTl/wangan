package Models

import (
	"fmt"
	"gorm.io/gorm"
	"wangxin2.0/databases"
)

var PocModel *pocmodel

type pocmodel struct {
	BaseModel
}

type Wx_pocs struct {
	Id                int       `json:"id" gorm:"column:id;primary_key"`
	Filename          string    `json:"filename" gorm:"column:filename"`
	Name              string    `json:"name" gorm:"column:name"`
	Author            string    `json:"author" gorm:"column:author"`
	Level             int       `json:"level" gorm:"column:level"`
	Description       string    `json:"description" gorm:"column:description"`
	Impact            string    `json:"impact" gorm:"column:impact"`
	Recommendation    string    `json:"recommendation" gorm:"column:recommendation"`
	wangan_query      string    `json:"wangan_query" gorm:"column:wangan_query"`
	Has_exp           int       `json:"has_exp" gorm:"column:has_exp"`
	Homepage          string    `json:"homepage" gorm:"column:homepage"`
	Product           string    `json:"product" gorm:"column:product"`
	Disclosure_date   string    `json:"disclosure_date" gorm:"column:disclosure_date"`
	Cvss_score        string    `json:"cvss_score" gorm:"column:cvss_score"`
	Goby_query        string    `json:"goby_query" gorm:"column:goby_query"`
	Exp_tips          string    `json:"exp_tips" gorm:"column:exp_tips"`
	Exp_params        string    `json:"exp_params" gorm:"column:exp_params"`
	Cve_ids           string    `json:"cve_ids" gorm:"column:cve_ids"`
	References        string    `json:"references" gorm:"column:references"`
	Vul_type          string    `json:"vul_type" gorm:"column:vul_type"`
	Attack_surfaces   string    `json:"attack_surfaces" gorm:"column:attack_surfaces"`
	Status            string    `json:"status" gorm:"column:status"`
	Deleted_at        string    `json:"deleted_at" gorm:"column:deleted_at"`
	CreatedAt         LocalTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         LocalTime `json:"updated_at" gorm:"column:updated_at"`
	Scan_aim_filename string    `json:"scan_aim_filename" `
}

func (p pocmodel) FindPocById(pocId int) (res *Wx_pocs) {
	var wxPoc Wx_pocs
	databases.Db.Where("id=?", pocId).Find(&wxPoc)
	return &wxPoc
}
func (p pocmodel) FindPocByName(fileName string) (res *Wx_pocs) {
	var wxPoc Wx_pocs
	databases.Db.Where("filename=?", fileName).Find(&wxPoc)
	return &wxPoc
}

const LOG4J2_FILENAME string = "CVD-2021-15875.json"
const SPRING_CORE_FILENAME string = "CVD-2022-1090.json"

func (w pocmodel) Find(queryParams map[string]interface{}, whereInParams ...map[string]interface{}) (returnResult Wx_pocs) {
	returnData := Wx_pocs{}
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

func (w *Wx_pocs) getScanAimFilenameAttribute() {
	w.Scan_aim_filename = fmt.Sprintf("poc_%v.csv", w.Filename)
}
func (w *Wx_pocs) AfterFind(tx *gorm.DB) (err error) {
	w.getScanAimFilenameAttribute()
	return
}
