package Request

import (
	"github.com/go-playground/validator/v10"
	"log"
	"net"
	"regexp"
)

var IpValidrule validator.Func = func(fl validator.FieldLevel) bool {
	inputData, ok := fl.Field().Interface().(string)
	if ok && inputData != "" {
		ip := net.ParseIP(inputData)
		if ip == nil {
			return false
		}
	}
	return true
}

var IpDataValidrule validator.Func = func(fl validator.FieldLevel) bool {
	inputData, ok := fl.Field().Interface().([]string)
	if ok {
		for _, inputIp := range inputData {
			addr, err := net.ResolveIPAddr("ip", inputIp)
			if err != nil {
				log.Println("Error", err.Error())
				return false
			}
			parseIp := addr.String()
			ip := net.ParseIP(parseIp)
			if ip == nil {
				return false
			}
		}
	}
	return true
}

var KeywordVal validator.Func = func(fl validator.FieldLevel) bool {
	keyword, ok := fl.Field().Interface().(string)
	if ok && keyword != "" {
		res, err := regexp.MatchString("^\\s*[\\x{4e00}-\\x{9fa5}A-Za-z0-9=\"*.]+\\s*$", keyword)
		if res == false || err != nil {
			log.Print("Error", err)
			return false
		}
	}
	return true
}

var TimeRangeVal validator.Func = func(fl validator.FieldLevel) bool {
	timerange, ok := fl.Field().Interface().([]string)
	if ok {
		for _, key := range timerange {
			res, err := regexp.MatchString("^([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01])) (([01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d)$", key)
			if res == false || err != nil {
				log.Println("Error", "日期时间输入格式出错")
				return false
			}
		}
	}
	return true
}

var IdsVal validator.Func = func(fl validator.FieldLevel) bool {
	ids, ok := fl.Field().Interface().([]int)
	if !ok {
		log.Println("Error", "参数错误")
		return false
	}
	for _, id := range ids {
		if id < 1 {
			log.Println("Error", "参数错误")
			return false
		}
	}
	return true
}
