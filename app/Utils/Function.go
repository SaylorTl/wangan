package Utils

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
	"wangxin2.0/app/Constant"
	"wangxin2.0/databases"
)

var F *FunctionClass

type FunctionClass struct {
}

func InitFunction() {
	func() {
		F = &FunctionClass{}
	}()
}

type sliceError struct {
	msg string
}

func (e *sliceError) Error() string {
	return e.msg
}

// 截取字符串
func (f FunctionClass) Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// 获取上级目录
func (f FunctionClass) GetParentDirectory(dirctory string) string {
	return f.Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

// 是否为域名
func (f FunctionClass) IsDomain(domain string) bool {
	var match bool
	//开头并且域名中间有/的情况
	IsLine := "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}(/)"
	//开头并且域名中间没有/的情况
	NotLine := "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}"
	match, _ = regexp.MatchString(IsLine, domain)
	if !match {
		match, _ = regexp.MatchString(NotLine, domain)
	}
	return match
}

func (f FunctionClass) HttpBuildQuery(params map[string]interface{}, parentKey string) (param_str string) {
	params_arr := make([]string, 0)
	for k, v := range params {
		if vals, ok := v.(map[string]interface{}); ok {
			if parentKey != "" {
				k = fmt.Sprintf("%s[%s]", parentKey, k)
			}
			params_arr = append(params_arr, f.HttpBuildQuery(vals, k))
		} else {
			if parentKey != "" {
				params_arr = append(params_arr, fmt.Sprintf("%s[%s]=%s", parentKey, k, gconv.String(v)))
			} else {
				params_arr = append(params_arr, fmt.Sprintf("%s=%s", k, gconv.String(v)))
			}
		}
	}
	param_str = strings.Join(params_arr, "&")
	return param_str
}

// 字符串反转函数
func (f FunctionClass) Reverse(str string) string {
	rs := []rune(str)
	len := len(rs)
	var tt []rune

	tt = make([]rune, 0)
	for i := 0; i < len; i++ {
		tt = append(tt, rs[len-i-1])
	}
	return string(tt[0:])
}

// 补白函数
func (f FunctionClass) StrPad(input string, padLength int, padString string, padType string) string {
	output := ""
	inputLen := len(input)
	if inputLen >= padLength {
		return input
	}
	padStringLen := len(padString)
	needFillLen := padLength - inputLen
	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}
	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
}

// 提取结构体中某列
func (f FunctionClass) StructColumn(structSlice []interface{}, key string) []interface{} {
	rt := reflect.TypeOf(structSlice)
	rv := reflect.ValueOf(structSlice)
	if rt.Kind() == reflect.Slice { //切片类型
		var sliceColumn []interface{}
		elemt := rt.Elem() //获取切片元素类型
		for i := 0; i < rv.Len(); i++ {
			inxv := rv.Index(i)
			if elemt.Kind() == reflect.Struct {
				for i := 0; i < elemt.NumField(); i++ {
					if elemt.Field(i).Name == key {
						strf := inxv.Field(i)
						switch strf.Kind() {
						case reflect.String:
							sliceColumn = append(sliceColumn, strf.String())
						case reflect.Float64:
							sliceColumn = append(sliceColumn, strf.Float())
						case reflect.Int, reflect.Int64:
							sliceColumn = append(sliceColumn, strf.Int())
						default:
							//do nothing
						}
					}
				}
			}
		}
		return sliceColumn
	}
	return nil
}

// 提取map中某列
func (f FunctionClass) MapColumn(mapSlice []interface{}, key string) []interface{} {
	sliceColumn := []interface{}{}
	if len(mapSlice) > 0 {
		for _, value := range mapSlice {
			if _, ok := value.(map[string]interface{})[key]; ok {
				sliceColumn = append(sliceColumn, value.(map[string]interface{})[key])
			}
		}
		return sliceColumn
	}
	return nil
}

// 提取map中某列
func (f FunctionClass) MapColumnIndex(mapSlice []interface{}, key string) map[string]interface{} {
	sliceColumn := make(map[string]interface{})
	for _, value := range mapSlice {
		value = value.(map[string]interface{})
		sliceColumn[value.(map[string]string)[key]] = value
	}
	return sliceColumn
}

// 判断是否在数组中
func (f FunctionClass) InArray(target interface{}, str_array interface{}) bool {
	switch str_array.(type) {
	case []string:
		for _, element := range str_array.([]string) {
			if target == element {
				return true
			}
		}
		return false
	case []int:
		for _, element := range str_array.([]int) {
			if target == element {
				return true
			}
		}
		return false
	case []interface{}:
		for _, element := range str_array.([]interface{}) {
			if target == element.(string) {
				return true
			}
		}
		return false
	default:
		for _, element := range str_array.([]interface{}) {
			if target == element {
				return true
			}
		}
		return false
	}
}

