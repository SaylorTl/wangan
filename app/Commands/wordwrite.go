package Commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
	"wangxin2.0/app/Constant"
	"wangxin2.0/app/Models"
)

func init() {
	rootCmd.AddCommand(worldwriteCmd)
}

var worldwriteCmd = &cobra.Command{
	Use:   "worldwrite",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, err := os.Open(Constant.GoAbsoulePath + "/datav/world.geo.json")
		if err != nil {
			fmt.Println(Constant.READ_JSON_ERROR)
			return
		}
		jsonData, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(Constant.READ_JSON_ERROR)
			return
		}
		var worldJson map[string][]map[string]map[string]interface{}
		json.Unmarshal(jsonData, &worldJson)
		//if unmarErr != nil {
		//	Constant.Error(c, unmarErr.Error())
		//	return
		//}
		insertData := &Models.Wx_geoatlas{
			Id:        0,
			Parent_id: 0,
			Lft:       1,
			Rgt:       7203,
			Adcode:    400000,
			Name:      "shijie",
			Level:     1,
			Initials:  "0",
			Pinyin:    "0",
			CreatedAt: Models.LocalTime(time.Time{}.Local()),
			UpdatedAt: Models.LocalTime(time.Time{}.Local()),
		}
		shijie_id, err := Models.GeoatlaModel.Addgeotlas(insertData)
		Lft := 7203
		Rft := 7215
		if err != nil {
			fmt.Println(Constant.ADD_FAIL)
			return
		}
		for _, value := range worldJson["features"] {
			adcode := int(value["properties"]["adcode"].(float64))
			name := value["properties"]["name"].(string)
			if "China" == name || "Taiwan" == name {
				continue
			}
			insertData := &Models.Wx_geoatlas{Id: 0,
				Parent_id: shijie_id,
				Lft:       Lft,
				Rgt:       Rft,
				Adcode:    adcode,
				Name:      name,
				Level:     1,
				Initials:  "0",
				Pinyin:    "0",
				CreatedAt: Models.LocalTime(time.Time{}.Local()),
				UpdatedAt: Models.LocalTime(time.Time{}.Local()),
			}
			res_id, _ := Models.GeoatlaModel.Addgeotlas(insertData)
			Lft = Rft + 1
			Rft = Lft + 1
			fmt.Print(res_id)
		}
		updateRes, err := Models.GeoatlaModel.Updatageotlas(shijie_id, Models.Wx_geoatlas{Rgt: Rft})
		if updateRes == 0 && err != nil {
			fmt.Println(Constant.UPDATE_FAIL)
			return
		}
		updateChinaRes, err := Models.GeoatlaModel.Updatageotlas(1, Models.Wx_geoatlas{Parent_id: shijie_id})
		if updateChinaRes == 0 && err != nil {
			fmt.Println(Constant.UPDATE_FAIL)
			return
		}
		fmt.Println(Constant.ADD_SUCCESS)
	},
}
