package Models

import (
	"database/sql/driver"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"wangxin2.0/app/Constant"
)

type BaseModel struct {
	table Model
}

type Model struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	CreatedAt *LocalTime `json:"created_at"`
	UpdatedAt *LocalTime `json:"updated_at"`
	DeletedAt *LocalTime `json:"deleted_at" sql:"index"`
}

type Pagination struct {
	Size int    `json:"size" form:"size" uri:"size"`
	Page int    `json:"page" form:"page" uri:"page"`
	Sort string `json:"sort" form:"sort" uri:"sort"`
}

func GeneratePaginationFromRequest(c *gin.Context) (pagination Pagination) {
	if err := c.ShouldBind(&pagination); err != nil {
		log.Print(Constant.PARAMS_BIND_ERROR)
		panic(Constant.PARAMS_BIND_ERROR)
	}
	// 校验参数
	if pagination.Size < 0 {
		pagination.Size = 2
	}
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if len(pagination.Sort) == 0 {
		pagination.Sort = "created_at desc"
	}
	return
}
func (b BaseModel) Init() {
	AssetModel = &assetmodel{}
	CompaniesModel = &companiesmodel{}
	ClueModel = &cluemodel{}
	ClueAssetsModel = &clueassetsmodel{}
	RecomGroupModel = &recomgroupmodel{}
	IpsModel = &ipsmodel{}
	AssetsChangeHistoryModel = &assetschangehistorymodel{}
	GeoatlaModel = &geoatlamodel{}
	UserModel = &usermodel{}
	TaskModel = &taskmodel{}
	PocModel = &pocmodel{}
	IndustryModel = &industrymodel{}
	ThreatsModel = &threatsmodel{}
	ClueSubDomainModel = &cluesubdomainmodel{}
	ScanAimsModel = &scanaimsmodel{}
	ScanSubdomainsModel = &scansubdomainsmodel{}
	wanganeeserviceModel = &wanganeeservicemodel{}
	wanganeesubdomainModel = &wanganeesubdomainmodel{}
	RulesModel = &rulesmodel{}
	SecondCategoryRulesModel = &secondcategoryrulesmodel{}
	SecondCategoriesModel = &secondcategoriesmodel{}
	CategoriesRulesModel = &categoriesrulesmodel{}
	CategoriesModel = &categoriesmodel{}
}

// 1. 创建 time.Time 类型的副本 LocalTime；
type LocalTime time.Time

const TimeFormat = "2006-01-02 15:04:05"

const IsoTimeFormat = "2006-01-02T15:04:05.00000000Z"

// 2. 为 LocalTime 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

// 3. 为 LocalTime 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
func (t *LocalTime) String() string {
	return time.Time(*t).Format(TimeFormat)
}