// 三目运算符
func (f FunctionClass) If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func (f FunctionClass) StrctToSlice(data any) map[string]interface{} {
	rType := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	ss := make(map[string]interface{}, v.NumField())
	for i, n := 0, v.NumField(); i < n; i++ { // 常见的 for 循环，支持初始化语句。
		name := rType.Field(i).Name
		ss[name] = v.Field(i)
	}
	return ss
}

// Struct2map 方法2：通过反射将struct转换成map
func (f FunctionClass) StructTomap(obj any) (data map[string]interface{}) {
	// 通过反射将结构体转换成map
	data = make(map[string]any)
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	for i := 0; i < objT.NumField(); i++ {
		fileName, ok := objT.Field(i).Tag.Lookup("to")
		is_continue := false
		switch objV.Field(i).Type().String() {
		case "int", "int8", "int16", "uint", "uintprt", "float32", "float64":
			if objV.Field(i).Interface() == nil || objV.Field(i).Interface() == 0 {
				is_continue = true
			}
		case "string":
			if objV.Field(i).Interface() == nil || objV.Field(i).Interface() == "" {
				is_continue = true
			}
		case "[]map[string]string":
			if objV.Field(i).Interface() == nil || objV.Field(i).Interface() == "" {
				is_continue = true
			}
		case "bool":
			if objV.Field(i).Interface() == nil {
				is_continue = true
			}
		default:
			if reflect.ValueOf(objV.Field(i).Interface()).IsNil() {
				is_continue = true
			}
		}
		if is_continue {
			continue
		}
		value := objV.Field(i).Interface()
		if reflect.TypeOf(value).String() == "*int" {
			value = *value.(*int)
		}
		if ok {
			data[fileName] = value
		} else {
			data[objT.Field(i).Name] = value
		}
	}
	return data
}

func (f FunctionClass) String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// 复制map
func (f FunctionClass) C(map1 map[string]interface{}) map[string]interface{} {
	mp2 := make(map[string]interface{}, len(map1))
	for k, v := range map1 {
		mp2[k] = v
	}
	return mp2
}

func (f FunctionClass) Unique(originals interface{}) (interface{}, error) {
	//数组去除
	switch slice := originals.(type) {
	case []string:
		result := make([]string, 0)
		m := make(map[string]bool) //map的值不重要
		for _, v := range originals.([]string) {
			if _, ok := m[v]; !ok {
				result = append(result, v)
				m[v] = true
			}
		}
		return result, nil
	case []int:
		result := make([]int, 0)
		m := make(map[int]bool) //map的值不重要
		for _, v := range originals.([]int) {
			if _, ok := m[v]; !ok {
				result = append(result, v)
				m[v] = true
			}
		}
		return result, nil
	case []int64:
		result := make([]int64, 0)
		m := make(map[int64]bool) //map的值不重要
		for _, v := range originals.([]int64) {
			if _, ok := m[v]; !ok {
				result = append(result, v)
				m[v] = true
			}
		}
		return result, nil
	case []interface{}:
		result := make([]interface{}, 0)
		m := make(map[interface{}]bool) //map的值不重要
		for _, v := range originals.([]interface{}) {
			if _, ok := m[v]; !ok {
				result = append(result, v)
				m[v] = true
			}
		}
		return result, nil
	default:
		err := f.Errorf("Unknown type: %T", slice)
		return nil, err
	}
}

func (f FunctionClass) UniqueMap(originals interface{}, indexVal interface{}) (interface{}, error) {
	switch slice := originals.(type) {
	case []map[string]interface{}:
		//数组去除
		m := make(map[interface{}]bool) //
		temp_rule_info := make([]map[string]interface{}, 0)
		for _, vv := range originals.([]map[string]interface{}) {
			switch indexslice := vv[indexVal.(string)].(type) {
			case string:
				if _, ok := m[vv[indexVal.(string)].(string)]; !ok {
					temp_rule_info = append(temp_rule_info, vv)
					m[vv[indexVal.(string)].(string)] = true
				}
			case int:
				if _, ok := m[vv[indexVal.(string)].(int)]; !ok {
					temp_rule_info = append(temp_rule_info, vv)
					m[vv[indexVal.(string)].(int)] = true
				}
			default:
				err := f.Errorf("Unknown type: %T", indexslice)
				return nil, err
			}

		}
		return temp_rule_info, nil

	case []map[int]interface{}:
		m := make(map[int]bool)
		temp_rule_info := make([]map[int]interface{}, 0)
		for _, vv := range originals.([]map[int]interface{}) {
			if _, ok := m[vv[indexVal.(int)].(int)]; !ok {
				temp_rule_info = append(temp_rule_info, vv)
				m[vv[indexVal.(int)].(int)] = true
			}
		}
		return temp_rule_info, nil
	default:
		err := f.Errorf("Unknown type: %T", slice)
		return nil, err
	}
}

