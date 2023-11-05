package Controllers

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"github.com/thinkeridea/go-extend/exnet"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Http/Request"
	"wangxin2.0/app/Models"
	"wangxin2.0/app/Services"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

const (
	email          = ""
	key            = ""
	field          = ""
	field_stats    = "
	base_url       = "https://wangan.info/api/v1/search/all"
	base_url_stats = "https://wangan.info/api/v1/search/stats"
	next_url       = "https://wangan.info/api/v1/search/next"
)

type wangansearchController struct {
	BaseController
}
type wangansearchShowValidate struct {
	Keyword    string `form:"keyword" to:"keyword" zh_cn:"关键词" binding:"required"`
	Page       int    `form:"page" to:"page" zh_cn:"页数" binding:"omitempty,min=1"`
	Size       int    `form:"size" to:"size" zh_cn:"页码" binding:"omitempty,min=10,max=500"`
	Full       *bool  `form:"full" to:"full" binding:"required"`
	Company_id int    `form:"company_id" to:"company_id"`
}
type wangansearchStatsValidate struct {
	Keyword string `form:"keyword" to:"keyword" zh_cn:"关键词" binding:"required"`
}
type wangansearchImport struct {
	Ip               string `json:"ip"`
	Port             string `json:"port"`
	Protocol         string `json:"protocol"`
	Country          string `json:"country"`
	City             string `json:"city"`
	Title            string `json:"title"`
	Banner           string `json:"banner"`
	Body             string `json:"body"`
	Server           string `json:"server"`
	Domain           string `json:"domain"`
	Country_name     string `json:"country_name"`
	Region           string `json:"region"`
	Longitude        string `json:"longitude"`
	Latitude         string `json:"latitude"`
	Product          string `json:"product"`
	Product_category string `json:"product_category"`
	Header           string `json:"header"`
	Fid              string `json:"fid"`
	Structinfo       string `json:"structinfo"`
	Icp              string `json:"icp"`
	Os               string `json:"os"`
}

