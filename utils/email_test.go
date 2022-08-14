package utils

import (
	"testing"
)

// @Author: Derek
// @Description:
// @Date: 2022/8/13 23:12
// @Version 1.0

func TestSendEmail(t *testing.T) {
	username := "xx@qq.com"
	host := "smtp.qq.com"
	password := "客户端专用密码"
	port := 465

	subject := "主题"
	content := "内容"
	contentType := "text/html"
	attach := "" // 附件
	to := []string{"xx@qq.com"}

	message := NewEmailMessage(username, subject, contentType, content, attach, to, nil)
	email := NewEmailClient(host, username, password, port, message)

	ok, err := email.SendMessage()
	if err != nil || !ok {
		t.Logf("发送邮件失败了: %s", err)
	}

}
