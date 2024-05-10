/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/5/10 15:59
 * @desc: about the role of class.
 */

package i18n

import "errors"

var errBadNetwork = errors.New("bad network, please check your internet connection")
var errBadRequest = errors.New("bad request, request on google translate api isn't working")
