/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/7/13 6:13 PM
 * @desc: about the role of class.
 */

package jwts

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type CustomClaims struct {
	UserInfo interface{}
	jwt.StandardClaims
}

func GenerateToken(userInfo interface{}, secret string, issuer string, audience string, expiredMinutes int64) (string, error) {
	hmacSampleSecret := []byte(secret)
	token := jwt.New(jwt.SigningMethodHS256)
	nowTime := time.Now().Unix()
	token.Claims = CustomClaims{
		UserInfo: userInfo,
		StandardClaims: jwt.StandardClaims{
			NotBefore: nowTime,
			ExpiresAt: nowTime + expiredMinutes*60,
			Issuer:    issuer,
			Audience:  audience,
		},
	}
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func ParseToken(tokenString string, secret string) (*CustomClaims, error) {
	var hmacSampleSecret = []byte(secret)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(*CustomClaims)
	return claims, nil
}

func RefreshToken(tokenString string, secret string, addMinutes int64) (string, error) {
	var hmacSampleSecret = []byte(secret)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return "", err
	}
	token.Claims.(*CustomClaims).ExpiresAt = time.Now().Unix() + addMinutes*60
	newTokenString, err := token.SignedString(hmacSampleSecret)
	return newTokenString, err
}
