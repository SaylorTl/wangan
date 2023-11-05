package Models

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/thinkeridea/go-extend/exnet"
	"log"
	"strconv"
	"strings"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Utils"
	"wangxin2.0/databases"
)

var AssetModel *assetmodel

type assetmodel struct {
	BaseModel
}

type Assets struct {
	Keyword        string    `json:"keyword" from:"keyword" to:"keyword" `
	Is_ipv6        int       `json:"is_ipv6" from:"is_ipv6" to:"is_ipv6"`
	Rule_tags      int       `json:"rule_tags" from:"rule_tags" to:"rule_tags"`
	Createtime     LocalTime `json:"createtime" from:"createtime" to:"createtime"`
	Lastupdatetime LocalTime `json:"lastupdatetime" from:"lastupdatetime" to:"lastupdatetime"`
}

func (a assetmodel) AssetExtractClue(queryParams map[string]interface{}) []string {
	boolQuery := elastic.NewBoolQuery()
	shouldQuery := elastic.NewBoolQuery()
	if keyword, ok := queryParams["keyword"]; ok {
		keywordShould := shouldQuery.Should(elastic.NewWildcardQuery("ip.ip_raw", "*"+keyword.(string)+"*"),
			elastic.NewWildcardQuery("host", "*"+keyword.(string)+"*"),
			elastic.NewWildcardQuery("industry", "*"+keyword.(string)+"*"),
			elastic.NewWildcardQuery("ports", "*"+keyword.(string)+"*"),
			elastic.NewWildcardQuery("protocols", "*"+keyword.(string)+"*"),
			elastic.NewWildcardQuery("rule_tags", "*"+keyword.(string)+"*"))
		boolQuery.Filter(keywordShould)
	}
	boolQuery.Filter(elastic.NewTermQuery("is_ipv6", false))
	if has_log4j, ok := queryParams["has_log4j"]; ok {
		boolQuery.Filter(elastic.NewTermQuery("has_log4j", has_log4j))
	}
	if created_at_range, ok := queryParams["created_at_range"]; ok {
		start, _ := time.Parse("2006-01-02 15:04:05", created_at_range.([]string)[0])
		boolQuery.Filter(elastic.NewRangeQuery("createtime").Gte(start.Format("2006-01-02") + " 00:00:00"))
		end, _ := time.Parse("2006-01-02 15:04:05", created_at_range.([]string)[1])
		boolQuery.Filter(elastic.NewRangeQuery("createtime").Lte(end.Format("2006-01-02") + " 23:59:59"))
	}
	if update_at_range, ok := queryParams["update_at_range"]; ok {
		update_start, _ := time.Parse("2006-01-02 15:04:05", update_at_range.([]string)[0])
		boolQuery.Filter(elastic.NewRangeQuery("lastupdatetime").Gte(update_start.Format("2006-01-02") + " 00:00:00"))
		update_end, _ := time.Parse("2006-01-02 15:04:05", update_at_range.([]string)[1])
		boolQuery.Filter(elastic.NewRangeQuery("lastupdatetime").Lte(update_end.Format("2006-01-02") + " 23:59:59"))
	}

	conf := databases.Conf
	ctx := context.Background()
	searchResult, err := databases.Esdb.Es.Scroll().
		Index(conf.Section("es").Key("ES_wangxin_ASSETS_INDEX").String()). // 设置索引名
		Type(conf.Section("es").Key("ES_wangxin_ASSETS_TYPE").String()).
		Query(boolQuery). // 设置查询条件
		Scroll("2m").
		Size(10).     // 设置分页参数 - 每页大小
		Pretty(true). // 查询结果返回可读性较好的JSON格式
		Do(ctx)
	returnData := []string{}
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	if 0 == searchResult.Hits.TotalHits {
		return returnData
	}
	for {
		for _, val := range searchResult.Hits.Hits {
			temData := make(map[string]interface{})
			err := json.Unmarshal(*val.Source, &temData)
			if err != nil {
				panic(Constant.JSON_UNENCODE_FAIL)
			}
			returnData = append(returnData, temData["ip"].(string))
		}
		scrollData, _ := databases.Esdb.Es.Scroll("2m").ScrollId(searchResult.ScrollId).Do(ctx)
		searchResult = scrollData
		if nil == searchResult || 0 == len(searchResult.Hits.Hits) {
			break
		}
	}
	return returnData
}

