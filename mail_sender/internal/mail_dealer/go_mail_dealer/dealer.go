package gomaildealer

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type Dealer struct {
	author    string
	md        *gomail.Dialer
	mail_body string
}

func New(host string, port int, login, password string, mailBody string) *Dealer {
	return &Dealer{
		author:    login,
		md:        gomail.NewDialer("smtp."+host, port, login, password),
		mail_body: mailBody,
	}
}

func (d *Dealer) Send(recipient string) error {
	mess := gomail.NewMessage()
	mess.SetHeader("Subject", "Welcome to Our Service!")
	mess.SetHeader("From", d.author)
	mess.SetBody("text/html", strings.ReplaceAll(string(d.mail_body), "[USER_NAME]", d.md.Username))
	mess.SetAddressHeader("To", recipient, recipient)
	err := d.md.DialAndSend(mess)
	if err != nil {
		return err
	}
	return nil
}
