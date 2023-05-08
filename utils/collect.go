package utils

import "encoding/json"

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
