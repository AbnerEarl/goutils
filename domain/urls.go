/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/16 17:06
 * @desc: about the role of class.
 */

package domain

import (
	"github.com/AbnerEarl/goutils/httpc"
	"regexp"
	"strings"
	"time"
)

func GetTopUrl(url string) string {
	patter := `^((http://)|(https://))?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}(/)`
	reg := regexp.MustCompile(patter)
	domain := reg.FindString(url)
	return domain
}

func GetAllUrl(url string) []string {
	result, err := httpc.RequestString(url, "GET", nil, nil, 20*time.Second)
	if err != nil {
		return nil
	}
	return ExtraUrls(result)
}
func GetAllUrlByProxy(proxyUrl, requestUrl string) []string {
	result, err := httpc.RequestStringByProxy(proxyUrl, requestUrl, "GET", nil, nil, 20*time.Second)
	if err != nil {
		return nil
	}
	return ExtraUrls(result)
}

func ExtraUrls(html string) []string {
	var urls []string
	// 正则表达式匹配标签中的href属性
	regRuler := `[^>]*?\s+?href="([^"]*)"`
	// 正则调用规则
	reg := regexp.MustCompile(regRuler)
	if matches := reg.FindAllStringSubmatch(html, -1); matches != nil {
		for _, match := range matches {
			if len(match) > 1 {
				url := match[1]
				// 过滤掉非绝对路径的URL
				if !strings.HasPrefix(url, "/") && strings.HasPrefix(url, "http") {
					urls = append(urls, url)
				}
			}
		}
	}
	return urls
}

func GetAllDomain(url string) []string {
	allUrl := GetAllUrl(url)
	urlMap := map[string]struct{}{}
	for _, u := range allUrl {
		domain, err := GetDomain(u)
		if err != nil {
			continue
		}
		urlMap[domain] = struct{}{}
	}
	var urls []string
	for k, _ := range urlMap {
		urls = append(urls, k)
	}
	return urls
}

func GetAllDomainByProxy(proxyUrl, requestUrl string) []string {
	allUrl := GetAllUrlByProxy(proxyUrl, requestUrl)
	urlMap := map[string]struct{}{}
	for _, u := range allUrl {
		domain, err := GetDomain(u)
		if err != nil {
			continue
		}
		urlMap[domain] = struct{}{}
	}
	var urls []string
	for k, _ := range urlMap {
		urls = append(urls, k)
	}
	return urls
}

func GetDomain(url string) (string, error) {
	url = strings.Split(url, "//")[1]
	url = strings.Split(url, "/")[0]
	return Domain(url)
}

func GetDomainInfo(url string) (*DomainName, error) {
	url = strings.Split(url, "//")[1]
	url = strings.Split(url, "/")[0]
	return Parse(url)
}

func GetAllDomainInfo(url string) []*DomainName {
	allUrl := GetAllUrl(url)
	urlMap := map[string]*DomainName{}
	for _, u := range allUrl {
		domain, err := GetDomainInfo(u)
		if err != nil {
			continue
		}
		urlMap[domain.String()] = domain
	}
	var domains []*DomainName
	for _, v := range urlMap {
		domains = append(domains, v)
	}
	return domains
}

func GetAllDomainInfoByProxy(proxyUrl, requestUrl string) []*DomainName {
	allUrl := GetAllUrlByProxy(proxyUrl, requestUrl)
	urlMap := map[string]*DomainName{}
	for _, u := range allUrl {
		domain, err := GetDomainInfo(u)
		if err != nil {
			continue
		}
		urlMap[domain.String()] = domain
	}
	var domains []*DomainName
	for _, v := range urlMap {
		domains = append(domains, v)
	}
	return domains
}
