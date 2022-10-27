/*
发送电子邮件
*/

package email

import (
	"net/smtp"
	"strconv"
)

type Email struct {
	SenderName string
	SenderAddr string // sender email address
	Password   string // sender email password
	Host       string // SMTP server host
	Port       int    // SMTP server port
}

// 构建消息体
func msgConstructor(sendName, sendAddress, receiverName, subject, content string) []byte {
	return []byte("From: " + sendName + " <" + sendAddress + ">\n" +
		"To: " + receiverName + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		content)
}

// Send email
func (e *Email) Send(receiver, subject, content string) error {
	auth := smtp.PlainAuth("", e.SenderAddr, e.Password, e.Host)
	text := msgConstructor(e.SenderName, e.SenderAddr, receiver, subject, content)

	err := smtp.SendMail(e.Host+":"+strconv.Itoa(e.Port),
		auth,
		e.SenderAddr,
		[]string{receiver},
		text)
	return err
}
