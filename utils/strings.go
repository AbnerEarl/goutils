package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"unicode"
	"unsafe"
)

func StrLength(str string) int {
	count := 0
	for _, c := range str {
		if unicode.Is(unicode.Han, c) {
			count++
		} else {
			count++
		}
	}
	return count
}

func ChEnLength(str string) int {
	count := 0
	for _, c := range str {
		if unicode.Is(unicode.Han, c) {
			count += 2
		} else {
			count++
		}
	}
	return count
}

func CheckPasswordRule(password string) bool {
	if len(password) < 11 {
		return false
	}
	var digit, upper, lower, special bool
	for _, c := range password {
		if unicode.IsDigit(c) {
			digit = true
		} else if unicode.IsUpper(c) {
			upper = true
		} else if unicode.IsLower(c) {
			lower = true
		} else {
			special = true
		}
	}
	if !digit || !upper || !lower || !special {
		return false
	}
	return true
}

func SubnetMatch(subnet string) (string, string, error) {
	ip, ipv4Net, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", "", err
	}
	return ipv4Net.String(), ip.String(), err
}

func CheckIP(ip string) string {
	result := net.ParseIP(ip)
	if result != nil {
		return result.String()
	}
	return ""
}

func Byte2Any(b []byte, t reflect.Type) interface{} {
	data := Byte2Str(b)
	switch t.Kind() {
	case reflect.Int:
		res, _ := strconv.Atoi(data)
		return res
	case reflect.Int64:
		res, _ := strconv.ParseUint(data, 10, 64)
		return res
	case reflect.Uint:
		res, _ := strconv.ParseUint(data, 10, 64)
		return res
	case reflect.Uint64:
		res, _ := strconv.ParseUint(data, 10, 64)
		return res
	case reflect.Float64:
		res, _ := strconv.ParseFloat(data, 64)
		return res
	case reflect.Bool:
		res, _ := strconv.ParseBool(data)
		return res
	default:
		return data
	}
}

func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Str2Uint64(s string) uint64 {
	res, _ := strconv.ParseUint(s, 10, 64)
	return res
}

func Str2Float(s string) float64 {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	//return math.Trunc(res) * 1e-2
	res, _ = NewFromFloat(res).Round(2).Float64()
	return res
}

func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

func IsEmpty(data interface{}) bool {
	if data == nil || data == "" {
		return true
	}
	if reflect.TypeOf(data).Kind() == reflect.String && len(fmt.Sprint(data)) < 1 {
		return true
	} else if reflect.TypeOf(data).Kind() == reflect.Slice {
		arr := data.([]interface{})
		if len(arr) < 1 {
			return true
		}
	} else if reflect.TypeOf(data).Kind() == reflect.Struct {
		ins1 := reflect.New(reflect.TypeOf(data)).Interface()
		bytes, _ := json.Marshal(data)
		json.Unmarshal(bytes, &ins1)
		ins2 := reflect.New(reflect.TypeOf(data)).Interface()
		return reflect.DeepEqual(ins1, ins2)
	} else if reflect.TypeOf(data).Kind() == reflect.Map {
		bytes, err := json.Marshal(data)
		if err != nil {
			return true
		}
		dataMap := map[string]interface{}{}
		json.Unmarshal(bytes, &dataMap)
		if len(dataMap) == 0 {
			return true
		}
	}
	return false
}
