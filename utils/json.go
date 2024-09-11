/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/9/10 16:35
 * @desc: about the role of class.
 */

package utils

import (
	"encoding/json"
	"os"
	"strconv"
)

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func Yaml2JsonForString(content string) (map[string]interface{}, error) {
	var body interface{}
	if err := Unmarshal([]byte(content), &body); err != nil {
		return nil, err
	}

	body = convert(body)

	bys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(bys, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Yaml2JsonForByte(content []byte) (map[string]interface{}, error) {
	var body interface{}
	if err := Unmarshal(content, &body); err != nil {
		return nil, err
	}

	body = convert(body)

	bys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(bys, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Yaml2JsonForFile(filePath string) (map[string]interface{}, error) {
	var body interface{}
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err = Unmarshal(buffer, &body); err != nil {
		return nil, err
	}

	body = convert(body)

	bys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(bys, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Json2MapForFile(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := map[string]interface{}{}
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Json2ListMapForFile(filePath string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []map[string]interface{}
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Json2MapForFileConv(filePath string) (map[string]interface{}, error) {
	bys, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	jsonStr, err := strconv.Unquote("\"" + string(bys) + "\"")
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Json2ListMapForFileConv(filePath string) ([]map[string]interface{}, error) {
	bys, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	jsonStr, err := strconv.Unquote("\"" + string(bys) + "\"")
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
