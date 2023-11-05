package Models

import (
	"github.com/jinzhu/gorm"
	"wangxin2.0/databases"
)

var UserModel *usermodel

type usermodel struct {
	BaseModel
}

type Wx_user struct {
	Id       int
	Username string
	Password string
	Phone    string
	Email    string
	Adcode   string
	Status   int8
	Level    int16
}

func (w usermodel) FindUser(queryData map[string]interface{}) (returnResult map[string]interface{}, err error) {
	var returnData map[string]interface{}
	dbError := databases.Db.Model(&Wx_user{}).Where(queryData).First(&returnData).Error
	if dbError != nil && dbError != gorm.ErrRecordNotFound {
		return nil, dbError
	}
	return returnData, nil
}
