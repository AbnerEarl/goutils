/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/8/23 4:14 PM
 * @desc: about the role of class.
 */

package esc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AbnerEarl/goutils/utils"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"os"
	"strings"
)

type EsClient7 struct {
	*elasticsearch7.Client
}

type Response7 struct {
	*esapi.Response
}

func InitV7EsClientCloud(cloudId, apiKey string) (*EsClient7, error) {
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		CloudID: cloudId,
		APIKey:  apiKey,
	})
	return &EsClient7{es7}, err
}

func InitV7EsClientCACert(addrs []string, username, password, caCertFile string) (*EsClient7, error) {
	cert, _ := os.ReadFile(caCertFile)
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: addrs, // []string{"http://localhost:9200"}
		Username:  username,
		Password:  password,
		CACert:    cert,
	})
	return &EsClient7{es7}, err
}

func InitV7EsClientFingerprint(addrs []string, username, password, certFingerprint string) (*EsClient7, error) {
	/**
	The certificate fingerprint can be calculated using openssl x509 with the certificate file:
	openssl x509 -fingerprint -sha256 -noout -in /path/to/http_ca.crt
	If you don’t have access to the generated CA file from Elasticsearch you can use the following script to output
	the root CA fingerprint of the Elasticsearch instance with openssl s_client:

	# Replace the values of 'localhost' and '9200' to the
	# corresponding host and port values for the cluster.
	openssl s_client -connect localhost:9200 -servername localhost -showcerts </dev/null 2>/dev/null \
	  | openssl x509 -fingerprint -sha256 -noout -in /dev/stdin

	The output of openssl x509 will look something like this:
	SHA256 Fingerprint=A5:2D:D9:35:11:E8:C6:04:5E:21:F1:66:54:B7:7C:9E:E0:F3:4A:EA:26:D9:F4:03:20:B5:31:C4:74:67:62:28
	 *
	*/
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses:              addrs, // []string{"http://localhost:9200"}
		Username:               username,
		Password:               password,
		CertificateFingerprint: certFingerprint,
	})
	return &EsClient7{es7}, err
}

func InitV7EsClientToken(addrs []string, token string) (*EsClient7, error) {
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses:    addrs, // []string{"http://localhost:9200"}
		ServiceToken: token,
	})
	return &EsClient7{es7}, err
}

func InitV7EsClient(addrs []string) (*EsClient7, error) {
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: addrs, // []string{"http://localhost:9200"}
	})
	return &EsClient7{es7}, err
}

func InitV7EsClientBasic(addrs []string, username, password string) (*EsClient7, error) {
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: addrs, // []string{"http://localhost:9200"}
		Username:  username,
		Password:  password,
	})
	return &EsClient7{es7}, err
}

