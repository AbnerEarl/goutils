package captcha

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"strings"
)

// 设置自带的store
var store = base64Captcha.DefaultMemStore

func CaptMake() (id, b64s string, err error) {
	// 生成验证码
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, store)
	lid, lb64s, lerr := captcha.Generate()
	return lid, lb64s, lerr
}

func CaptVerify(id string, capt string) bool {
	// 验证captcha是否正确
	if len(strings.TrimSpace(id)) > 0 && len(strings.TrimSpace(capt)) > 0 {
		return store.Verify(id, capt, true)
	}
	return false
	//return true
}

func String2Md5(data string) string {
	md5Value := md5.New()
	md5Value.Write([]byte(data))
	result := hex.EncodeToString(md5Value.Sum(nil))
	return result
}
