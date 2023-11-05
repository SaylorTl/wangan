package Services

import (
	"fmt"
	"strconv"
	"strings"
	"wangxin2.0/app/Utils"
)

var GarBleService garbleservice

type garbleservice struct {
}

func (b garbleservice) _idToString(num int) string {
	convertStr := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	strList := strings.Split(convertStr, "")
	str := ""
	for num != 0 {
		tmp := num % 32
		str += strList[tmp]
		num = int(num / 32)
	}
	return str
}

func (b garbleservice) IdToString(id int) string {
	INIT_NUM := 123456789
	newId := id + INIT_NUM
	newStr := fmt.Sprintf("%010d", newId)
	num1 := string(newStr[0]) + string(newStr[2]) + string(newStr[6]) + string(newStr[9])
	num2, _ := strconv.Atoi(num1)
	num3 := string(newStr[1]) + string(newStr[3]) + string(newStr[4]) + string(newStr[5]) + string(newStr[7]) + string(newStr[8])
	num4, _ := strconv.Atoi(num3)
	str1 := b._idToString(num2)
	str1 = Utils.F.Reverse(str1)
	str2 := b._idToString(num4)
	str2 = Utils.F.Reverse(str2)
	newStr1 := Utils.F.StrPad(str1, 3, "U", "RIGHT")
	newStr2 := Utils.F.StrPad(str2, 4, "L", "RIGHT")
	encodeStr := newStr1 + newStr2
	return encodeStr
}
