package Services

import (
	"encoding/json"
	"fmt"
	"github.com/thinkeridea/go-extend/exnet"
	"log"
	"os/exec"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var ClueSubdomainService cluesubdomainservice

type cluesubdomainservice struct {
}

func (a cluesubdomainservice) Extract(user_id int) (bool, string) {
	//构建异步队列，防并发场景
	params := make(map[string]interface{})
	params["status"] = 3
	rangeData := Models.ClueSubDomainModel.Getall(params)
	for _, data := range rangeData {
		queryParams := map[string]interface{}{}
		queryParams["domain"] = data.Domain
		queryParams["ip"] = data.Ip
		returnData := Models.ScanSubdomainsModel.Show(queryParams)
		//未知资产发异步发现
		SubdomainParams := make(map[string]interface{})
		SubdomainParams["subdomain_clue_id"] = data.Subdomain_clue_id
		SubdomainParams["domain"] = data.Domain
		SubdomainParams["path"] = data.Path
		SubdomainParams["user_id"] = user_id
		SubdomainParams["company_id"] = 0
		if company_id, ok := returnData["company_id"]; ok {
			SubdomainParams["company_id"] = int(company_id.(float64))
		}
		databases.Redis.RPush("subdomain:clue", Utils.F.MapToJson(SubdomainParams))
		SubdomainUpdataParams := make(map[string]interface{})
		SubdomainUpdataParams["status"] = 0
		SubdomainUpdataParams["subdomain_clue_id"] = data.Subdomain_clue_id
		Models.ClueSubDomainModel.UpdateOrCreate(SubdomainUpdataParams)
	}
	go a.createClueDomainPool()
	return true, Constant.EXTRACT_SUCCESS
}

// 创建提炼池
func (a cluesubdomainservice) createClueDomainPool() {
	err := Utils.SubdomainPool.Submit(a.ExtractClue())
	log.Print(err)
}

type cluDomainFunc func()

// 创建提炼池
func (a cluesubdomainservice) ExtractClue() cluDomainFunc {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Print("subdomainError", r)
				runningParams := make(map[string]interface{})
				runningParams["status"] = 1
				runnningArr := Models.ClueSubDomainModel.Getall(runningParams)
				for _, runningdata := range runnningArr {
					updateData := make(map[string]interface{})
					updateData["subdomain_clue_id"] = runningdata.Subdomain_clue_id
					updateData["status"] = 4
					_, ferr := Models.ClueSubDomainModel.UpdateOrCreate(updateData)
					if ferr != nil {
						log.Print(ferr)
					}
				}
			}
		}()
		for {
			Subdomainstring, wanganerror := databases.Redis.LPop("subdomain:clue").Result()
			log.Print("当前待扫描数据", Subdomainstring)
			if nil != wanganerror {
				log.Print(wanganerror)
				break
			}
			// 写数据
			Subdomain := make(map[string]interface{})
			_ = json.Unmarshal([]byte(Subdomainstring), &Subdomain)
			log.Print("当前待扫描数据", Subdomain)

			//删除对照数组中元素
			showData := Models.ClueSubDomainModel.Show(map[string]interface{}{"subdomain_clue_id": int(Subdomain["subdomain_clue_id"].(float64))})
			if 0 == len(showData) {
				continue
			}

			//将域名状态设置为探测中
			updateData := make(map[string]interface{})
			updateData["subdomain_clue_id"] = int(Subdomain["subdomain_clue_id"].(float64))
			updateData["status"] = 1
			_, uerr := Models.ClueSubDomainModel.UpdateOrCreate(Utils.F.C(updateData))
			if uerr != nil {
				log.Print("已删除记录跳过！", updateData)
				fmt.Println(uerr)
				continue
			}
			goScannerExecutableCmd := "./sudomy -d " + Subdomain["domain"].(string) + "  -pS -rS"
			cmd := exec.Command("bash", "-c", goScannerExecutableCmd)
			log.Print("正在探测中", goScannerExecutableCmd)
			cmd.Dir = "/wangxin/Sudomy"
			cmdErr := cmd.Run()
			if cmdErr != nil {
				log.Print("sudomy fail", cmdErr)
				continue
			}
			data, _ := Utils.F.ReadTxt(Subdomain["path"].(string) + "/Subdomain_Resolver.txt")
			livedata, _ := Utils.F.ReadTxt(Subdomain["path"].(string) + "/Live_hosts_pingsweep.txt")
			liveDomains := []string{}
			for _, livevalue := range livedata {
				liveDomains = append(liveDomains, livevalue[0].(string))
			}
			Subdomain_resolve := []map[string]interface{}{}
			for _, value := range data {
				if !Utils.F.InArray(value[0].(string), liveDomains) {
					continue
				}
				tmp := map[string]interface{}{}
				tmp["domain"] = value[0]
				tmp["ip"] = value[1]
				tmp["subdomain_clue_id"] = int(Subdomain["subdomain_clue_id"].(float64))
				tmp["ip_segment"] = value[1].(string)
				tmp["ip_type"] = Constant.IP_TYPE_V4
				tmp["createtime"] = time.Now().Format("2006-01-02 15:04:05")
				tmp["lastupdatetime"] = time.Now().Format("2006-01-02 15:04:05")
				tmp["belong_user_id"] = int(Subdomain["user_id"].(float64))
				start_ip, _ := exnet.IPString2Long(value[1].(string))
				tmp["start_ip"] = int64(start_ip)
				tmp["end_ip"] = int64(start_ip)
				tmp["ip_source"] = Constant.SOURCE_DOMAIN_DETECT
				tmp["company_id"] = int(Subdomain["company_id"].(float64))
				tmp["domain"] = value[0].(string)
				tmp["user_id"] = int(Subdomain["user_id"].(float64))
				if 0 != tmp["company_id"] {
					Models.IpsModel.UpdateOrCreate(Utils.F.C(tmp))
				}
				Subdomain_resolve = append(Subdomain_resolve, tmp)
			}
			log.Print("Subdomain_resolve", Subdomain_resolve)
			Models.ScanAimsModel.InsertScanAims(Subdomain_resolve)
			Models.ScanSubdomainsModel.InsertScanSubdomains(Subdomain_resolve)
			//将域名状态设置为探测完成
			updateData["status"] = 2
			_, ferr := Models.ClueSubDomainModel.UpdateOrCreate(Utils.F.C(updateData))
			if ferr != nil {
				continue
			}
		}
		return
	}
}
