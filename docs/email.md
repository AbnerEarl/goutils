电子邮件在日常工作中有很大用途，凡项目或任务，有邮件来往可避免扯皮背锅。而在一些自动化的应用场合，也使用得广泛，特别是系统监控方面，如果在资源使用达到警戒线之前自动发邮件通知运维人员，能消除隐患于前期，而不至于临时临急去做善后方案。

对于多人协合（不管是不是异地）场合，邮件也有用武之地，当有代码或文档更新时，自动发邮件通知项目成员或领导，提醒各方人员知晓并及时更新。

说到发邮件，不得不提用程序的方式实现。下面就来为大家介绍一下怎么使用Go语言来实现发送电子邮件。Go语言拥有大量的库，非常方便使用。

Go语言使用 emails 包来发送邮箱，代码如下所示：

```go
package main

import (
"strings"
"github.com/AbnerEarl/goutils/emails"
)

type EmailParam struct {
// ServerHost 邮箱服务器地址，如腾讯邮箱为smtp.qq.com
ServerHost string
// ServerPort 邮箱服务器端口，如腾讯邮箱为465
ServerPort int
// FromEmail　发件人邮箱地址
FromEmail string
// FromPasswd 发件人邮箱密码（注意，这里是明文形式），TODO：如果设置成密文？
FromPasswd string
// Toers 接收者邮件，如有多个，则以英文逗号(“,”)隔开，不能为空
Toers string
// CCers 抄送者邮件，如有多个，则以英文逗号(“,”)隔开，可以为空
CCers string
}

// 全局变量，因为发件人账号、密码，需要在发送时才指定
// 注意，由于是小写，外面的包无法使用
var serverHost, fromEmail, fromPasswd string
var serverPort int

var m *emails.Message

func InitEmail(ep *EmailParam) {
toers := []string{}

    serverHost = ep.ServerHost
    serverPort = ep.ServerPort
    fromEmail = ep.FromEmail
    fromPasswd = ep.FromPasswd
   
    m = emails.NewMessage()
   
    if len(ep.Toers) == 0 {
        return
    }

    for _, tmp := range strings.Split(ep.Toers, ",") {
        toers = append(toers, strings.TrimSpace(tmp))
    }
   
    // 收件人可以有多个，故用此方式
    m.SetHeader("To", toers...)

    //抄送列表
    if len(ep.CCers) != 0 {
        for _, tmp := range strings.Split(ep.CCers, ",") {
            toers = append(toers, strings.TrimSpace(tmp))
        }
        m.SetHeader("Cc", toers...)
    }

    // 发件人
    // 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
    m.SetAddressHeader("From", fromEmail, "")
}

// SendEmail body支持html格式字符串
func SendEmail(subject, body string) {
// 主题
m.SetHeader("Subject", subject)

    // 正文
    m.SetBody("text/html", body)

    d := emails.NewPlainDialer(serverHost, serverPort, fromEmail, fromPasswd)
    // 发送
    err := d.DialAndSend(m)
    if err != nil {
        panic(err)
    }
}

func main() {
serverHost := "smtp.qq.com"
serverPort := 465
fromEmail := "xxxxxxx@qq.com"     //发件人邮箱
fromPasswd := "xumkkzfscmxxxxxx"    //授权码

    myToers := "xxxxxxx@qq.com" // 收件人邮箱，逗号隔开
    myCCers := "" //"readchy@163.com"
   
    subject := "这是主题"
    body := `这是正文<br>
             Hello <a href = "http://c.biancheng.net/">C语言中文网</a>`
    // 结构体赋值
    myEmail := &EmailParam {
        ServerHost: serverHost,
        ServerPort: serverPort,
        FromEmail:  fromEmail,
        FromPasswd: fromPasswd,
        Toers:      myToers,
        CCers:      myCCers,
    }
   
    InitEmail(myEmail)
    SendEmail(subject, body)
}


```
使用自定义客户端发放邮件需要以下两个要素:

发送方的邮箱必须开启 stmt 和 pop3 通道，以 qq 邮箱为例，登陆 qq 邮箱 -> 设置 -> 账户 -> 开启 pop3 和 stmt 服务


开启后会获得该账户的授权码，如果忘记也可以重新生成。


可以支持很多种邮箱，例如outlook，封装之后代码示例：

```go

import (
"strings"
"github.com/AbnerEarl/goutils/emails"
)

func SendEmail(serverHost string, serverPort int, fromEmail, fromPasswd string, toers, ccers []string, subject, body string) error {
	var m *emails.Message
	m = emails.NewMessage()
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

	d := emails.NewPlainDialer(serverHost, serverPort, fromEmail, fromPasswd)
	// 发送
	err := d.DialAndSend(m)
	return err
}


func main() {
    serverHost := "smtp.office365.com"
    serverPort := 587
    fromEmail := "ilovemitu@outlook.com" //发件人邮箱
    fromPasswd := "xxxxxx"  //授权码
    
    myToers := []string{"ilovemitu@outlook.com"} // 收件人邮箱
    myCCers := []string{}                        //"readchy@163.com"
    
    subject := "这是主题"
    body := `这是正文<br>Hello <a href = "http://c.biancheng.net/">C语言中文网</a>`
    // 结构体赋值
    utils.SendEmail(serverHost, serverPort, fromEmail, fromPasswd, myToers, myCCers, subject, body)
}

```

