package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"wangxin2.0/routes"
)

func TestFind(t *testing.T) {
	globalrouter := routes.InitRouter()
	var w *httptest.ResponseRecorder
	w = Get("/api/v2/find", globalrouter)
	assert.Equal(t, 400, w.Code)
	//result := w.Body
	body, _ := ioutil.ReadAll(w.Body)

	m := make(map[string]interface{})

	err := json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		return
	}
	assert.Equal(t, float64(0), m["code"])

}