func (a assetmodel) searchByAfter(searchAfter string, query_ip string) *elastic.SearchResult {
	ctx := context.Background()
	conf := databases.Conf
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(elastic.NewTermQuery("ip", query_ip))
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_SERVICE_INDEX").String()
	SCAN_SUBDOMAIN_INDEX := conf.Section("es").Key("ES_wangxin_SUBDOMAIN_INDEX").String()
	searchService := databases.Esdb.Es.Search(SCAN_INDEX, SCAN_SUBDOMAIN_INDEX).
		Query(boolQuery).
		// 每页的个数
		Size(2000).
		Sort("_id", false).
		Pretty(true)
	// 为了兼容第一次请求，判断searchAfter是否为空
	if searchAfter != "" {
		searchService.SearchAfter(searchAfter)
	}
	// 执行查询
	res, err := searchService.Do(ctx)
	if err != nil {
		log.Print(Constant.JSON_UNENCODE_FAIL)
	}
	// 本例根据Sn字段排序，所以返回最后一条数据的Sn
	return res
}

func (a assetmodel) TaggedAsset(queryip string) assetimportFunc {
	return func() {
		searchResult := a.searchByAfter("", queryip)
		log.Print("searchResult", searchResult)
		if nil == searchResult {
			return
		}
		company_tags := make([]string, 0)
		rule_infos := make([]map[string]interface{}, 0)
		cat_tags := make([]string, 0)
		rule_tags := make([]string, 0)
		port_strings := make([]int, 0)
		value := make(map[string]interface{})
		if nil != searchResult && 0 != searchResult.Hits.TotalHits {
			i := 0
			for {
				for _, hits := range searchResult.Hits.Hits {
					i++
					tmpdata := make(map[string]interface{})
					err := json.Unmarshal(*hits.Source, &tmpdata)
					if err != nil {
						log.Print(Constant.JSON_UNENCODE_FAIL)
						return
					}
					log.Print("product", tmpdata["product"])
					if _, ok := tmpdata["product"]; !ok {
						continue
					}
					//如果是subdomain单独判断是否是https或http
					for _, product := range tmpdata["product"].([]interface{}) {
						if "" == product {
							continue
						}
						rule_item := make(map[string]interface{})
						rule_val := RulesModel.FindOne(product.(string))
						rule_item["rule_id"] = rule_val.Id
						rule_item["second_cat_tag"] = rule_val.WxSecondCategoryRules.WxSecondCategories.Title
						rule_item["first_cat_tag"] = ""
						if "" != rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry {
							first_cat_id, _ := strconv.Atoi(rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry)
							WxFirstCategories := SecondCategoriesModel.Find(first_cat_id)
							rule_item["first_cat_tag"] = WxFirstCategories.Title
						}
						rule_item["level_code"] = rule_val.Level_code
						rule_item["company"] = rule_val.Company
						rule_item["soft_hard_code"] = rule_val.Soft_hard_code
						rule_item["title"] = rule_val.Product
						rule_item["ports"] = append(port_strings, Utils.F.Ti(tmpdata["port"]))
						cat_tags = append(cat_tags, rule_item["second_cat_tag"].(string))
						rule_tags = append(rule_tags, rule_item["title"].(string))
						company_tags = append(company_tags, rule_item["company"].(string))
						rule_infos = append(rule_infos, rule_item)
					}
				}
				if len(searchResult.Hits.Hits) < 2000 {
					break
				}
				SeqNo := searchResult.Hits.Hits[len(searchResult.Hits.Hits)-1].Sort[0]
				tempSearchResult := a.searchByAfter(SeqNo.(string), queryip)
				if nil == tempSearchResult || 0 == len(tempSearchResult.Hits.Hits) {
					break
				}
				log.Print("tagged query_ip 翻页", queryip)
				searchResult = tempSearchResult
			}
			value["company_tags"], _ = Utils.F.Unique(company_tags)
			value["rule_tags"], _ = Utils.F.Unique(rule_tags)
			value["cat_tags"], _ = Utils.F.Unique(cat_tags)
			value["rule_infos"], _ = Utils.F.UniqueMap(rule_infos, "rule_id")

			//计算软硬件数量
			first_tag_num := make([]map[string]interface{}, 0)
			second_tag_num := make([]map[string]interface{}, 0)
			firsr_cat_arr := make(map[string][]interface{})
			second_cat_arr := make(map[string][]interface{})
			for _, rule_kk := range rule_infos {
				if 0 != rule_kk["soft_hard_code"] {
					first_tag_name := rule_kk["first_cat_tag"].(string)
					firsr_cat_arr[first_tag_name] = append(firsr_cat_arr[first_tag_name], rule_kk)
					second_tag_name := rule_kk["second_cat_tag"].(string)
					second_cat_arr[second_tag_name] = append(second_cat_arr[second_tag_name], rule_kk)
				}
			}
			value["software_num"] = 0
			for first_kk, first_rule := range firsr_cat_arr {
				first_rule_map := make(map[string]interface{})
				first_rule_map["first_cat_tag"] = first_kk
				first_rule_map["num"] = len(first_rule)
				if first_rule[0].(map[string]interface{})["soft_hard_code"] == 2 {
					first_rule_map["software_num"] = first_rule_map["num"]
					value["software_num"] = value["software_num"].(int) + first_rule_map["num"].(int)
					first_rule_map["hardware_num"] = 0
				} else {
					first_rule_map["hardware_num"] = first_rule_map["num"]
					first_rule_map["software_num"] = 0
				}
				first_tag_num = append(first_tag_num, first_rule_map)
			}
			value["first_tag_num"] = first_tag_num
			for second_kk, second_rule := range second_cat_arr {
				second_rule_map := make(map[string]interface{})
				second_rule_map["second_cat_tag"] = second_kk
				second_rule_map["num"] = len(second_rule)
				if second_rule[0].(map[string]interface{})["soft_hard_code"] == 2 {
					second_rule_map["software_num"] = second_rule_map["num"]
					second_rule_map["hardware_num"] = 0
				} else {
					second_rule_map["hardware_num"] = second_rule_map["num"]
					second_rule_map["software_num"] = 0
				}
				second_tag_num = append(second_tag_num, second_rule_map)
			}

			value["second_tag_num"] = second_tag_num
			log.Print("tagged query_ip 入库  ", queryip, " ", value)
			if 0 == len(rule_tags) {
				return
			}
			_id := queryip
			log.Print("tagged query_ip 入库", _id, " ", time.Now())
			a.RetryInsertAssetFromScan(_id, value, 0)
		}
		return
	}
}

