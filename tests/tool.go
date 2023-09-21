/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/9/14 10:16 AM
 * @desc: about the role of class.
 */

package tests

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func MapToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + url.QueryEscape(key) + "=" + url.QueryEscape(val)
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

func Case(name string, t *testing.T, fn func(items ...interface{})) {
	Convey(name, t, fn)
}
