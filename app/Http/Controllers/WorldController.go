package Controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"wangxin2.0/app/Constant"
)

type WorldController struct {
}

func (w WorldController) Show(c *gin.Context) {
	jsonFile, err := os.Open(Constant.GoAbsoulePath + "/datav/world.geo.json")
	if err != nil {
		Constant.Error(c, Constant.READ_JSON_ERROR)
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		Constant.Error(c, Constant.READ_JSON_ERROR)
		return
	}
	var worldJson map[string]*json.RawMessage
	unmarErr := json.Unmarshal(jsonData, &worldJson)
	if unmarErr != nil {
		return
	}
	Constant.Success(c, "success", worldJson)
}
