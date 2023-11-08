/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/11/7 2:57 PM
 * @desc: about the role of class.
 */

package emails

import (
	"mime"
	"mime/quotedprintable"
	"strings"
)

var newQPWriter = quotedprintable.NewWriter

type mimeEncoder struct {
	mime.WordEncoder
}

var (
	bEncoding     = mimeEncoder{mime.BEncoding}
	qEncoding     = mimeEncoder{mime.QEncoding}
	lastIndexByte = strings.LastIndexByte
)
