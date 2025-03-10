package smail

import (
	"blog/pkg/configs"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

var config configs.EmailConfig

func InitSmtp(EmailConfig configs.EmailConfig) {
	config = EmailConfig
}

// SendEmail 发送邮件
// to:对方邮箱 subject:邮箱主题 isHTML:是否是html格式 text:文本信息
func SendEmail(to string, subject string, isHTML bool, text string) error {

	var e = email.NewEmail()

	e.From = fmt.Sprintf("%s <%s>", "", config.Username)

	e.To = []string{to}

	e.Subject = subject

	if isHTML {
		e.HTML = []byte(text)
	} else {
		e.Text = []byte(text)
	}

	if err := e.Send(config.Addr, smtp.PlainAuth("", config.Username, config.Password, config.Host)); err != nil {
		return err
	}

	return nil
}
