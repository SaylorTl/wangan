package Services

import (
	"strconv"
	"wangxin2.0/app/Models"
)

var RecommGroupService recommgroupservice

type recommgroupservice struct {
}

func (a recommgroupservice) UpdateClueCount(group_id int) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		queryParams := make(map[string]interface{})
		queryParams["recom_group_id"] = group_id
		count := Models.ClueModel.Count(queryParams)
		updateParams := make(map[string]interface{})
		stringCount := strconv.FormatInt(count, 10)
		intCount, _ := strconv.Atoi(stringCount)
		updateParams["id"] = group_id
		updateParams["clue_count"] = intCount
		Models.RecomGroupModel.UpdateOrCreate(updateParams)
	}()
}

func (a recommgroupservice) UpdateAllClueCount() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		insertData := make(map[string]interface{})
		AllData := Models.RecomGroupModel.GetAllGroup(insertData)
		for _, val := range AllData {
			queryParams := make(map[string]interface{})
			queryParams["recom_group_id"] = val["id"]
			count := Models.ClueModel.Count(queryParams)
			updateParams := make(map[string]interface{})
			stringCount := strconv.FormatInt(count, 10)
			intCount, _ := strconv.Atoi(stringCount)
			updateParams["id"] = val["id"]
			updateParams["clue_count"] = intCount
			Models.RecomGroupModel.UpdateOrCreate(updateParams)
		}
	}()
}
