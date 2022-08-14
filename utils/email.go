package utils

// @Author: Derek
// @Description: Send Email.
// @Date: 2022/8/13 23:03
// @Version 1.0

import (
	"gopkg.in/gomail.v2"
)

type EmailMessage struct {
	From        string
	To          []string
	Cc          []string
	Subject     string
	ContentType string
	Content     string
	Attach      string
}

// NewEmailMessage return EmailMessage.
// from: 发件人
// subject: 标题
// contentType: 内容的类型 text/plain text/html
// attach: 附件
// to: 收件人
// cc: 抄送人
func NewEmailMessage(from, subject, contentType, content, attach string, to, cc []string) *EmailMessage {
	return &EmailMessage{
		From:        from,
		Subject:     subject,
		ContentType: contentType,
		Content:     content,
		To:          to,
		Cc:          cc,
		Attach:      attach,
	}
}

// EmailClient 发送客户端
type EmailClient struct {
	Host     string
	Port     int
	Username string
	Password string
	Message  *EmailMessage
}

// NewEmailClient 返回一个邮件客户端
// host smtp地址
// username 用户名
// password 密码
// port 端口
func NewEmailClient(host, username, password string, port int, message *EmailMessage) *EmailClient {
	return &EmailClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Message:  message,
	}
}

// SendMessage 发送邮件
func (c *EmailClient) SendMessage() (bool, error) {
	dm := gomail.NewMessage()
	dm.SetHeader("From", dm.FormatAddress(c.Message.From, "xx官方"))
	dm.SetHeader("To", c.Message.To...)

	if len(c.Message.Cc) != 0 {
		dm.SetHeader("Cc", c.Message.Cc...)
	}

	dm.SetHeader("Subject", c.Message.Subject)
	dm.SetBody(c.Message.ContentType, c.Message.Content)

	if c.Message.Attach != "" {
		dm.Attach(c.Message.Attach)
	}

	e := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	//if 587 == c.Port {
	//e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	//}
	if err := e.DialAndSend(dm); err != nil {
		return false, err
	}
	return true, nil
}
