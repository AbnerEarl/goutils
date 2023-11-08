/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/11/7 3:14 PM
 * @desc: about the role of class.
 */

package emails

import "fmt"

func SendEmail(serverHost string, serverPort int, fromEmail, fromPasswd string, toers, ccers []string, subject, body string) error {
	var m *Message
	m = NewMessage()
	// 收件人可以有多个，故用此方式
	if len(toers) > 0 {
		m.SetHeader("To", toers...)
	} else {
		return fmt.Errorf("the reciever is not nil")
	}

	//抄送列表
	if len(ccers) > 0 {
		m.SetHeader("Cc", ccers...)
	}

	// 发件人，第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", fromEmail, "")

	// 主题
	m.SetHeader("Subject", subject)

	// 正文
	m.SetBody("text/html", body)

	d := NewPlainDialer(serverHost, serverPort, fromEmail, fromPasswd)
	// 发送
	err := d.DialAndSend(m)
	return err
}