func (a assetmodel) GatherAsset(queryip string, user_id int, company_id int) assetimportFunc {
	return func() {
		searchResult := a.searchByAfter("", queryip)
		if nil == searchResult {
			return
		}
		if nil != searchResult && 0 != searchResult.Hits.TotalHits {
			_id := ""
			value := make(map[string]interface{})
			port_lists := make([]interface{}, 0)
			hosts := make([]interface{}, 0)
			ports := make([]interface{}, 0)
			protocols := make([]interface{}, 0)

			company_tags := make([]string, 0)
			rule_infos := make([]map[string]interface{}, 0)
			cat_tags := make([]string, 0)
			rule_tags := make([]string, 0)
			port_strings := make([]int, 0)

			i := 0
			for {
				for _, hits := range searchResult.Hits.Hits {
					i++
					tmpdata := make(map[string]interface{})
					err := json.Unmarshal(*hits.Source, &tmpdata)
					if err != nil {
						log.Print(Constant.JSON_UNENCODE_FAIL)
						return
					}
					//如果是subdomain单独判断是否是https或http
					if "wanganee_subdomain" == hits.Index {
						if strings.Hwangxinrefix(hits.Id, "https://") {
							tmpdata["protocol"] = "https"
						} else {
							tmpdata["protocol"] = "http"
						}
					}
					value["createtime"] = time.Now().Format("2006-01-02 15:04:05")
					value["lastupdatetime"] = tmpdata["lastupdatetime"].(string)
					value["belong_user_id"] = user_id
					value["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS
					value["company_id"] = company_id
					value["is_ipv6"] = false
					value["ip"] = tmpdata["ip"].(string)
					if _, ok := tmpdata["geoip"]; ok {
						value["geoip"] = tmpdata["geoip"]
					}

					UserQueryData := make(map[string]interface{})
					UserQueryData["id"] = company_id
					company := CompaniesModel.Find(UserQueryData)
					value["geoatla_lft_id"] = company.Geoatla_lft_id

					_id = tmpdata["ip"].(string)
					if _, ok := tmpdata["host"]; ok && Utils.F.IsDomain(tmpdata["host"].(string)) {
						if "" != tmpdata["host"].(string) && !Utils.F.InArray(tmpdata["host"].(string), hosts) {
							hosts = append(hosts, tmpdata["host"].(string))
						}
					}
					value["hosts"] = hosts
					if _, ok := tmpdata["protocol"]; ok {
						tmpdata["port"] = Utils.F.Ts(tmpdata["port"])
						if "" != tmpdata["port"].(string) && !Utils.F.InArray(tmpdata["port"].(string), ports) {
							ports = append(ports, tmpdata["port"].(string))
						}
					}
					value["ports"] = ports
					if _, ok := tmpdata["protocol"]; ok {
						if "" != tmpdata["protocol"].(string) && !Utils.F.InArray(tmpdata["protocol"].(string), protocols) {
							protocols = append(protocols, tmpdata["protocol"].(string))
						}
					}

					value["protocols"] = protocols
					//port_list聚合
					if _, ok := tmpdata["protocol"]; ok {
						port_list := make(map[string]interface{})
						port_list["protocol"] = tmpdata["protocol"].(string)
						port_list["port"] = tmpdata["port"].(string)
						port_list["banner"] = ""
						if _, bok := tmpdata["banner"]; bok {
							port_list["banner"] = tmpdata["banner"].(string)
						}
						is_exist_port := false
						for kk, tmpportlist := range port_lists {
							if tmpportlist.(map[string]interface{})["port"].(string) == tmpdata["port"].(string) {
								is_exist_port = true
								if _, bok := tmpdata["banner"]; bok && "" == tmpportlist.(map[string]interface{})["banner"].(string) {
									port_lists[kk].(map[string]interface{})["banner"] = tmpdata["banner"].(string)
									port_lists[kk].(map[string]interface{})["protocol"] = tmpdata["protocol"].(string)
								}
							}
						}
						if false == is_exist_port {
							port_lists = append(port_lists, port_list)
						}
					}
					//如果是subdomain单独判断是否是https或http
					for _, product := range tmpdata["product"].([]interface{}) {
						if "" == product {
							continue
						}
						rule_item := make(map[string]interface{})
						rule_val := RulesModel.FindOne(product.(string))
						rule_item["rule_id"] = rule_val.Id
						rule_item["second_cat_tag"] = rule_val.WxSecondCategoryRules.WxSecondCategories.Title
						rule_item["first_cat_tag"] = ""
						if "" != rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry {
							first_cat_id, _ := strconv.Atoi(rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry)
							WxFirstCategories := SecondCategoriesModel.Find(first_cat_id)
							rule_item["first_cat_tag"] = WxFirstCategories.Title
						}
						rule_item["level_code"] = rule_val.Level_code
						rule_item["company"] = rule_val.Company
						rule_item["soft_hard_code"] = rule_val.Soft_hard_code
						rule_item["title"] = rule_val.Product
						rule_item["ports"] = append(port_strings, Utils.F.Ti(tmpdata["port"]))
						cat_tags = append(cat_tags, rule_item["second_cat_tag"].(string))
						rule_tags = append(rule_tags, rule_item["title"].(string))
						company_tags = append(company_tags, rule_item["company"].(string))
						rule_infos = append(rule_infos, rule_item)
					}
					value["port_list"] = port_lists
				}
				if len(searchResult.Hits.Hits) < 2000 {
					break
				}
				SeqNo := searchResult.Hits.Hits[len(searchResult.Hits.Hits)-1].Sort[0]
				tempSearchResult := a.searchByAfter(SeqNo.(string), queryip)
				if nil == tempSearchResult || 0 == len(tempSearchResult.Hits.Hits) {
					break
				}
				log.Print("query_ip 翻页", queryip)
				searchResult = tempSearchResult
			}

			value["company_tags"], _ = Utils.F.Unique(company_tags)
			value["rule_tags"], _ = Utils.F.Unique(rule_tags)
			value["cat_tags"], _ = Utils.F.Unique(cat_tags)
			value["rule_infos"], _ = Utils.F.UniqueMap(rule_infos, "rule_id")

			//计算软硬件数量
			first_tag_num := make([]map[string]interface{}, 0)
			second_tag_num := make([]map[string]interface{}, 0)
			firsr_cat_arr := make(map[string][]interface{})
			second_cat_arr := make(map[string][]interface{})
			for _, rule_kk := range rule_infos {
				first_tag_name := rule_kk["first_cat_tag"].(string)
				firsr_cat_arr[first_tag_name] = append(firsr_cat_arr[first_tag_name], rule_kk)
				second_tag_name := rule_kk["second_cat_tag"].(string)
				second_cat_arr[second_tag_name] = append(second_cat_arr[second_tag_name], rule_kk)
			}
			value["software_num"] = 0
			for first_kk, first_rule := range firsr_cat_arr {
				first_rule_map := make(map[string]interface{})
				first_rule_map["first_cat_tag"] = first_kk
				first_rule_map["num"] = len(first_rule)
				if first_rule[0].(map[string]interface{})["soft_hard_code"] == 2 {
					first_rule_map["software_num"] = first_rule_map["num"]
					value["software_num"] = value["software_num"].(int) + first_rule_map["num"].(int)
					first_rule_map["hardware_num"] = 0
				} else {
					first_rule_map["hardware_num"] = first_rule_map["num"]
					first_rule_map["software_num"] = 0
				}
				first_tag_num = append(first_tag_num, first_rule_map)
			}
			value["first_tag_num"] = first_tag_num
			for second_kk, second_rule := range second_cat_arr {
				second_rule_map := make(map[string]interface{})
				second_rule_map["second_cat_tag"] = second_kk
				second_rule_map["num"] = len(second_rule)
				if second_rule[0].(map[string]interface{})["soft_hard_code"] == 2 {
					second_rule_map["software_num"] = second_rule_map["num"]
					second_rule_map["hardware_num"] = 0
				} else {
					second_rule_map["hardware_num"] = second_rule_map["num"]
					second_rule_map["software_num"] = 0
				}
				second_tag_num = append(second_tag_num, second_rule_map)
			}

			value["second_tag_num"] = second_tag_num
			log.Print("query_ip 入库", queryip, " ", time.Now())
			a.RetryInsertAssetFromScan(_id, value, 0)
		}
		return
	}
}

type assetimportFunc func()

func (a assetmodel) InsertAsset(queryip string, user_id int, company_id int) assetimportFunc {
	return func() {
		AssetboolQuery := elastic.NewBoolQuery()
		AssetboolQuery.Filter(elastic.NewTermQuery("ip_segment", queryip))
		searchResult := a.searchByAfter("", queryip)
		if nil == searchResult {
			return
		}
		log.Print("query_ip ", queryip, " ", time.Now())
		if nil != searchResult && 0 != int(searchResult.Hits.TotalHits) {
			_id := ""
			value := make(map[string]interface{})
			port_lists := make([]interface{}, 0)
			hosts := make([]interface{}, 0)
			ports := make([]interface{}, 0)
			protocols := make([]interface{}, 0)

			company_tags := make([]string, 0)
			rule_infos := make([]map[string]interface{}, 0)
			cat_tags := make([]string, 0)
			rule_tags := make([]string, 0)
			port_strings := make(map[string][]string)
			products_string := make([]string, 0)
			i := 0
			for {
				for _, hits := range searchResult.Hits.Hits {
					i++
					tmpdata := make(map[string]interface{})
					err := json.Unmarshal(*hits.Source, &tmpdata)
					if err != nil {
						log.Print(Constant.JSON_UNENCODE_FAIL)
						return
					}
					value["createtime"] = time.Now().Format("2006-01-02 15:04:05")
					value["lastupdatetime"] = tmpdata["lastupdatetime"].(string)
					value["belong_user_id"] = user_id
					value["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS
					value["company_id"] = company_id
					value["is_ipv6"] = false
					value["ip"] = tmpdata["ip"].(string)
					if _, ok := tmpdata["geoip"]; ok {
						value["geoip"] = tmpdata["geoip"]
					}
					UserQueryData := make(map[string]interface{})
					UserQueryData["id"] = user_id
					userData, _ := UserModel.FindUser(UserQueryData)
					geoatla := GeoatlaModel.Find(map[string]interface{}{"adcode": userData["adcode"]})
					value["geoatla_lft_id"] = geoatla.Lft

					_id = tmpdata["ip"].(string)
					if _, ok := tmpdata["host"]; ok && Utils.F.IsDomain(tmpdata["host"].(string)) {
						if "" != tmpdata["host"].(string) && !Utils.F.InArray(tmpdata["host"].(string), hosts) {
							hosts = append(hosts, tmpdata["host"].(string))
						}
					}
					value["hosts"] = hosts
					if _, ok := tmpdata["protocol"]; ok {
						tmpdata["port"] = Utils.F.Ts(tmpdata["port"])
						if "" != tmpdata["port"].(string) && !Utils.F.InArray(tmpdata["port"].(string), ports) {
							ports = append(ports, tmpdata["port"].(string))
						}
					}
					value["ports"] = ports
					if _, ok := tmpdata["protocol"]; ok {
						if "" != tmpdata["protocol"].(string) && !Utils.F.InArray(tmpdata["protocol"].(string), protocols) {
							protocols = append(protocols, tmpdata["protocol"].(string))
						}
					}

					value["protocols"] = protocols
					//port_list聚合
					if _, ok := tmpdata["protocol"]; ok {
						for kk, tmpportlist := range port_lists {
							if tmpportlist.(map[string]interface{})["port"].(string) == tmpdata["port"].(string) {
								port_lists = append(port_lists[:kk], port_lists[kk+1:]...)
							}
						}
						port_list := make(map[string]interface{})
						port_list["protocol"] = tmpdata["protocol"].(string)
						port_list["port"] = tmpdata["port"].(string)
						port_list["banner"] = ""
						if _, bok := tmpdata["banner"]; bok {
							port_list["banner"] = tmpdata["banner"].(string)
						}
						port_lists = append(port_lists, port_list)
					}
					value["port_list"] = port_lists
					log.Print("product", tmpdata["product"])
					if _, ok := tmpdata["product"]; !ok {
						continue
					}
					//如果是subdomain单独判断是否是https或http
					for _, product := range tmpdata["product"].([]interface{}) {
						if "" == product {
							continue
						}
						//收集组件
						products_string = append(products_string, product.(string))
						tmp_products_string, _ := Utils.F.Unique(products_string)
						products_string = tmp_products_string.([]string)

						//收集组件端口
						port_strings[product.(string)] = append(port_strings[product.(string)], tmpdata["port"].(string))
						tmp_port_strings, _ := Utils.F.Unique(port_strings[product.(string)])
						port_strings[product.(string)] = tmp_port_strings.([]string)
					}

				}
				if len(searchResult.Hits.Hits) < 2000 {
					break
				}
				SeqNo := searchResult.Hits.Hits[len(searchResult.Hits.Hits)-1].Sort[0]
				tempSearchResult := a.searchByAfter(SeqNo.(string), queryip)
				if nil == tempSearchResult || 0 == len(tempSearchResult.Hits.Hits) {
					break
				}
				log.Print("query_ip 翻页", queryip)
				searchResult = tempSearchResult
			}
			for _, product := range products_string {
				rule_item := make(map[string]interface{})
				rule_val := RulesModel.FindOne(product)
				if "" == rule_val.Product {
					continue
				}
				rule_item["rule_id"] = rule_val.Id
				rule_item["second_cat_tag"] = rule_val.WxSecondCategoryRules.WxSecondCategories.Title
				rule_item["first_cat_tag"] = ""
				if "" != rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry {
					first_cat_id, _ := strconv.Atoi(rule_val.WxSecondCategoryRules.WxSecondCategories.Ancestry)
					WxFirstCategories := SecondCategoriesModel.Find(first_cat_id)
					rule_item["first_cat_tag"] = WxFirstCategories.Title
				}
				rule_item["level_code"] = rule_val.Level_code
				rule_item["company"] = rule_val.Company
				rule_item["soft_hard_code"] = rule_val.Soft_hard_code
				rule_item["title"] = rule_val.Product
				rule_item["ports"] = port_strings[product]
				cat_tags = append(cat_tags, rule_item["second_cat_tag"].(string))
				rule_tags = append(rule_tags, rule_item["title"].(string))
				company_tags = append(company_tags, rule_item["company"].(string))
				rule_infos = append(rule_infos, rule_item)
			}

			value["company_tags"], _ = Utils.F.Unique(company_tags)
			value["rule_tags"], _ = Utils.F.Unique(rule_tags)
			value["cat_tags"], _ = Utils.F.Unique(cat_tags)
			value["rule_infos"], _ = Utils.F.UniqueMap(rule_infos, "rule_id")

			//计算软硬件数量
			first_tag_num := make([]map[string]interface{}, 0)
			second_tag_num := make([]map[string]interface{}, 0)
			firsr_cat_arr := make(map[string][]interface{})
			second_cat_arr := make(map[string][]interface{})
			for _, rule_kk := range rule_infos {
				first_tag_name := rule_kk["first_cat_tag"].(string)
				firsr_cat_arr[first_tag_name] = append(firsr_cat_arr[first_tag_name], rule_kk)
				second_tag_name := rule_kk["second_cat_tag"].(string)
				second_cat_arr[second_tag_name] = append(second_cat_arr[second_tag_name], rule_kk)
			}
			value["software_num"] = 0
			for first_kk, first_rule := range firsr_cat_arr {
				first_rule_map := make(map[string]interface{})
				first_rule_map["first_cat_tag"] = first_kk
				first_rule_map["num"] = len(first_rule)
				first_rule_map["hardware_num"] = 0
				first_rule_map["software_num"] = 0
				if first_rule[0].(map[string]interface{})["soft_hard_code"] == 2 {
					first_rule_map["software_num"] = first_rule_map["software_num"].(int) + 1
					value["software_num"] = value["software_num"].(int) + first_rule_map["num"].(int)
				}
				if first_rule[0].(map[string]interface{})["soft_hard_code"] == 1 {
					first_rule_map["hardware_num"] = first_rule_map["hardware_num"].(int) + 1
				}
				first_tag_num = append(first_tag_num, first_rule_map)
			}
			value["first_tag_num"] = first_tag_num
			for second_kk, second_rule := range second_cat_arr {
				second_rule_map := make(map[string]interface{})
				second_rule_map["second_cat_tag"] = second_kk
				second_rule_map["num"] = len(second_rule)
				second_rule_map["hardware_num"] = 0
				second_rule_map["software_num"] = 0
				if second_rule[0].(map[string]interface{})["soft_hard_code"] == 2 {
					second_rule_map["software_num"] = second_rule_map["software_num"].(int) + 1
				}
				if second_rule[0].(map[string]interface{})["soft_hard_code"] == 1 {
					second_rule_map["hardware_num"] = second_rule_map["hardware_num"].(int) + 1
				}
				second_tag_num = append(second_tag_num, second_rule_map)
			}
			value["second_tag_num"] = second_tag_num
			log.Print("query_ip 入库", queryip, " ", time.Now())
			a.RetryInsertScanAims(_id, value, 0)
		}
		return
	}
}

func (a assetmodel) RetryInsertScanAims(_id string, value map[string]interface{}, retrynum int) bool {
	if retrynum >= 10 {
		return false
	}
	conf := databases.Conf
	updateParams := Utils.F.C(value)
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(elastic.NewTermQuery("_id", _id))
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_ASSETS_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_ASSETS_TYPE").String()
	updateParams["port_size"] = len(updateParams["ports"].([]interface{}))
	updateParams["task_ids"] = []string{"0"}
	scerr := databases.Esdb.Upsert(context.Background(), SCAN_INDEX, SCAN_TYPE, _id, updateParams)
	if scerr != nil {
		time.Sleep(5 * time.Second)
		retrynum++
		log.Print("Upsert_asset_retry"+strconv.Itoa(retrynum)+":", scerr)
		a.RetryInsertScanAims(_id, value, retrynum)
	}
	start_ip, _ := exnet.IPString2Long(value["ip"].(string))
	value["ip_segment"] = value["ip"].(string)
	value["start_ip"] = int64(start_ip)
	value["end_ip"] = int64(start_ip)

	value["ip_type"] = Constant.IP_TYPE_V4
	value["belong_user_id"] = updateParams["belong_user_id"]
	value["ip_source"] = Constant.SOURCE_wanganSYNC_ASSETS

	ScanAimsModel.SyncScanAims([]map[string]interface{}{Utils.F.C(value)})
	for _, domain := range value["hosts"].([]interface{}) {
		domainParams := Utils.F.C(value)
		domainParams["domain"] = domain
		ScanSubdomainsModel.SyncScanSubdomains([]map[string]interface{}{domainParams})
	}
	return true
}

func (a assetmodel) RetryInsertAssetFromScan(_id string, value map[string]interface{}, retrynum int) bool {
	if retrynum >= 10 {
		return false
	}
	conf := databases.Conf
	updateParams := Utils.F.C(value)
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(elastic.NewTermQuery("_id", _id))
	SCAN_INDEX := conf.Section("es").Key("ES_wangxin_ASSETS_INDEX").String()
	SCAN_TYPE := conf.Section("es").Key("ES_wangxin_ASSETS_TYPE").String()
	if _, ok := updateParams["ports"]; ok {
		updateParams["port_size"] = len(updateParams["ports"].([]interface{}))
	}
	updateParams["task_ids"] = []string{"0"}
	scerr := databases.Esdb.Upsert(context.Background(), SCAN_INDEX, SCAN_TYPE, _id, updateParams)
	if scerr != nil {
		time.Sleep(5 * time.Second)
		retrynum++
		log.Print("Upsert_asset_retry"+strconv.Itoa(retrynum)+":", scerr)
		a.RetryInsertScanAims(_id, value, retrynum)
	}
	return true
}
