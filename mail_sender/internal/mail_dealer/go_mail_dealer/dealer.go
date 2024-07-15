package gomaildealer

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type Dealer struct {
	md        *gomail.Dialer
	mail_body string
}

func New(host string, port int, login, password string, mailBody string) *Dealer {
	return &Dealer{
		md:        gomail.NewDialer("smtp."+host, port, login, password),
		mail_body: mailBody,
	}
}

func (d *Dealer) Send(recepient string) error {
	mess := gomail.NewMessage()
	mess.SetHeader("Subject", "Welcome to Our Service!")
	mess.SetHeader("From", "andreysavinandreas@yandex.com")
	mess.SetBody("text/html", strings.ReplaceAll(string(d.mail_body), "[USER_NAME]", d.md.Username))
	mess.SetAddressHeader("To", recepient, recepient)
	err := d.md.DialAndSend(mess)
	if err != nil {
		return err
	}
	return nil
}