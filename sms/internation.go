/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/22 11:01
 * @desc: about the role of class.
 */

package sms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func SendWithTwilio(accountSid, authToken, from, to, message string) (map[string]interface{}, error) {
	// visit https://www.twilio.com/ to get auth info.

	requestUrl := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	//Set up the data for the text message
	messageDate := url.Values{}
	messageDate.Set("To", to)
	messageDate.Set("From", from)
	messageDate.Set("Body", message)
	msgDataReader := *strings.NewReader(messageDate.Encode())

	//Create HTTP request client
	client := &http.Client{}
	request, _ := http.NewRequest("POST", requestUrl, &msgDataReader)
	request.SetBasicAuth(accountSid, authToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	return res, err
}

func SendWithTextMagic() {
	// visit https://textmagic.com/ to get auth info.

}

func SendWithCourier() {
	// visit https://www.courier.com/ to get auth info.
}
