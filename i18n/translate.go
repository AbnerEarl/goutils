/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/5/10 14:04
 * @desc: about the role of class.
 */

package i18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"strings"
	"time"
)

type DictRequestHS struct {
	Source         string   `json:"source"`
	Words          []string `json:"words"`
	SourceLanguage string   `json:"source_language"`
	TargetLanguage string   `json:"target_language"`
}
type DictResponseHS struct {
	Details []struct {
		Detail string `json:"detail"`
		Extra  string `json:"extra"`
	} `json:"details"`
}

type DictResponseHSData struct {
	Result []struct {
		Ec struct {
			Basic struct {
				Explains []struct {
					Pos   string `json:"pos"`
					Trans string `json:"trans"`
				} `json:"explains"`
			} `json:"basic"`
		} `json:"ec"`
	} `json:"result"`
}

func QueryHS(word, srcLang, tarLang string) (string, map[string]string, error) {
	if len(srcLang) < 2 {
		srcLang = "en"
	}
	if len(tarLang) < 2 {
		tarLang = "zh"
	}

	client := &http.Client{}
	request := DictRequestHS{"youdao", []string{word}, srcLang, tarLang}
	buf, err := json.Marshal(request)
	if err != nil {
		return "", nil, err
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/detail/v1/?msToken=&X-Bogus=DFSzswVOQDaibrQ3tJHN7cppgiFh&_signature=_02B4Z6wo00001g0lO6gAAIDD-FrRNX0w-.4NJT8AAOfuf7", data)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16888064002651706; i18next=zh-CN; s_v_web_id=verify_ljtrq6kx_UW3ieIzP_8gQX_4abc_B8D8_AoHwuLysn026; ttcid=db98bce9149b4f09b905a71503d9331e36")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/?category=&home_language=zh&source_language=detect&target_language=zh&text=bad")
	req.Header.Set("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Microsoft Edge";v="114"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.43")
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("bad StatusCode: %d, body: %s", resp.StatusCode, string(bodyText))
	}

	var dictResponse DictResponseHS

	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		return "", nil, err
	}

	item := dictResponse.Details[0]
	jsonStr := item.Detail

	var HSData DictResponseHSData
	err = json.Unmarshal([]byte(jsonStr), &HSData)
	if err != nil {
		return "", nil, err
	}

	result := map[string]string{}
	value := ""
	for i, it := range HSData.Result[0].Ec.Basic.Explains {
		//fmt.Println(it.Pos, it.Trans)
		result[it.Pos] = it.Trans
		if i == 0 {
			s := strings.Split(it.Trans, "；")[0]
			s = strings.Split(s, ";")[0]
			s = strings.Split(s, "（")[0]
			s = strings.Split(s, "(")[0]
			value = strings.TrimSpace(s)
		}
	}
	return value, result, nil
}

type DictRequestCY struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponseCY struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type DictResponseBD struct {
	Errno int `json:"errno"`
	Data  []struct {
		K string `json:"k"`
		V string `json:"v"`
	} `json:"data"`
}