func (f wangansearchController) Lists(c *gin.Context) {
	var params wangansearchShowValidate
	if err := c.ShouldBindWith(&params, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	qbase64 := base64.StdEncoding.EncodeToString([]byte(params.Keyword))
	queryParams := url.Values{}
	Url, _ := url.Parse(base_url)
	queryParams.Set("email", email)
	queryParams.Set("key", key)
	queryParams.Set("qbase64", qbase64)
	queryParams.Set("fields", field)
	queryParams.Set("page", strconv.Itoa(params.Page))
	queryParams.Set("size", strconv.Itoa(params.Size))
	queryParams.Set("full", strconv.FormatBool(*params.Full))
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = queryParams.Encode()
	result, _ := url.QueryUnescape(Url.String())

	// 忽略 https 证书验证
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Get(result)
	if err != nil {
		log.Println("wangansearch get error", err)
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	data := make(map[string]interface{})
	_ = json.Unmarshal(body, &data)
	Constant.Success(c, "成功", data)
	return
}

func (f wangansearchController) Stats(c *gin.Context) {
	var params wangansearchStatsValidate
	if err := c.ShouldBindWith(&params, binding.Query); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	qbase64 := base64.StdEncoding.EncodeToString([]byte(params.Keyword))
	queryParams := url.Values{}
	Url, _ := url.Parse(base_url_stats)
	queryParams.Set("email", email)
	queryParams.Set("key", key)
	queryParams.Set("qbase64", qbase64)
	queryParams.Set("fields", field_stats)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = queryParams.Encode()
	result, _ := url.QueryUnescape(Url.String())
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Get(result)
	if err != nil {
		log.Println("wangansearchstats get error", err)
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	data := make(map[string]interface{})
	_ = json.Unmarshal(body, &data)
	Constant.Success(c, "成功", data)
	return
}

func (f wangansearchController) Import(c *gin.Context) {
	var params wangansearchShowValidate
	if err := c.ShouldBindWith(&params, binding.Form); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	wanganvalue, wanganerror := databases.Redis.Get("wangansearch:totalsize").Result()
	if "" != wanganvalue || wanganerror != redis.Nil {
		// 查询不到数据
		Constant.Error(c, Constant.wangan_SYNC_TASK_EXISTS)
		return
	}
	queryParams := map[string]interface{}{}
	queryParams["user_id"] = f.GetUserId(c)
	queryParams["keyword"] = params.Keyword
	queryParams["base_url"] = next_url
	queryParams["email"] = email
	queryParams["key"] = key
	queryParams["fields"] = field
	queryParams["full"] = *params.Full
	queryParams["company_id"] = params.Company_id
	go Services.wangansearchService.Find(Utils.F.C(queryParams))
	Constant.Success(c, "成功", "")
	return
}

func (f wangansearchController) SyncFileData(c *gin.Context) {
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
	go func() {
		for _, filepath := range files {
			subdomainstream := Utils.NewJSONStream()
			for {
				is_task_running, _ := databases.Redis.Get("wangansync:is_task_running").Result()
				if "" == is_task_running {
					log.Print("break")
					break
				}
			}
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
			subdomainstream.Start(filepath)
			databases.Redis.Set("wangansync:is_task_running", 1, 172800*time.Second)
		}
	}()
	Constant.Success(c, "运行成功", "")
	return
}

func (f wangansearchController) AggregationAsset(c *gin.Context) {
	//计算聚合总数
	go func() {
		conf := databases.Conf
		ctx := context.Background()
		SCAN_INDEX := conf.Section("es").Key("ES_wangxin_SERVICE_INDEX").String()
		SCAN_SUBDOMAIN_INDEX := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_INDEX").String()
		//资产聚合查询
		searchResult, _ := databases.Esdb.Es.Scroll().
			Index(SCAN_INDEX, SCAN_SUBDOMAIN_INDEX).
			FetchSourceContext(elastic.NewFetchSourceContext(true).Include("ip")).
			Query(elastic.NewMatchAllQuery()).
			Scroll("60m").
			Size(2000).
			Pretty(true).
			Do(ctx) // 执行请求
		aggpage := 1
		user_id := f.GetUserId(c)
		company_id := 0
		scroll_id := searchResult.ScrollId
		for {
			temArr := make([]interface{}, 0)
			for _, val := range searchResult.Hits.Hits {
				tmpMap := make(map[string]interface{})
				json.Unmarshal(*val.Source, &tmpMap)
				temArr = append(temArr, tmpMap["ip"])
			}
			ipArr, _ := Utils.F.Unique(temArr)
			for _, ipval := range ipArr.([]interface{}) {
				Utils.ImportSyncAssetPool.Submit(Models.AssetModel.InsertAsset(ipval.(string), user_id, company_id))
			}
			aggpage++
			scroll_id = searchResult.ScrollId
			scrollData, _ := databases.Esdb.Es.Scroll("60m").ScrollId(searchResult.ScrollId).Do(ctx)
			log.Print("insert asset process", aggpage)
			searchResult = scrollData
			if nil == searchResult || 0 == len(searchResult.Hits.Hits) {
				break
			}
		}
		databases.Esdb.Es.ClearScroll(scroll_id)
	}()
	Constant.Success(c, "运行成功", "")
	return
}

func (f wangansearchController) ClearSubdomain(c *gin.Context) {
	//计算聚合总数
	go func() {
		conf := databases.Conf
		ctx := context.Background()
		SCAN_SUBDOMAIN_INDEX := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_INDEX").String()
		SCAN_SUBDOMAIN_TYPE := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_TYPE").String()
		//资产聚合查询
		searchResult, _ := databases.Esdb.Es.Scroll().
			Index(SCAN_SUBDOMAIN_INDEX).
			FetchSourceContext(elastic.NewFetchSourceContext(true).Include("host")).
			Query(elastic.NewMatchAllQuery()).
			Scroll("30m").
			Size(2000).
			Pretty(true).
			Do(ctx) // 执行请求
		aggpage := 1
		scroll_id := searchResult.ScrollId
		for {
			bulkRequest := databases.Esdb.Es.Bulk()
			for _, val := range searchResult.Hits.Hits {
				tmpMap := make(map[string]interface{})
				json.Unmarshal(*val.Source, &tmpMap)
				log.Print("host", tmpMap["host"])

				if strings.Hwangxinrefix(tmpMap["host"].(string), "https://") || strings.Hwangxinrefix(tmpMap["host"].(string), "https://") {
					link := tmpMap["host"].(string)
					newUrl, _ := url.Parse(link)
					newUrlHost := newUrl.Hostname()
					if Utils.F.IsDomain(newUrlHost) {
						tmpMap["isdomain"] = true
						tmpMap["host"] = newUrlHost
					}
				} else {
					if Utils.F.IsDomain(tmpMap["host"].(string)) {
						tmpMap["isdomain"] = true
					}
				}
				if strings.Hwangxinrefix(val.Id, "https://") {
					tmpMap["protocol"] = "https"
				} else {
					tmpMap["protocol"] = "http"
				}
				updatereq := elastic.NewBulkUpdateRequest().Index(SCAN_SUBDOMAIN_INDEX).Type(SCAN_SUBDOMAIN_TYPE).Id(val.Id).Doc(tmpMap)
				bulkRequest = bulkRequest.Add(updatereq)
			}
			_, subdomainErr := bulkRequest.Refresh("true").Do(ctx)
			if nil != subdomainErr {
				log.Print("subdomainErr", subdomainErr)
				continue
			}
			aggpage++
			scroll_id = searchResult.ScrollId
			scrollData, _ := databases.Esdb.Es.Scroll("30m").ScrollId(searchResult.ScrollId).Do(ctx)
			log.Print("subdomain clear", aggpage)
			searchResult = scrollData
			if nil == searchResult || 0 == len(searchResult.Hits.Hits) {
				log.Print("searchResult", searchResult)
				break
			}
		}
		databases.Esdb.Es.ClearScroll(scroll_id)
	}()
	Constant.Success(c, "运行成功", "")
	return
}

func (f wangansearchController) Running(c *gin.Context) {
	wanganvalue, wanganerror := databases.Redis.Get("wangansearch:page").Result()
	if "" != wanganvalue || wanganerror != redis.Nil {
		// 查询不到数据
		Constant.Error(c, Constant.wangan_SYNC_TASK_EXISTS)
		return
	}
	Constant.Success(c, "当前无导入任务", "")
	return
}

func (f wangansearchController) SingleImport(c *gin.Context) {
	var params []wangansearchImport
	if err := c.ShouldBindJSON(&params); err != nil {
		Constant.Error(c, Request.ValidatorError(err))
		return
	}
	fmt.Printf("%#v", params)
	Subdomain_resolve := make([]map[string]interface{}, 0)
	for _, value := range params {
		tmp := map[string]interface{}{}
		tmp["ip_segment"] = value.Ip
		tmp["ip"] = value.Ip
		tmp["ip_type"] = Constant.IP_TYPE_V4
		tmp["createtime"] = time.Now().Format("2006-01-02 15:04:05")
		tmp["lastupdatetime"] = time.Now().Format("2006-01-02 15:04:05")
		tmp["belong_user_id"] = f.GetUserId(c)
		tmp["subdomain_clue_id"] = 0
		start_ip, _ := exnet.IPString2Long(value.Ip)
		tmp["start_ip"] = int64(start_ip)
		tmp["end_ip"] = int64(start_ip)
		tmp["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS
		tmp["company_id"] = 0
		tmp["user_id"] = f.GetUserId(c)
		tmp["port"] = value.Port
		tmp["protocol"] = value.Protocol
		tmp["country"] = value.Country
		tmp["city"] = value.City
		tmp["title"] = value.Title
		tmp["banner"] = value.Banner
		tmp["body"] = value.Body
		tmp["server"] = value.Server
		tmp["domain"] = value.Domain
		tmp["country_name"] = value.Country_name
		tmp["region"] = value.Region
		tmp["longitude"] = value.Longitude
		tmp["latitude"] = value.Latitude
		tmp["product"] = value.Product
		tmp["product_category"] = value.Product_category
		tmp["header"] = value.Header
		tmp["fid"] = value.Fid
		tmp["structinfo"] = value.Structinfo
		tmp["icp"] = value.Icp
		tmp["os"] = value.Os
		Subdomain_resolve = append(Subdomain_resolve, tmp)
	}
	go Services.wangansearchService.Sync(Subdomain_resolve)
	Constant.Success(c, "成功", "")
	return
}

func (f wangansearchController) GatherAssetFromScan(c *gin.Context) {
	go func() {
		conf := databases.Conf
		ctx := context.Background()
		SCAN_ASSET_INDEX := conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_INDEX").String()
		SCAN_ASSET_TYPE := conf.Section("es").Key("ES_wangxin_WX_SCAN_AIMS_TYPE").String()
		user_id := f.GetUserId(c)
		boolQuery := elastic.NewBoolQuery()
		boolQuery.MustNot(elastic.NewTermQuery("company_id", 0))
		AssetSearchResult, _ := databases.Esdb.Es.Scroll().
			Index(SCAN_ASSET_INDEX).
			Type(SCAN_ASSET_TYPE).
			Query(boolQuery).
			Scroll("15m").
			Size(1000).
			Pretty(true).
			Do(ctx) // 执行请求
		scroll_id := AssetSearchResult.ScrollId
		aggpage := 1
		for {
			for _, val := range AssetSearchResult.Hits.Hits {
				tmpMap := make(map[string]interface{})
				json.Unmarshal(*val.Source, &tmpMap)
				Utils.ImportAssetPool.Submit(Models.AssetModel.GatherAsset(tmpMap["ip_segment"].(string), user_id, Utils.F.Ti(tmpMap["company_id"])))
			}
			aggpage++
			scroll_id = AssetSearchResult.ScrollId
			scrollData, _ := databases.Esdb.Es.Scroll("15m").ScrollId(AssetSearchResult.ScrollId).Do(ctx)
			log.Print("insert asset process", aggpage)
			AssetSearchResult = scrollData
			if nil == AssetSearchResult || 0 == len(AssetSearchResult.Hits.Hits) {
				break
			}
		}
		databases.Esdb.Es.ClearScroll(scroll_id)
	}()
	Constant.Success(c, "运行成功！", "")
}

func (f wangansearchController) TaggedAsset(c *gin.Context) {
	go func() {
		conf := databases.Conf
		ctx := context.Background()
		SCAN_ASSET_INDEX := conf.Section("es").Key("ES_wangxin_ASSETS_INDEX").String()
		SCAN_ASSET_TYPE := conf.Section("es").Key("ES_wangxin_ASSETS_TYPE").String()
		boolQuery := elastic.NewBoolQuery()
		AssetSearchResult, _ := databases.Esdb.Es.Scroll().
			Index(SCAN_ASSET_INDEX).
			Type(SCAN_ASSET_TYPE).
			Query(boolQuery).
			Scroll("120m").
			Size(1000).
			Pretty(true).
			Do(ctx) // 执行请求

		aggpage := 1
		scroll_id := AssetSearchResult.ScrollId
		retryNum := 0
		for {
			if 0 == retryNum {
				temArr := make([]interface{}, 0)
				for _, val := range AssetSearchResult.Hits.Hits {
					tmpMap := make(map[string]interface{})
					json.Unmarshal(*val.Source, &tmpMap)
					temArr = append(temArr, tmpMap["ip"])
				}
				ipArr, _ := Utils.F.Unique(temArr)
				for _, ipval := range ipArr.([]interface{}) {
					Utils.ImportAssetPool.Submit(Models.AssetModel.TaggedAsset(ipval.(string)))
				}
			}
			aggpage++
			scroll_id = AssetSearchResult.ScrollId
			log.Print("insert asset process start", aggpage)
			scrollData, scroll_err := databases.Esdb.Es.Scroll("120m").ScrollId(scroll_id).Do(ctx)
			log.Print("insert asset process end", aggpage)
			if nil != scroll_err && retryNum < 4 {
				retryNum++
				log.Print("insert asset process", scroll_err)
				continue
			}
			retryNum = 0
			AssetSearchResult = scrollData
			if nil == AssetSearchResult || 0 == len(AssetSearchResult.Hits.Hits) {
				log.Print("insert asset process break", AssetSearchResult)
				break
			}
		}
		databases.Esdb.Es.ClearScroll(scroll_id)
	}()
	Constant.Success(c, "运行成功！", "")
}
