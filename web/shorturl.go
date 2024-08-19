/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/8/16 09:26
 * @desc: about the role of class.
 */

package web

import (
	"encoding/json"
	"errors"
	"github.com/AbnerEarl/goutils/esc"
	"github.com/AbnerEarl/goutils/redisc"
	"math"
	"strings"
	"time"
)

var (
	SHORT_URL_INFO_INDEX          = "short_url_info_index"
	SHORT_URL_INFO_INDEX_DOC_TYPE = "short_url_info_index_doc_type"

	BASE_STR  = "RScdTU567rstu89ABCDijnopEFGefwxHIJ01klm2ghqv34KLMNOPQVWXYZabyz"
	BASE_CHAR = []string{"R", "S", "c", "d", "T", "U", "5", "6", "7", "r", "s", "t", "u", "8", "9", "A", "B", "C", "D", "i", "j", "n", "o", "p", "E", "F", "G", "e", "f", "w", "x", "H", "I", "J", "0", "1", "k", "l", "m", "2", "g", "h", "q", "v", "3", "4", "K", "L", "M", "N", "O", "P", "Q", "V", "W", "X", "Y", "Z", "a", "b", "y", "z"}
)

type ShortUrl struct {
	BaseUrl  string                  // The url or domain name of the server, e.g. http://127.0.0.1:8080/api/d/
	Rediscli *redisc.RedisCli        // redis Standalone Client or Cluster Either one can be initialized
	Redisclu *redisc.RedisClusterCli // redis Standalone Client or Cluster Either one can be initialized
	Esc7     *esc.EsClient7          // Just choose any one of the different versions of ES to initialize
	Esc8     *esc.EsClient8          // Just choose any one of the different versions of ES to initialize
}

func (su *ShortUrl) GenShortUrl(longUrl string) (string, error) {
	var id uint64
	if su.Rediscli != nil {
		id = GenAutoId(su.Rediscli)
	} else if su.Redisclu != nil {
		id = GenAutoIdByClu(su.Redisclu)
	} else {
		return "", errors.New("the redis client is not initialized")
	}
	shortCode, err := Int2Char(id, 62)
	if err != nil {
		return "", err
	}
	shortUrl := su.BaseUrl + shortCode

	doc := map[string]interface{}{
		"id":       id,
		"doc_type": SHORT_URL_INFO_INDEX_DOC_TYPE,
		"body": map[string]interface{}{
			"short_code": shortCode,
			"long_url":   longUrl,
			"timestamp":  time.Now(),
		},
	}
	if su.Esc7 != nil {
		exist, _ := su.Esc7.ExistsIndexBool(SHORT_URL_INFO_INDEX)
		if !exist {
			_, err = su.Esc7.CreateIndex(SHORT_URL_INFO_INDEX)
			if err != nil {
				return "", err
			}
		}

		_, err = su.Esc7.CreateDocMap(SHORT_URL_INFO_INDEX, doc)
		if err != nil {
			return "", err
		}
	} else if su.Esc8 != nil {
		exist, _ := su.Esc8.ExistsIndexBool(SHORT_URL_INFO_INDEX)
		if !exist {
			_, err = su.Esc8.CreateIndex(SHORT_URL_INFO_INDEX)
			if err != nil {
				return "", err
			}
		}
		_, err = su.Esc8.CreateDocMap(SHORT_URL_INFO_INDEX, doc)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("the elasticsearch client is not initialized")
	}
	return shortUrl, nil
}

func (su *ShortUrl) ConvLongUrl(shortUrl string) (string, error) {
	if strings.Contains(shortUrl, "?") {
		shortUrl = strings.Split(shortUrl, "?")[0]
	}
	shortCode := shortUrl
	if strings.Contains(shortUrl, "/") {
		ss := strings.Split(shortUrl, "/")
		for i := len(ss) - 1; i >= 0; i-- {
			if len(ss[i]) > 2 {
				shortCode = ss[i]
				break
			}
		}
	}
	id, err := Char2Int(shortCode, 62)
	if err != nil {
		return "", err
	}
	params := map[string]interface{}{
		"doc_type": SHORT_URL_INFO_INDEX_DOC_TYPE,
		"id":       id,
	}
	longUrl := ""
	if su.Esc7 != nil {
		r, err := su.Esc7.RetrieveDocMap(SHORT_URL_INFO_INDEX, params)
		if err != nil {
			return "", err
		}
		if r.IsError() {
			return "", err
		}
		defer r.Body.Close()
		var searchHits map[string]interface{}
		if err = json.NewDecoder(r.Body).Decode(&searchHits); err != nil {
			return "", err
		}

		if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
			if hitsList, ok := hits["hits"].([]interface{}); ok {
				for _, hit := range hitsList {
					source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
					if !ok {
						continue
					} else {
						url, ok := source["long_url"]
						if ok {
							longUrl = url.(string)
							break
						}
					}
				}
			}
		}

	} else if su.Esc8 != nil {
		r, err := su.Esc8.RetrieveDocMap(SHORT_URL_INFO_INDEX, params)
		if err != nil {
			return "", err
		}
		if r.IsError() {
			return "", err
		}
		defer r.Body.Close()
		var searchHits map[string]interface{}
		if err = json.NewDecoder(r.Body).Decode(&searchHits); err != nil {
			return "", err
		}
		if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
			if hitsList, ok := hits["hits"].([]interface{}); ok {
				for _, hit := range hitsList {
					source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
					if !ok {
						continue
					} else {
						url, ok := source["long_url"]
						if ok {
							longUrl = url.(string)
							break
						}
					}
				}
			}
		}
	} else {
		return "", errors.New("the elasticsearch client is not initialized")
	}
	return longUrl, nil
}

func Char2Int(code string, base uint64) (uint64, error) {
	ss := strings.Split(code, "")
	var result uint64
	for i := 0; i < len(ss); i++ {
		n := strings.Index(BASE_STR, ss[i])
		if n < 0 {
			return 0, errors.New("url error, unknown character")
		}
		result = result + uint64(n)*uint64(math.Pow(float64(base), float64(i)))
	}
	return result, nil
}

func Int2Char(number, base uint64) (string, error) {
	if base > uint64(len(BASE_CHAR)) {
		return "", errors.New("the value of base exceeds the maximum value")
	}
	var result []string
	for number > 0 {
		index := number % base
		number = number / base
		result = append(result, BASE_CHAR[index])
	}
	return strings.Join(result, ""), nil
}
