package utils

import (
	"encoding/json"
	"reflect"
)

func RemoveDuplicateElement(s []interface{}) []interface{} {
	result := make([]interface{}, 0, len(s))
	temp := map[string]struct{}{}
	for _, item := range s {
		bytes, _ := json.Marshal(item)
		key := string(bytes)
		if _, ok := temp[key]; !ok {
			temp[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Append(slice interface{}, elements interface{}) []interface{} {
	var result []interface{}
	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		bys, _ := json.Marshal(slice)
		json.Unmarshal(bys, &result)

	} else {
		return result
	}
	if reflect.TypeOf(elements).Kind() == reflect.Slice {
		bys, _ := json.Marshal(elements)
		var arr []interface{}
		json.Unmarshal(bys, &arr)
		for _, el := range arr {
			result = append(result, el)
		}
	} else {
		result = append(result, elements)
	}
	return result
}