更多详细使用举例：

```go

package emails_test

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"time"

	"github.com/AbnerEarl/goutils/emails"
)

func Example() {
	m := emails.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", "bob@example.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.Attach("/home/Alex/lolcat.jpg")

	d := emails.NewDialer("smtp.example.com", 587, "user", "123456")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// A daemon that listens to a channel and sends all incoming messages.
func Example_daemon() {
	ch := make(chan *emails.Message)

	go func() {
		d := emails.NewDialer("smtp.example.com", 587, "user", "123456")

		var s emails.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err := emails.Send(s, m); err != nil {
					log.Print(err)
				}
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						panic(err)
					}
					open = false
				}
			}
		}
	}()

	// Use the channel in your program to send emails.

	// Close the channel to stop the mail daemon.
	close(ch)
}

// Efficiently send a customized newsletter to a list of recipients.
func Example_newsletter() {
	// The list of recipients.
	var list []struct {
		Name    string
		Address string
	}

	d := emails.NewDialer("smtp.example.com", 587, "user", "123456")
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := emails.NewMessage()
	for _, r := range list {
		m.SetHeader("From", "no-reply@example.com")
		m.SetAddressHeader("To", r.Address, r.Name)
		m.SetHeader("Subject", "Newsletter #1")
		m.SetBody("text/html", fmt.Sprintf("Hello %s!", r.Name))

		if err := emails.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r.Address, err)
		}
		m.Reset()
	}
}

// Send an email using a local SMTP server.
func Example_noAuth() {
	m := emails.NewMessage()
	m.SetHeader("From", "from@example.com")
	m.SetHeader("To", "to@example.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Hello!")

	d := emails.Dialer{Host: "localhost", Port: 587}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Send an email using an API or postfix.
func Example_noSMTP() {
	m := emails.NewMessage()
	m.SetHeader("From", "from@example.com")
	m.SetHeader("To", "to@example.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Hello!")

	s := emails.SendFunc(func(from string, to []string, msg io.WriterTo) error {
		// Implements you email-sending function, for example by calling
		// an API, or running postfix, etc.
		fmt.Println("From:", from)
		fmt.Println("To:", to)
		return nil
	})

	if err := emails.Send(s, m); err != nil {
		panic(err)
	}
	// Output:
	// From: from@example.com
	// To: [to@example.com]
}

var m *emails.Message

func ExampleSetCopyFunc() {
	m.Attach("foo.txt", emails.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write([]byte("Content of foo.txt"))
		return err
	}))
}

func ExampleSetHeader() {
	h := map[string][]string{"Content-ID": {"<foo@bar.mail>"}}
	m.Attach("foo.jpg", emails.SetHeader(h))
}

func ExampleRename() {
	m.Attach("/tmp/0000146.jpg", emails.Rename("picture.jpg"))
}

func ExampleMessage_AddAlternative() {
	m.SetBody("text/plain", "Hello!")
	m.AddAlternative("text/html", "<p>Hello!</p>")
}

func ExampleMessage_AddAlternativeWriter() {
	t := template.Must(template.New("example").Parse("Hello {{.}}!"))
	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		return t.Execute(w, "Bob")
	})
}

func ExampleMessage_Attach() {
	m.Attach("/tmp/image.jpg")
}

func ExampleMessage_Embed() {
	m.Embed("/tmp/image.jpg")
	m.SetBody("text/html", `<img src="cid:image.jpg" alt="My image" />`)
}

func ExampleMessage_FormatAddress() {
	m.SetHeader("To", m.FormatAddress("bob@example.com", "Bob"), m.FormatAddress("cora@example.com", "Cora"))
}

func ExampleMessage_FormatDate() {
	m.SetHeaders(map[string][]string{
		"X-Date": {m.FormatDate(time.Now())},
	})
}

func ExampleMessage_SetAddressHeader() {
	m.SetAddressHeader("To", "bob@example.com", "Bob")
}

func ExampleMessage_SetBody() {
	m.SetBody("text/plain", "Hello!")
}

func ExampleMessage_SetDateHeader() {
	m.SetDateHeader("X-Date", time.Now())
}

func ExampleMessage_SetHeader() {
	m.SetHeader("Subject", "Hello!")
}

func ExampleMessage_SetHeaders() {
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress("alex@example.com", "Alex")},
		"To":      {"bob@example.com", "cora@example.com"},
		"Subject": {"Hello"},
	})
}

func ExampleSetCharset() {
	m = emails.NewMessage(emails.SetCharset("ISO-8859-1"))
}

func ExampleSetEncoding() {
	m = emails.NewMessage(emails.SetEncoding(emails.Base64))
}

func ExampleSetPartEncoding() {
	m.SetBody("text/plain", "Hello!", emails.SetPartEncoding(emails.Unencoded))
}

```


