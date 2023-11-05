package Constant

const ASSET_SCAN_AIM_FILE_ABSOLUTE_PATH string = "app/scan/%v/%v/asset_task/%v/ip_range.csv" //资产扫描目标文件路径
const ASSET_SCAN_AIM_FILE_NAME string = "ip_range.csv"                                       //资产扫描目标文件名
const ASSET_SCAN_AIM_FILE_PATH string = "app/scan/%v/%v/asset_task/%v/"                      //资产扫描目标文件路径
const LOOPHOLE_SCAN_AIM_FILE_PATH string = "app/scan/%v/%v/loophole_task/%v/"                //漏洞扫描目标文件夹路径
const TASK_PAUSH_CONF_FILE_PATH string = "app/pause/task/%v/%v/task.conf"                    //漏洞扫描目标文件夹路径
const TASK_PAUSH_CONF_PATH string = "app/pause/task/%v/%v/"                                  //漏洞扫描目标文件夹路径
const REPORT_PATH string = "storage/reports/%v/%v"

const SOURCE_IP_UPLOAD int = 0         //页面上传
const SOURCE_PAGE_INPUT_ADD int = 1    //页面添加
const SOURCE_UNKNOWN_ASSETS int = 2    //未知资产
const SOURCE_wanganSYNC_ASSETS int = 3 //wangan同步资产
const SOURCE_FORADAR_DETECT int = 4    //单位资产探测
const SOURCE_DOMAIN_DETECT int = 5     //域名资产资产探测
