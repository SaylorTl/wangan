package Models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var TaskModel *taskmodel

type taskmodel struct {
	BaseModel
}

type Wx_tasks struct {
	Id                           uint64    `gorm:"column:id;primary_key;auto_increment" json:"id"`
	User_id                      uint64    `gorm:"column:user_id;not null" json:"user_id"`
	Name                         string    `gorm:"column:name;not null" json:"name"`
	Classification               int8      `gorm:"column:classification;not null;default:0" json:"classification"`
	Progress                     float64   `gorm:"column:progress;not null;default:0.00000" json:"progress"`
	Status                       int8      `gorm:"column:status;not null;default:0" json:"status"`
	Type                         int8      `gorm:"column:type;not null;default:0" json:"type"`
	Step                         int8      `gorm:"column:step;not null;default:0" json:"step"`
	ScanParam                    string    `gorm:"column:scan_param;not null" json:"scan_param"`
	Created_at                   LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                    LocalTime `gorm:"column:updated_at" json:"updated_at"`
	DispatchedAt                 LocalTime `gorm:"column:dispatched_at" json:"dispatched_at"`
	DeletedAt                    LocalTime `gorm:"column:deleted_at" json:"deleted_at" sql:"index"`
	Asset_scan_absolute_path     string    `json:"asset_scan_absolute_path" `
	Report_absolute_path         string    `json:"report_absolute_path"`
	Asset_scan_complete_filename string    `json:"asset_scan_complete_filename"`
	Poc_scan_absolute_path       string    `json:"poc_scan_absolute_path" `
}

func (t taskmodel) FindTaskById(taskId int) (res *Wx_tasks) {
	var wxTask Wx_tasks
	databases.Db.Where("id=?", taskId).Find(&wxTask)
	return &wxTask
}

func (t taskmodel) UpdateColumns(task *Wx_tasks, param map[string]interface{}) (res *Wx_tasks) {
	var wxTask Wx_tasks
	if err := databases.Db.Model(task).UpdateColumns(param).Error; err != nil {
		// 错误处理
		log.Panic(err)
	}
	return &wxTask
}

func (t taskmodel) Find(queryParams map[string]interface{}, whereInParams ...map[string]interface{}) (returnResult Wx_tasks) {
	returnData := Wx_tasks{}
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

func (w taskmodel) GetGoScannerExecutableCmd(task Wx_tasks, poc Wx_pocs, target string) string {
	var scanParam map[string]interface{}
	_ = json.Unmarshal([]byte(task.ScanParam), &scanParam)
	log.Printf("%#v", scanParam)
	bandwidth, _ := strconv.Atoi(scanParam["bandwidth"].(string))
	concurrent := bandwidth / 50
	if LOG4J2_FILENAME != poc.Filename && SPRING_CORE_FILENAME != poc.Filename {
		fileName := poc.Scan_aim_filename

		rPath := task.Poc_scan_absolute_path + fileName

		return fmt.Sprintf("%v -n %v -m %v -r %v 2>&1", Utils.F.GoScannerPath("goscanner"), concurrent, poc.Filename, rPath)
	}
	return fmt.Sprintf("%v -n %v -m %v -t %v 2>&1", Utils.F.GoScannerPath("goscanner"), concurrent, poc.Filename, target)
}

func (w *Wx_tasks) AfterFind(tx *gorm.DB) (err error) {
	w.getPocScanAbsolutePathAttribute()
	w.getAssetScanAbsolutePathAttribute()
	w.getAssetScanCompleteFilenameAttribute()
	w.getReportAbsolutePathAttribute()
	return
}

func (w *Wx_tasks) getPocScanAbsolutePathAttribute() {
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", w.Created_at.String(), time.Local)
	create_at := time.Unix(stamp.Unix(), 0).Format("200601/02")
	w.Poc_scan_absolute_path = Constant.wangxinAbsoulePath + "/storage/" + fmt.Sprintf(Constant.LOOPHOLE_SCAN_AIM_FILE_PATH, w.User_id, create_at, w.Id)
}

func (w *Wx_tasks) getAssetScanAbsolutePathAttribute() {
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", w.Created_at.String(), time.Local)
	create_at := time.Unix(stamp.Unix(), 0).Format("200601/02")
	w.Report_absolute_path = Constant.wangxinAbsoulePath + "/storage/" + fmt.Sprintf(Constant.ASSET_SCAN_AIM_FILE_PATH, w.User_id, create_at, w.Id)
}

func (w *Wx_tasks) getAssetScanCompleteFilenameAttribute() {
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", w.Created_at.String(), time.Local)
	create_at := time.Unix(stamp.Unix(), 0).Format("200601/02")
	w.Asset_scan_complete_filename = Constant.wangxinAbsoulePath + "/storage/" + fmt.Sprintf(Constant.ASSET_SCAN_AIM_FILE_ABSOLUTE_PATH, w.User_id, create_at, w.Id)
}

func (w *Wx_tasks) getReportAbsolutePathAttribute() {
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", w.Created_at.String(), time.Local)
	create_at := time.Unix(stamp.Unix(), 0).Format("200601/02")
	w.Report_absolute_path = Constant.wangxinAbsoulePath + "/storage/" + fmt.Sprintf(Constant.REPORT_PATH, create_at, w.Id)
}
