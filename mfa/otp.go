/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/22 16:32
 * @desc: about the role of class.
 */

package mfa

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/AbnerEarl/goutils/files"
	"net/url"
	"strings"
	"time"
)

func GetSecret(baseSize uint) string {
	randomStr := randStr(baseSize)
	return strings.ToUpper(randomStr)
}

func randStr(strSize uint) string {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, strSize)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

type FreeOtpAuthenticator2FaSha1 struct {
	Secret string `json:"secret"` //The base32NoPaddingEncodedSecret parameter is an arbitrary key value encoded in Base32 according to RFC 3548. The padding specified in RFC 3548 section 2.2 is not required and should be omitted.
	Period int64  `json:"period"` //更新周期单位秒，根据客户端程序设置，一般为：30
	Digits uint   `json:"digits"` //数字数量,根据客户端程序设置，一般为：6 到 8
	Label  string `json:"label"`  //格式为："应用名称:账号名称"
	Issuer string `json:"issuer"` //作者信息
}

func (m *FreeOtpAuthenticator2FaSha1) QrString() (qr string) {
	//规范文档 https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	//otpauth://totp/ACME%20Co:john.doe@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30
	return fmt.Sprintf(`otpauth://totp/%s?secret=%s&issuer=%s&algorithm=SHA1&digits=%d&period=%d`, url.QueryEscape(m.Label), m.Secret, url.QueryEscape(m.Issuer), m.Digits, m.Period)
}

func (m *FreeOtpAuthenticator2FaSha1) VerifyCode(code uint32) bool {
	// 为了考虑时间误差，判断前当前时间及前后30秒时间
	// 当前google值
	if m.getCode(m.Secret, 0) == code {
		return true
	}

	// 前30秒google值
	if m.getCode(m.Secret, -m.Period) == code {
		return true
	}

	// 后30秒google值
	if m.getCode(m.Secret, m.Period) == code {
		return true
	}

	return false
}

// 获取Google Code
func (m *FreeOtpAuthenticator2FaSha1) getCode(secret string, offset int64) uint32 {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix() + offset
	return m.oneTimePassword(key, toBytes(epochSeconds/m.Period))
}

// from https://github.com/robbiev/two-factor-auth/blob/master/main.go
func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func (m *FreeOtpAuthenticator2FaSha1) oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	d := uint32(1)
	//取十进制的余数
	for i := uint(0); i < m.Digits && i < 8; i++ {
		d *= 10
	}
	// one million is the first number with 7 digits so the remainder of the division will always return < 7 digits
	pwd := number % d
	return pwd
}

// MakeGoogleAuthenticator 获取key&t对应的验证码
// key 秘钥
// t 1970年的秒
func (m *FreeOtpAuthenticator2FaSha1) MakeGoogleAuthenticator(key string, t int64) (string, error) {
	hs, e := hmacSha1(key, t/m.Period)
	if e != nil {
		return "", e
	}
	num := lastBit4byte(hs)
	d := uint32(1)
	//取十进制的余数
	for i := uint(0); i < m.Digits && i < 8; i++ {
		d *= 10
	}
	v := num % d
	intFormat := fmt.Sprintf("%%0%dd", m.Digits) //数字长度补零
	return fmt.Sprintf(intFormat, v), nil
}

// MakeGoogleAuthenticatorForNow 获取key对应的验证码
func (m *FreeOtpAuthenticator2FaSha1) MakeGoogleAuthenticatorForNow(key string) (string, error) {
	return m.MakeGoogleAuthenticator(key, time.Now().Unix())
}

func lastBit4byte(hmacSha1 []byte) uint32 {
	if len(hmacSha1) != sha1.Size {
		return 0
	}
	offsetBits := int8(hmacSha1[len(hmacSha1)-1]) & 0x0f
	p := (uint32(hmacSha1[offsetBits]) << 24) | (uint32(hmacSha1[offsetBits+1]) << 16) | (uint32(hmacSha1[offsetBits+2]) << 8) | (uint32(hmacSha1[offsetBits+3]) << 0)
	return p & 0x7fffffff
}

func hmacSha1(key string, t int64) ([]byte, error) {
	decodeKey, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(key)
	if err != nil {
		return nil, err
	}

	cData := make([]byte, 8)
	binary.BigEndian.PutUint64(cData, uint64(t))

	h1 := hmac.New(sha1.New, decodeKey)
	_, e := h1.Write(cData)
	if e != nil {
		return nil, e
	}
	return h1.Sum(nil), nil
}

func ClientResourcePackage() map[string]map[string]string {
	dirPath := files.GetAbPath() + "mfa/resource/"
	return map[string]map[string]string{
		"google_authenticator": {
			"ios":     "https://apps.apple.com/tw/app/google-authenticator/id388497605",
			"android": dirPath + "google-authenticator.apk.zip",
		},
		"freeotp_authenticator": {
			"ios":     "https://apps.apple.com/us/app/freeotp-authenticator/id872559395",
			"android": dirPath + "free-otp.apk.zip",
		},
		"microsoft_authenticator": {
			"ios":     "https://apps.apple.com/us/app/microsoft-authenticator/id983156458",
			"android": dirPath + "microsoft-authenticator.apk.zip",
		},
	}
}
