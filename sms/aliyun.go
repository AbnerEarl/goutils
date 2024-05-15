/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/5/15 10:18
 * @desc: about the role of class.
 */

package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	one       sync.Once
	SmsClient = &smsClient{}
)

type smsConfig struct {
	RegionId        string
	AccessKeyId     string
	AccessKeySecret string
}

type smsClient struct {
	client *sdk.Client
	config *smsConfig
}

type SMSContent struct {
	PhoneNumbers  string
	SignName      string
	TemplateCode  string
	TemplateParam string // JSON串
}

// SendSms https://api.aliyun.com/#/?product=Dysmsapi&api=SendSms&params={}&tab=DEMO&lang=GO
func (sms *smsClient) SendSms(content *SMSContent) bool {
	if sms.client == nil {
		one.Do(func() {
			sms.config = &smsConfig{
				// https://help.aliyun.com/document_detail/40654.html?spm=a2c6h.13066369.0.0.54a120f8MncWFW
				RegionId:        "cn-hangzhou",
				AccessKeyId:     "",
				AccessKeySecret: "",
			}
			client, err := sdk.NewClientWithAccessKey(sms.config.RegionId, sms.config.AccessKeyId, sms.config.AccessKeySecret)
			if err != nil {
				panic(err)
			}
			sms.client = client
		})
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = sms.config.RegionId

	// 自定义
	request.QueryParams["PhoneNumbers"] = content.PhoneNumbers
	request.QueryParams["SignName"] = content.SignName
	request.QueryParams["TemplateCode"] = content.TemplateCode
	request.QueryParams["TemplateParam"] = content.TemplateParam

	response, err := sms.client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpContentString())

	return response.IsSuccess()
}

// GenValidateCode 生成随机数验证码
func GenValidateCode(len int) string {
	numbers := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < len; i++ {
		fmt.Fprintf(&sb, "%d", numbers[rand.Intn(10)])
	}
	return sb.String()
}

func StartSend() {
	content := &SMSContent{
		PhoneNumbers:  "110",
		SignName:      "派大星",
		TemplateCode:  "SMS_173230275",
		TemplateParam: `{"code":"` + GenValidateCode(6) + `"}`,
	}
	SmsClient.SendSms(content)
}
