/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/5/15 10:20
 * @desc: about the role of class.
 */

package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiURL = "https://api.uomg.com/api/long2dwz?%s"
)

type ShortResponse struct {
	Code     int    `json:"code"`
	ShortURL string `json:"ae_url"`
}

// GetShortURL https://www.free-api.com/doc/300
func GetShortURL(longURL string) (string, error) {

	params := url.Values{}
	params.Add("url", longURL)
	params.Add("dwzapi", "urlcn")

	resp, err := http.Get(fmt.Sprintf(apiURL, params.Encode()))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	shortResponse := new(ShortResponse)
	err = json.Unmarshal(body, shortResponse)

	return shortResponse.ShortURL, err
}