func QueryBD(word, transLang string) (string, map[string]string, error) {
	if len(transLang) < 3 {
		transLang = "en2zh"
	}
	client := &http.Client{}
	var data = strings.NewReader("kw=" + word)
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/sug", data)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,de;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	//req.Header.Set("Cookie", "PSTM=1700618016; BIDUPSID=001E381FAF9F12569F89E62DBEF27690; BAIDUID=4FA55789C172962A13C116AA64C29D2B:FG=1; BDUSS=JBRTQyUmY0cEFYa0o3WnNFMmV3czVwQU45dHduRWNTYU9JY0VyUXhTaTdKOGhsRUFBQUFBJCQAAAAAAAAAAAEAAAD1h209xMeyu7Tm1Nq1xNKy0O0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALuaoGW7mqBlM1; BDUSS_BFESS=JBRTQyUmY0cEFYa0o3WnNFMmV3czVwQU45dHduRWNTYU9JY0VyUXhTaTdKOGhsRUFBQUFBJCQAAAAAAAAAAAEAAAD1h209xMeyu7Tm1Nq1xNKy0O0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALuaoGW7mqBlM1; BDSFRCVID=TaDOJexroG3KIhJt2Q1au8KfXgKKvV3TDYLEOwXPsp3LGJLVYEd8EG0Pt_NFmZK-oxmHogKK0mOTH6KF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tJAtVC-5JCP3H48k-4QEbbQH-UnLq-QRLgOZ04n-ah05KtbIjRQNj5kZ5H7iKqvMW20j0h7m3UTdsq76Wh35K5tTQP6rLf5C2574KKJxbnQcDPbxLPKhy58BhUJiB5OLBan7_qvIXKohJh7FM4tW3J0ZyxomtfQxtNRJ0DnjtpChbC_RjjuMjTQbeU5eetjK2CntsJOOaCkhqxJOy4oWK441DgcJKj5D0GQWBPKyWbQiepvoD-Jc3M04X-o9-hvT-54e2p3FBUQPEJvRQft20b0EXUnJbUJaJm-jLJ7jWhk2Dq72y5jvQlRX5q79atTMfNTJ-qcH0KQpsIJM5-DWbT8IjHCJq6tjfnKeoIvM5nL_HPJYhP-_-P4Dep3i0URZ5mAqoq8MJ-bJH4jGMPOG5jkUbR7AbhLJQJTnaIQhtD38HnvJbtRh5UIUBUQxqlQ43bRT-pCy5KJvMh61yxRjhP-UyN3LWh37Je3lMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMoPPafD8abDIlj6taePDyqx5Ka43tHD7yWCv9-CJcOR59K4nnDUrbhMTDKPjQQ5IOabQ95prDSJb43MOZKxLg5n7Tbb8eBgvZ2UQw366Dsq0x0bO5DDuOQqoattIetIOMahkMal7xOKQoQlPK5JkgMx6MqpQJQeQ-5KQN3KJmfbL9bT3YjjTLjNuetTLHfKresJoq2RbhKROvhjRx-f0gyxoObtRxt23ABpc2tfDaHtK9yPjibUPU5RDeLU3kBgT9LMnx--t58h3_XhjZQt3bQttjQn3dBecLhqkEa-nVOb7TyU42hf47yaji0q4Hb6b9BJcjfU5MSlcNLTjpQT8r5MDOK5OuJRQ2QJ8BtCD5MIjP; H_WISE_SIDS=40445_40500_40080_60142_60175; BDSFRCVID_BFESS=TaDOJexroG3KIhJt2Q1au8KfXgKKvV3TDYLEOwXPsp3LGJLVYEd8EG0Pt_NFmZK-oxmHogKK0mOTH6KF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF_BFESS=tJAtVC-5JCP3H48k-4QEbbQH-UnLq-QRLgOZ04n-ah05KtbIjRQNj5kZ5H7iKqvMW20j0h7m3UTdsq76Wh35K5tTQP6rLf5C2574KKJxbnQcDPbxLPKhy58BhUJiB5OLBan7_qvIXKohJh7FM4tW3J0ZyxomtfQxtNRJ0DnjtpChbC_RjjuMjTQbeU5eetjK2CntsJOOaCkhqxJOy4oWK441DgcJKj5D0GQWBPKyWbQiepvoD-Jc3M04X-o9-hvT-54e2p3FBUQPEJvRQft20b0EXUnJbUJaJm-jLJ7jWhk2Dq72y5jvQlRX5q79atTMfNTJ-qcH0KQpsIJM5-DWbT8IjHCJq6tjfnKeoIvM5nL_HPJYhP-_-P4Dep3i0URZ5mAqoq8MJ-bJH4jGMPOG5jkUbR7AbhLJQJTnaIQhtD38HnvJbtRh5UIUBUQxqlQ43bRT-pCy5KJvMh61yxRjhP-UyN3LWh37Je3lMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMoPPafD8abDIlj6taePDyqx5Ka43tHD7yWCv9-CJcOR59K4nnDUrbhMTDKPjQQ5IOabQ95prDSJb43MOZKxLg5n7Tbb8eBgvZ2UQw366Dsq0x0bO5DDuOQqoattIetIOMahkMal7xOKQoQlPK5JkgMx6MqpQJQeQ-5KQN3KJmfbL9bT3YjjTLjNuetTLHfKresJoq2RbhKROvhjRx-f0gyxoObtRxt23ABpc2tfDaHtK9yPjibUPU5RDeLU3kBgT9LMnx--t58h3_XhjZQt3bQttjQn3dBecLhqkEa-nVOb7TyU42hf47yaji0q4Hb6b9BJcjfU5MSlcNLTjpQT8r5MDOK5OuJRQ2QJ8BtCD5MIjP; delPer=0; BAIDUID_BFESS=4FA55789C172962A13C116AA64C29D2B:FG=1; PSINO=7; BA_HECTOR=2k0g050k8g20a5ag8l8k848lds3jg41j3qs3c1v; ZFY=xwaPkYOQRAeOt:BzvcoXqNUCJAjFIJMnI2A4eogv6AxE:C; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; H_PS_PSSID=40445_40500_60119_40080_60142_60175_60269; H_WISE_SIDS_BFESS=40445_40500_40080_60142_60175; SL_G_WPT_TO=en; SL_GWPT_Show_Hide_tmp=1; SL_wptGlobTipTmp=1; ab_sr=1.0.1_ZDcyMTdhODVjZWIwZTE0NmJjYjg4ZTM0MmQ5ODZjNWExZDQ1ZjFkZWE1YmVjYjM5ZjMwYTE5NjEyNjFhYjc4ZWRhYjYzZTI4NzIzMWJkNTNkZTIxZTAzNzY4MTBkNmUzOTFmNDM5ZGMxZDE5N2YxMzgxOTM4MTU2NTM5ZmViYjg2NWMwNGY0MTFiMmM0NDA4NDE0MzU3ODgzOTVkZjFhYg==; RT=\"z=1&dm=baidu.com&si=f258d0cd-15bd-48c8-aaf5-a2e7e7665621&ss=lw0hypwo&sl=1&tt=2s9&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=3k6\"")
	req.Header.Set("Host", "fanyi.baidu.com")
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/mtpe-individual/multimodal?query="+word+"&lang="+transLang)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("bad StatusCode: %d, body: %s", resp.StatusCode, string(bodyText))
	}
	var dictResponse DictResponseBD
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		return "", nil, err
	}

	result := map[string]string{}
	value := ""
	for i, it := range dictResponse.Data {
		//fmt.Println(i, it.V)
		result[it.K] = it.V
		if i == 0 {
			s := strings.Split(it.V, "；")[0]
			s = strings.Split(s, ";")[0]
			s = strings.Split(s, "（")[0]
			s = strings.Split(s, "(")[0]
			if strings.Contains(s, " ") {
				s = strings.Split(s, " ")[1]
			}
			if strings.Contains(s, ".") {
				s = strings.Split(s, ".")[1]
			}
			value = strings.TrimSpace(s)
		}
	}

	return value, result, nil
}
func QueryCY(word, transType string) (string, map[string]string, error) {
	if len(transType) < 3 {
		transType = "en2zh"
	}
	client := &http.Client{}
	request := DictRequestCY{TransType: transType, Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		return "", nil, err
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("bad StatusCode: %d, body: %s", resp.StatusCode, string(bodyText))
	}
	var dictResponse DictResponseCY
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		return "", nil, err
	}
	result := map[string]string{"UK:": dictResponse.Dictionary.Prons.En, "US:": dictResponse.Dictionary.Prons.EnUs}
	value := ""
	for i, item := range dictResponse.Dictionary.Explanations {
		//fmt.Println(item)
		k := item[:strings.Index(item, ".")+1]
		v := item[strings.Index(item, ".")+1:]
		result[k] = v
		if i == 0 {
			s := strings.Split(v, "；")[0]
			s = strings.Split(s, ";")[0]
			s = strings.Split(s, "（")[0]
			s = strings.Split(s, "(")[0]
			value = strings.TrimSpace(s)
		}
	}

	return value, result, nil
}

var GoogleHost = "google.com"

// TranslationParams is a util struct to pass as parameter to indicate how to translate
type TranslationParams struct {
	From       string
	To         string
	Tries      int
	Delay      time.Duration
	GoogleHost string
}

// Translate translate a text using native tags offer by go language
func Translate(text string, from language.Tag, to language.Tag, googleHost ...string) (string, error) {
	if len(googleHost) != 0 && googleHost[0] != "" {
		GoogleHost = googleHost[0]
	}
	translated, err := translate(text, from.String(), to.String(), false, 2, 0)
	if err != nil {
		return "", err
	}

	return translated, nil
}

// TranslateWithParams translate a text with simple params as string
func TranslateWithParams(text string, params TranslationParams) (string, error) {
	if params.GoogleHost == "" {
		GoogleHost = "google.com"
	} else {
		GoogleHost = params.GoogleHost
	}
	translated, err := translate(text, params.From, params.To, true, params.Tries, params.Delay)
	if err != nil {
		return "", err
	}
	return translated, nil
}
