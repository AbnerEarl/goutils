/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/5/6 14:18
 * @desc: about the role of class.
 */

package httpc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

/**
 * @Description: Get WinNT name at Window NT
 * @param sName
 * @return string
 */
func getWinOsNameWithWinNT(sName string) string {
	osName := "Windows"
	types := map[string]string{
		"Windows NT 11":  "Windows 11",
		"Windows NT 10":  "Windows 10",
		"Windows NT 6.3": "Windows 8",
		"Windows NT 6.2": "Windows 8",
		"Windows NT 6.1": "Windows 7",
		"Windows NT 6.0": "Windows Vista/Server 2008",
		"Windows NT 5.2": "Windows Server 2003",
		"Windows NT 5.1": "Windows XP",
		"Windows NT 5":   "Windows 2000",
		"Windows NT 4":   "Windows NT4",
	}

	for keyWord, name := range types {
		if strings.Contains(sName, keyWord) {
			osName = name
			break
		}
	}

	return osName
}

// GetOsName 获取OS名称
func GetOsName(userAgent string) string {
	osName := "Unknown"
	if userAgent == "" {
		return osName
	}

	strRe, _ := regexp.Compile("(?i:\\((.*?)\\))")
	userAgent = strRe.FindString(userAgent)

	levelNames := ":micromessenger:dart:Windows NT:Windows Mobile:Windows Phone:Windows Phone OS:Macintosh|Macintosh:Mac OS:CrOS|CrOS:iPhone OS:iPad|iPad:OS:Android:Linux:blackberry:hpwOS:Series:Symbian:PalmOS:SymbianOS:J2ME:Sailfish:Bada:MeeGo:webOS|hpwOS:Maemo:"
	var regStrArr []string
	namesArr := strings.Split(strings.Trim(levelNames, ":"), ":")
	for _, name := range namesArr {
		regStrArr = append(regStrArr, fmt.Sprintf("(%s[\\s?\\/XxSs0-9_.]+)", name))
	}
	regexpStr := fmt.Sprintf("(?i:%s)", strings.Join(regStrArr, "|"))
	nameRe, _ := regexp.Compile(regexpStr)

	names := nameRe.FindAllString(userAgent, -1)
	name := ""
	for _, s := range names {
		if name == "" {
			name = strings.TrimSpace(s)
		} else if len(name) > 0 {
			if strings.Contains(name, "Macintosh") && s != "" {
				name = strings.TrimSpace(s)
			} else if strings.Contains(name, s) {
				name = strings.TrimSpace(s)
			} else if !strings.Contains(s, name) {
				if strings.Contains(name, "iPhone") ||
					strings.Contains(name, "iPad") {
					s = strings.Trim(s, "Mac OS X")
				}

				if s != "" {
					name += " " + strings.TrimSpace(s)
				}
			}
			break
		}

		if strings.Contains(name, "Windows NT") {
			name = getWinOsNameWithWinNT(name)
			break
		}
	}

	if name != "" {
		osName = name
	}

	return osName
}

// GetBrowserName 获取浏览器名称
func GetBrowserName(userAgent string) string {
	deviceName := "Unknown"

	levelNames := ":VivoBrowser:QQDownload:QQBrowser:QQ:MQQBrowser:MicroMessenger:TencentTraveler:LBBROWSER:TaoBrowser:BrowserNG:UCWEB:TwonkyBeamBrowser:NokiaBrowser:OviBrowser:NF-Browser:OneBrowser:Obigo:DiigoBrowser:baidubrowser:baiduboxapp:xiaomi:Redmi:MI:Lumia:Micromax:MSIEMobile:IEMobile:EdgiOS:Yandex:Mercury:Openwave:TouchPad:UBrowser:Presto:Maxthon:MetaSr:Trident:Opera:IEMobile:Edge:Chrome:Chromium:OPR:CriOS:Firefox:FxiOS:fennec:CrMo:Safari:Nexus One:Nexus S:Nexus:Blazer:teashark:bolt:HTC:Dell:Motorola:Samsung:LG:Sony:SonyST:SonyLT:SonyEricsson:Asus:Palm:Vertu:Pantech:Fly:Wiko:i-mobile:Alcatel:Nintendo:Amoi:INQ:ONEPLUS:Tapatalk:PDA:Novarra-Vision:NetFront:Minimo:FlyFlow:Dolfin:Nokia:Series:AppleWebKit:Mobile:Mozilla:Version:"

	var regStrArr []string
	namesArr := strings.Split(strings.Trim(levelNames, ":"), ":")
	for _, name := range namesArr {
		regStrArr = append(regStrArr, fmt.Sprintf("(%s[\\s?\\/0-9.]+)", name))
	}
	regexpStr := fmt.Sprintf("(?i:%s)", strings.Join(regStrArr, "|"))
	nameRe, _ := regexp.Compile(regexpStr)
	names := nameRe.FindAllString(userAgent, -1)

	level := 0
	for _, name := range names {
		replaceRe, _ := regexp.Compile("(?i:[\\s?\\/0-9.]+)")
		n := replaceRe.ReplaceAllString(name, "")
		l := strings.Index(levelNames, fmt.Sprintf(":%s:", n))
		if level == 0 {
			deviceName = strings.TrimSpace(name)
		}

		if l >= 0 && (level == 0 || level > l) {
			level = l
			deviceName = strings.TrimSpace(name)
		}
	}

	return deviceName
}

// GenerateFingerprint 从请求头中提取信息生成指纹
func GenerateFingerprint(r *http.Request) string {
	// 选择一些请求头字段作为浏览器指纹的基础
	headers := []string{
		r.Header.Get("User-Agent"),
		r.Header.Get("Accept"),
		r.Header.Get("Accept-Language"),
		r.Header.Get("Accept-Encoding"),
		r.Header.Get("DNT"),
		r.Header.Get("Connection"),
	}

	// 将 headers 排序，以保证一致性
	sort.Strings(headers)

	// 将排序后的 headers 拼接为一个单一的字符串
	fingerprint := strings.Join(headers, "|")

	// 使用 sha256 哈希生成一个固定长度的指纹
	hash := sha256.Sum256([]byte(fingerprint))

	// 将哈希值编码为 base64 字符串
	return base64.StdEncoding.EncodeToString(hash[:])
}