func (esc *EsClient7) CreateIndex(name string) (*Response7, error) {
	res, err := esc.Indices.Create(name)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) CreateIndexBodyMap(name string, body map[string]interface{}) (*Response7, error) {
	bys, _ := json.Marshal(body)
	res, err := esapi.IndicesCreateRequest{Index: name, Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) CreateIndexBodyString(name, body string) (*Response7, error) {
	bys, _ := json.Marshal(body)
	res, err := esapi.IndicesCreateRequest{Index: name, Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) ExistsIndex(name string) (*Response7, error) {
	res, err := esc.Indices.Exists([]string{name})
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveIndex(name string) (*Response7, error) {
	res, err := esc.Indices.Get([]string{name})
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) ExistsIndexBool(name string) (bool, error) {
	res, err := esc.Indices.Exists([]string{name})
	if err != nil {
		return false, err
	}
	if !res.IsError() {
		return true, nil
	}
	return false, nil
}

func (esc *EsClient7) RetrieveIndexBool(name string) (bool, error) {
	res, err := esc.Indices.Get([]string{name})
	if err != nil {
		return false, err
	}
	if !res.IsError() {
		return true, nil
	}
	return false, nil
}

func (esc *EsClient7) DeleteIndex(name string) (*Response7, error) {
	res, err := esc.Indices.Delete([]string{name})
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) CreateDoc(index, body string) (*Response7, error) {
	bys, _ := json.Marshal(body)
	res, err := esc.Index(index, bytes.NewReader(bys))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) CreateDocMap(index string, body map[string]interface{}) (*Response7, error) {
	bys, _ := json.Marshal(body)
	res, err := esc.Index(index, bytes.NewReader(bys))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) CreateDocBatchString(index string, docs []string) (*Response7, error) {
	bys, _ := json.Marshal(docs)
	res, err := esapi.BulkRequest{Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) CreateDocBatchMap(index string, docs []map[string]interface{}) (*Response7, error) {
	bys, _ := json.Marshal(docs)
	res, err := esapi.BulkRequest{Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDoc(index, id string) (*Response7, error) {
	res, err := esc.Get(index, id)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDoc2Map(index, id string) ([]map[string]interface{}, error) {
	res, err := esc.Get(index, id)
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDoc2Byte(index, id string) ([]byte, error) {
	res, err := esc.Get(index, id)
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err := json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocSql(sql string) (*Response7, error) {
	query := map[string]interface{}{"query": sql}
	bys, _ := json.Marshal(query)
	res, err := esapi.SQLQueryRequest{Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocSql2Map(sql string) ([]map[string]interface{}, error) {
	query := map[string]interface{}{"query": sql}
	bys, _ := json.Marshal(query)
	res, err := esapi.SQLQueryRequest{Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocSql2Byte(sql string) ([]byte, error) {
	query := map[string]interface{}{"query": sql}
	bys, _ := json.Marshal(query)
	res, err := esapi.SQLQueryRequest{Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err = json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocQuery(index, query string) (*Response7, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocQuery2Map(index, query string) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocQuery2Byte(index, query string) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err := json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocCountQuery(index, query string) (*Response7, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocCountQueryNum(index, query string) (uint64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return 0, err
	}
	var countResult map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0, err
	}
	count := uint64(countResult["count"].(float64))
	return count, nil
}

func (esc *EsClient7) RetrieveDocMatch(index string, params map[string]interface{}) (*Response7, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": params,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocMatch2Map(index string, params map[string]interface{}) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": params,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocMatch2Byte(index string, params map[string]interface{}) ([]byte, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": params,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err := json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocCountMatch(index string, params map[string]interface{}) (*Response7, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": params,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocCountMatchNum(index string, params map[string]interface{}) (uint64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": params,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return 0, err
	}
	var countResult map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0, err
	}
	count := uint64(countResult["count"].(float64))
	return count, nil
}

func (esc *EsClient7) RetrieveDocMap(index string, params map[string]interface{}) (*Response7, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, err
	}

	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocMap2Map(index string, params map[string]interface{}) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, err
	}

	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocMap2Byte(index string, params map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, err
	}

	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err := json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocCountMap(index string, params map[string]interface{}) (*Response7, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocCountMapNum(index string, params map[string]interface{}) (uint64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return 0, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return 0, err
	}
	var countResult map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0, err
	}
	count := uint64(countResult["count"].(float64))
	return count, nil
}

func (esc *EsClient7) RetrieveDocCountModel(index string, model interface{}, tagName string) (*Response7, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocCountModelNum(index string, model interface{}, tagName string) (uint64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0, err
	}
	res, err := esapi.CountRequest{Index: []string{index}, Body: &buf}.Do(context.Background(), esc)
	if err != nil {
		return 0, err
	}
	var countResult map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0, err
	}
	count := uint64(countResult["count"].(float64))
	return count, nil
}

func (esc *EsClient7) RetrieveDocModel(index string, model interface{}, tagName string) (*Response7, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocModel2Map(index string, model interface{}, tagName string) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocModel2Byte(index string, model interface{}, tagName string) ([]byte, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err := json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocMapList(pageSize, pageNo int, params map[string]interface{}, order, index string) (*Response7, error) {
	if pageSize < 1 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo < 1 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	query := map[string]interface{}{
		"from": offset,
		"size": pageSize,
	}
	orders := strings.Split(utils.StrTrim(order), " ")
	if len(orders) == 2 {
		sort := map[string]interface{}{
			orders[0]: map[string]interface{}{
				"order": orders[1],
			},
		}
		query["sort"] = sort
	}
	if params != nil {
		match := map[string]interface{}{
			"match": params,
		}
		query["query"] = match
	}
	bys, _ := json.Marshal(query)
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(bytes.NewReader(bys)))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocMapList2Map(pageSize, pageNo int, params map[string]interface{}, order, index string) ([]map[string]interface{}, error) {
	if pageSize < 1 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo < 1 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	query := map[string]interface{}{
		"from": offset,
		"size": pageSize,
	}
	orders := strings.Split(utils.StrTrim(order), " ")
	if len(orders) == 2 {
		sort := map[string]interface{}{
			orders[0]: map[string]interface{}{
				"order": orders[1],
			},
		}
		query["sort"] = sort
	}
	if params != nil {
		match := map[string]interface{}{
			"match": params,
		}
		query["query"] = match
	}
	bys, _ := json.Marshal(query)
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(bytes.NewReader(bys)))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocMapList2Byte(pageSize, pageNo int, params map[string]interface{}, order, index string) ([]byte, error) {
	if pageSize < 1 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo < 1 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	query := map[string]interface{}{
		"from": offset,
		"size": pageSize,
	}
	orders := strings.Split(utils.StrTrim(order), " ")
	if len(orders) == 2 {
		sort := map[string]interface{}{
			orders[0]: map[string]interface{}{
				"order": orders[1],
			},
		}
		query["sort"] = sort
	}
	if params != nil {
		match := map[string]interface{}{
			"match": params,
		}
		query["query"] = match
	}
	bys, _ := json.Marshal(query)
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(bytes.NewReader(bys)))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err = json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) RetrieveDocModelList(pageSize, pageNo int, model interface{}, tagName, order, index string) (*Response7, error) {
	if pageSize < 1 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo < 1 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	query := map[string]interface{}{
		"from": offset,
		"size": pageSize,
	}
	orders := strings.Split(utils.StrTrim(order), " ")
	if len(orders) == 2 {
		sort := map[string]interface{}{
			orders[0]: map[string]interface{}{
				"order": orders[1],
			},
		}
		query["sort"] = sort
	}
	if model != nil {
		match := map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		}
		query["query"] = match
	}
	bys, _ := json.Marshal(query)
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(bytes.NewReader(bys)))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) RetrieveDocModelList2Map(pageSize, pageNo int, model interface{}, tagName, order, index string) ([]map[string]interface{}, error) {
	if pageSize < 1 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo < 1 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	query := map[string]interface{}{
		"from": offset,
		"size": pageSize,
	}
	orders := strings.Split(utils.StrTrim(order), " ")
	if len(orders) == 2 {
		sort := map[string]interface{}{
			orders[0]: map[string]interface{}{
				"order": orders[1],
			},
		}
		query["sort"] = sort
	}
	if model != nil {
		match := map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		}
		query["query"] = match
	}
	bys, _ := json.Marshal(query)
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(bytes.NewReader(bys)))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	return result, nil
}

func (esc *EsClient7) RetrieveDocModelList2Byte(pageSize, pageNo int, model interface{}, tagName, order, index string) ([]byte, error) {
	if pageSize < 1 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo < 1 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	query := map[string]interface{}{
		"from": offset,
		"size": pageSize,
	}
	orders := strings.Split(utils.StrTrim(order), " ")
	if len(orders) == 2 {
		sort := map[string]interface{}{
			orders[0]: map[string]interface{}{
				"order": orders[1],
			},
		}
		query["sort"] = sort
	}
	if model != nil {
		match := map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		}
		query["query"] = match
	}
	bys, _ := json.Marshal(query)
	res, err := esc.Search(esc.Search.WithIndex(index), esc.Search.WithBody(bytes.NewReader(bys)))
	if err != nil {
		return nil, err
	}
	var searchHits map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&searchHits); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if hits, ok := searchHits["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
				if !ok {
					continue
				} else {
					result = append(result, source)
				}
			}
		}
	}
	bys, err = json.Marshal(result)
	return bys, err
}

func (esc *EsClient7) UpdateDoc(index, id string, body interface{}) (*Response7, error) {
	bys, _ := json.Marshal(body)
	res, err := esc.Update(index, id, bytes.NewReader(bys))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) UpdateDocMap(index, id string, params map[string]interface{}) (*Response7, error) {
	bys, _ := json.Marshal(params)
	res, err := esc.Update(index, id, bytes.NewReader(bys))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) UpdateDocModel(index, id string, model interface{}, tagName string) (*Response7, error) {
	bys, _ := json.Marshal(utils.Struct2MapNoZero(model, tagName))
	res, err := esc.Update(index, id, bytes.NewReader(bys))
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) UpdateDocQuery(index, query, script string) (*Response7, error) {
	body := map[string]interface{}{
		"query":  query,
		"script": script,
	}
	bys, _ := json.Marshal(body)
	res, err := esapi.UpdateByQueryRequest{Index: []string{index}, Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) UpdateDocQueryMap(index string, query, params map[string]interface{}) (*Response7, error) {
	uds := ""
	for k, _ := range params {
		uds += fmt.Sprintf("ctx._source.%s=params.%s;", k, k)
	}
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query,
		},
		"script": map[string]interface{}{
			"source": uds,
			"params": params,
			"lang":   "painless",
		},
	}
	bys, _ := json.Marshal(body)
	res, err := esapi.UpdateByQueryRequest{Index: []string{index}, Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) UpdateDocQueryModel(index string, whereModel, updateModel map[string]interface{}, tagName string) (*Response7, error) {
	query := utils.Struct2MapNoZero(whereModel, tagName)
	params := utils.Struct2MapNoZero(updateModel, tagName)
	uds := ""
	for k, _ := range params {
		uds += fmt.Sprintf("ctx._source.%s=params.%s;", k, k)
	}
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query,
		},
		"script": map[string]interface{}{
			"source": uds,
			"params": params,
			"lang":   "painless",
		},
	}
	bys, _ := json.Marshal(body)
	res, err := esapi.UpdateByQueryRequest{Index: []string{index}, Body: bytes.NewReader(bys)}.Do(context.Background(), esc)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) DeleteDoc(index, id string) (*Response7, error) {
	res, err := esc.Delete(index, id)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) DeleteDocMap(index string, params map[string]interface{}) (*Response7, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": params,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.DeleteByQuery([]string{index}, &buf)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}

func (esc *EsClient7) DeleteDocModel(index string, model interface{}, tagName string) (*Response7, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": utils.Struct2MapNoZero(model, tagName),
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := esc.DeleteByQuery([]string{index}, &buf)
	if err != nil {
		return nil, err
	}
	return &Response7{res}, nil
}