func (f FunctionClass) MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}
func (f FunctionClass) Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &sliceError{msg}
}
func (f FunctionClass) Isipv6(ips string) bool {
	ip := net.ParseIP(ips)
	return ip != nil && strings.Contains(ips, ":")
}

func (f FunctionClass) GetIpType(ips string) int {
	// 解析 IP 地址
	ip := net.ParseIP(ips)

	// 判断 IP 类型
	if ip.To4() != nil {
		// IPv4 类型
		return 0
	} else if ip.To16() != nil {
		// IPv6 类型
		return 1
	} else {
		// 其他类型，返回错误值或者默认值
		return 2
	}
}

func (f FunctionClass) MaxNum(arr []int) (max int, maxIndex int) {
	max = arr[0] //假设数组的第一位为最大值
	//常规循环，找出最大值
	for i := 0; i < len(arr); i++ {
		if max < arr[i] {
			max = arr[i]
			maxIndex = i
		}
	}
	return max, maxIndex
}

func (f FunctionClass) GoScannerPath(path string) string {
	basePath := Constant.wangxinAbsoulePath
	return filepath.Join(basePath, "goscanner", path)
}

func (f FunctionClass) SliceInterfaceToInt(inters []interface{}) []int {
	ints := []int{}
	for _, inter := range inters {
		ints = append(ints, int(inter.(float64)))
	}
	return ints
}

func (f FunctionClass) ParseIpAndPortByHost(host string) (string, int) {
	result, err := url.Parse(host)
	if err != nil {
		return "", 0
	}

	if result.Port() != "" {
		port, _ := strconv.Atoi(result.Port())
		return result.Hostname(), port
	}

	if result.Scheme == "https" {
		return result.Hostname(), 443
	}

	return result.Hostname(), 80
}

func (f FunctionClass) ReadTxt(file string) (map[int]map[int]interface{}, error) {
	var err error
	lf, err := os.Open(file)
	if err != nil {
		return map[int]map[int]interface{}{}, err
	}

	defer lf.Close()
	scanner := bufio.NewScanner(lf)
	scanner.Split(bufio.ScanLines)
	words := make(map[int]map[int]interface{})
	i := 0
	for scanner.Scan() {
		stringArrs := strings.Fields(scanner.Text())
		if err != nil {
			return map[int]map[int]interface{}{}, err
		}
		j := 0
		wordMap := make(map[int]interface{})
		for _, word := range stringArrs {
			wordMap[j] = word
			j++
		}
		words[i] = wordMap
		i++
	}
	return words, nil
}

func (f FunctionClass) DeleteSlice(a []string, elem string) []string {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func (f FunctionClass) Ts(originals interface{}) string {
	switch originals.(type) {
	case float64:
		return strconv.Itoa(int(originals.(float64)))
	case int:
		return strconv.Itoa(originals.(int))
	case int64:
		return strconv.Itoa(int(originals.(int64)))
	default:
		return originals.(string)
	}
}

func (f FunctionClass) Ti(originals interface{}) int {
	switch originals.(type) {
	case float64:
		return int(originals.(float64))
	case string:
		intorigin, _ := strconv.Atoi(originals.(string))
		return intorigin
	case int64:
		return int(originals.(int64))
	default:
		return originals.(int)
	}
}

func (f FunctionClass) SetProcess(progress float32) map[string]interface{} {
	processdata := make(map[string]interface{})
	listData := make(map[string]interface{})
	listData["progress"] = progress
	listData["type"] = Constant.wanganSEARCH_ASSETS
	processdata["Type"] = "user"
	processdata["Data"] = listData
	processdata["Message"] = "查询成功"
	processdata["Success"] = true
	bytesData, _ := json.Marshal(processdata)
	contentData := make(map[string]interface{})
	contentData["content"] = string(bytesData)

	posturl := databases.Conf.Section("goroutinepool").Key("WEBSOCKET_URL").String()
	if "" == posturl {
		posturl = "http://127.0.0.1:8444/api/v2/sendtoall"
	}
	returnData := f.XPost(posturl, contentData)
	return returnData
}

func (f FunctionClass) XPost(posturl string, queryData map[string]interface{}) map[string]interface{} {
	//发送post请求
	queryParams := url.Values{}
	for kk, vv := range queryData {
		queryParams.Set(kk, vv.(string))
	}
	requestBody := queryParams.Encode()
	request, err := http.NewRequest("POST", posturl, strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(request) //Do 方法发送请求，返回 HTTP 回复
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	returndata := make(map[string]interface{})
	_ = json.Unmarshal(respBytes, &returndata)
	return returndata
}

func (f FunctionClass) Xget(posturl string, queryData map[string]interface{}) map[string]interface{} {
	//如果参数中有中文参数,这个方法会进行URLEncode
	queryParams := url.Values{}
	Url, _ := url.Parse(posturl)
	for kk, vv := range queryData {
		queryParams.Set(kk, vv.(string))
	}
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
	return data
}
