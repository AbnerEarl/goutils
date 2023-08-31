package utils

import (
	"encoding/json"
	"reflect"
	"strings"
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

func Struct2Map(model interface{}, tagName string) map[string]interface{} {
	/**
	 * @author: yangchangjia
	 * @email 1320259466@qq.com
	 * @date: 2023/8/29 6:49 PM
	 * @desc: about the role of function.
	 * @param model: an struct or struct pointer.
	 * @param tagName: field conversion flag, default value is "", can be set, "json",etc.
	 * @return map
	 */
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if !strings.Contains(v.Field(i).String(), "time.Time") && (v.Field(i).Kind() == reflect.Struct || v.Field(i).Kind() == reflect.Ptr) {
			items := Struct2Map(v.Field(i).Interface(), tagName)
			for k, m := range items {
				data[k] = m
			}
		} else {
			if !IsEmpty(tagName) {
				tag := t.Field(i).Tag.Get(tagName)
				if tag == "-" {
					continue
				}
				if !IsEmpty(tag) {
					data[tag] = v.Field(i).Interface()
					continue
				}
			}
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return data
}

func Struct2MapNoZero(model interface{}, tagName string) map[string]interface{} {
	/**
	 * @author: yangchangjia
	 * @email 1320259466@qq.com
	 * @date: 2023/8/29 6:31 PM
	 * @desc: about the role of function.
	 * @param model: an struct or struct pointer.
	 * @param tagName: field conversion flag, default value is "", can be set, "json",etc.
	 * @return map
	 */

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if !strings.Contains(v.Field(i).String(), "time.Time") && (v.Field(i).Kind() == reflect.Struct || v.Field(i).Kind() == reflect.Ptr) {
			items := Struct2MapNoZero(v.Field(i).Interface(), tagName)
			for k, m := range items {
				data[k] = m
			}
		} else if !v.Field(i).IsZero() {
			if !IsEmpty(tagName) {
				tag := t.Field(i).Tag.Get(tagName)
				if tag == "-" {
					continue
				}
				if !IsEmpty(tag) {
					data[tag] = v.Field(i).Interface()
					continue
				}
			}
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return data
}
