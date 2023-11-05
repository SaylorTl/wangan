package Commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

func init() {
	rootCmd.AddCommand(jsonimportassetCmd)
}

var jsonimportassetCmd = &cobra.Command{
	Use:   "jsonimportasset",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		var files []string
		root := Constant.SubdomainFilePath
		err := filepath.Walk(root, func(pathStr string, info os.FileInfo, err error) error {
			if true == info.IsDir() {
				return nil
			}
			fileSuffix := path.Ext(pathStr)
			if ".json" != fileSuffix {
				return nil

			}
			files = append(files, pathStr)
			return nil
		})
		if nil != err {
			log.Print(err)
		}
		for _, filepath := range files {
			subdomainstream := Utils.NewJSONStream()
			//资产入库
			go func() {
				domaini := 0
				servicei := 0
				insertDomainArr := make([]map[string]interface{}, 0)
				insertServiceArr := make([]map[string]interface{}, 0)
				for data := range subdomainstream.Watch() {
					if nil != data.ReturnData && "wanganpro_subdomain" == data.ReturnData["_index"].(string) {
						domaini++
						insertDomainArr = append(insertDomainArr, data.ReturnData)
					}
					if nil != data.ReturnData && "wanganpro_service" == data.ReturnData["_index"].(string) {
						servicei++
						insertServiceArr = append(insertServiceArr, data.ReturnData)
					}
					if data.Error != nil {
						Utils.F.SetProcess(float32(50))
						if len(insertDomainArr) > 0 {
							Utils.ImportAssetPool.Submit(Models.wanganeesubdomainModel.InsertwanganSubdomainAsset(insertDomainArr))
						}
						if len(insertServiceArr) > 0 {
							Utils.ImportAssetPool.Submit(Models.wanganeeserviceModel.InsertwanganServiceAsset(insertServiceArr))
						}
						time.Sleep(5 * time.Second)
						databases.Redis.Del("wangansync:is_task_running")
						log.Print("subdomain error ", data.Error)
						break
					}
					Utils.F.SetProcess(data.Process)
					if domaini >= 250 {
						Utils.ImportAssetPool.Submit(Models.wanganeesubdomainModel.InsertwanganSubdomainAsset(insertDomainArr))
						domaini = 0
						insertDomainArr = make([]map[string]interface{}, 0)
					}
					if servicei >= 250 {
						Utils.ImportAssetPool.Submit(Models.wanganeeserviceModel.InsertwanganServiceAsset(insertServiceArr))
						servicei = 0
						insertServiceArr = make([]map[string]interface{}, 0)
					}
				}
			}()
			log.Print("reading file start", filepath)
			subdomainstream.Start(filepath)
			log.Print("reading file end", filepath)
			databases.Redis.Set("wangansync:is_task_running", 1, 172800*time.Second)
			for {
				is_task_running, _ := databases.Redis.Get("wangansync:is_task_running").Result()
				if "" == is_task_running {
					break
				}
			}
		}
		//资产聚合
		conf := databases.Conf
		ctx := context.Background()
		SCAN_INDEX := conf.Section("es").Key("ES_wangxin_SERVICE_INDEX").String()
		SCAN_SUBDOMAIN_INDEX := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_INDEX").String()
		//资产聚合查询
		searchResult, _ := databases.Esdb.Es.Scroll().
			Index(SCAN_INDEX, SCAN_SUBDOMAIN_INDEX).
			FetchSourceContext(elastic.NewFetchSourceContext(true).Include("ip")).
			Query(elastic.NewMatchAllQuery()).
			Scroll("45m").
			Size(1000).
			Pretty(true).
			Do(ctx) // 执行请求
		totalHits := searchResult.Hits.TotalHits
		aggpage := 1
		userQueryData := make(map[string]interface{})
		userQueryData["username"] = "admin"
		userData, _ := Models.UserModel.FindUser(userQueryData)
		user_id := userData["id"].(int)
		company_id := 0
		scroll_id := searchResult.ScrollId
		currenti := 0
		for {
			temArr := make([]interface{}, 0)
			for _, val := range searchResult.Hits.Hits {
				currenti++
				tmpMap := make(map[string]interface{})
				json.Unmarshal(*val.Source, &tmpMap)
				temArr = append(temArr, tmpMap["ip"])
			}
			ipArr, _ := Utils.F.Unique(temArr)
			for _, ipval := range ipArr.([]interface{}) {
				Utils.ImportSyncAssetPool.Submit(Models.AssetModel.InsertAsset(ipval.(string), user_id, company_id))
			}
			strprogress, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (float32(currenti)/float32(totalHits*2))*100), 32)
			currentProcess := float32(50) + float32(strprogress)
			Utils.F.SetProcess(currentProcess)

			aggpage++
			scroll_id = searchResult.ScrollId
			scrollData, _ := databases.Esdb.Es.Scroll("45m").ScrollId(searchResult.ScrollId).Do(ctx)
			log.Print("insert asset process", aggpage)
			searchResult = scrollData
			if nil == searchResult || 0 == len(searchResult.Hits.Hits) {
				log.Print("insert asset process", aggpage)
				break
			}
		}
		databases.Esdb.Es.ClearScroll(scroll_id)
		for {
			runningNums := Utils.ImportSyncAssetPool.Running()
			if 0 == runningNums {
				break
			}
		}
	},
}
